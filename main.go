package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var forbiddenPorts = map[string]bool{
	"20": true, "21": true, "22": true, "23": true, "25": true, "53": true, "80": true, "443": true, "3306": true, "5432": true, "27114": true,
}

// Helper to prompt and validate entry file

func main() {
	var containerName string
	if len(os.Args) > 1 {
		cmd := os.Args[1]
		if cmd == "delete" {
			if len(os.Args) > 2 {
				containerName = os.Args[2]
			}
			docifyDeleteParam(containerName)
			return
		} else if cmd == "logs" {
			if len(os.Args) > 2 {
				containerName = os.Args[2]
			}
			showContainerLogs(containerName)
			return
		} else if cmd == "show" {
			showAllContainers()
			return
		}
	}
	docifyRunWithSignal()
}

func detectBackendPort() string {
	// 1. Check .env file for PORT
	if fileExists(".env") {
		data, err := os.ReadFile(".env")
		if err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(strings.TrimSpace(line), "PORT=") {
					port := strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
					if _, err := strconv.Atoi(port); err == nil {
						return port
					}
				}
			}
		}
	}
	// 2. Fallback: Try to detect port from any .js file
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".js") {
			data, err := os.ReadFile(f.Name())
			if err == nil {
				lines := strings.Split(string(data), "\n")
				for _, line := range lines {
					// Node.js: listen(3000)
					if strings.Contains(line, "listen(") {
						parts := strings.Split(line, "listen(")
						if len(parts) > 1 {
							rest := parts[1]
							rest = strings.TrimSpace(rest)
							rest = strings.Trim(rest, ")")
							rest = strings.Split(rest, ",")[0]
							port := strings.Trim(rest, ") ")
							if _, err := strconv.Atoi(port); err == nil {
								return port
							}
						}
					}
					// Node.js: const/let/var PORT = 3000;
					lineTrim := strings.TrimSpace(line)
					if strings.HasPrefix(lineTrim, "const PORT =") || strings.HasPrefix(lineTrim, "let PORT =") || strings.HasPrefix(lineTrim, "var PORT =") {
						parts := strings.Split(lineTrim, "=")
						if len(parts) == 2 {
							port := strings.TrimSpace(strings.Trim(parts[1], ";"))
							if _, err := strconv.Atoi(port); err == nil {
								return port
							}
						}
					}
				}
			}
		}
	}
	// Default fallback
	return ""
}

func detectDBUsage() string {
	// Check for MongoDB usage in any .js file or .env
	if fileExists(".env") {
		data, err := os.ReadFile(".env")
		if err == nil && strings.Contains(string(data), "MONGO_URL") {
			return "mongodb"
		}
	}
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".js") {
			data, err := os.ReadFile(f.Name())
			if err == nil && (strings.Contains(string(data), "mongodb") || strings.Contains(string(data), "mongoose")) {
				return "mongodb"
			}
		}
	}
	return "none"
}

func docifyRunWithSignal() {
	// No Ctrl+C wait, script will exit after container run
	var response string
	var dbType string
	var entryFile string
	var projectType string
	dbType = detectDBUsage()
	fmt.Println("\n🚀 Welcome to Docify Docker Automation!")
	fmt.Println("----------------------------------------")
	if dbType == "none" {
		fmt.Print("🟢 Is Docker Desktop running? (y/n): ")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			if !isDockerRunning() {
				fmt.Println("❌ You entered 'y' but Docker Desktop is NOT running. Please start Docker Desktop.")
				return
			}
		} else {
			fmt.Println("❌ Please start Docker Desktop to continue.")
			return
		}
		projectType = detectProjectType()
		ef, ok := promptEntryFile(projectType)
		if !ok {
			return
		}
		entryFile = ef
	} else {
		fmt.Print("🟢 Is Docker Desktop running? (y/n): ")
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" {
			if !isDockerRunning() {
				fmt.Println("❌ You entered 'y' but Docker Desktop is NOT running. Please start Docker Desktop.")
				return
			}
		} else {
			fmt.Println("❌ Please start Docker Desktop to continue.")
			return
		}

		projectType = detectProjectType()
		ef, ok := promptEntryFile(projectType)
		if !ok {
			return
		}
		entryFile = ef

		fmt.Print("🟢 Is your local MongoDB server running? (y/n): ")
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("❌ Please start your local MongoDB server to continue.")
			return
		}
	}
	// projectType already detected above
	imageName := "docify-app"
	fmt.Println("\n----------------------------------------")
	fmt.Print("📦 Enter container name (default: docify-app): ")
	var containerName string
	fmt.Scanln(&containerName)
	fmt.Print("containerName: ", containerName)
	if containerName == "" {
		containerName = "docify-app"
	}

	backendPort := detectBackendPort()
	if backendPort == "" || !isValidPort(backendPort) {
		fmt.Print("🔢 I am unable to fetch port of your backend. Please enter the port your backend runs on: ")
		fmt.Scanln(&backendPort)
		if !isValidPort(backendPort) {
			fmt.Println("❌ Invalid or forbidden port.")
			return
		}
	}
	fmt.Printf("\n🔎 Detected backend port: %s\n", backendPort)
	dbType = detectDBUsage()
	if dbType == "mongodb" {
		fmt.Println("✅ Detected MongoDB usage. Your app will use local MongoDB at host.docker.internal:27017.")
	} else {
		fmt.Println("✅ No database detected in your project.")
	}

	// Dockerfile
	fmt.Println("\n----------------------------------------")
	if fileExists("Dockerfile") {
		os.Remove("Dockerfile")
	}
	generateDockerfile(projectType, backendPort, entryFile)
	fmt.Println("📦 Dockerfile generated/updated.")

	// Stop containers
	fmt.Println("\n----------------------------------------")
	fmt.Printf("🔎 Checking for Docker containers using port %s...\n", backendPort)
	stopCmd := exec.Command("powershell", "-Command", fmt.Sprintf("docker ps --filter 'publish=%s' --format '{{.ID}}' | ForEach-Object { docker stop $_ >$null 2>&1 } | Out-Null", backendPort))
	stopCmd.Stdout = nil
	stopCmd.Stderr = nil
	_ = stopCmd.Run()
	fmt.Printf("✅ Stopped any Docker container using port %s.\n", backendPort)

	// Remove container
	removeCmd := exec.Command("powershell", "-Command", fmt.Sprintf("docker rm -f %s >$null 2>&1", containerName))
	removeCmd.Stdout = nil
	removeCmd.Stderr = nil
	_ = removeCmd.Run()
	fmt.Printf("🗑️  Removed any container named %s.\n", containerName)

	// Build image
	fmt.Println("\n----------------------------------------")
	fmt.Printf("🔨 Building Docker image '%s'...\n", imageName)
	buildCmd := exec.Command("docker", "build", "-t", imageName, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Println("❌ Docker build failed:", err)
		return
	}
	fmt.Println("✅ Docker image built successfully!")

	// Run container
	fmt.Println("\n----------------------------------------")
	fmt.Println("🚀 Running Docker container...")
	envVars := os.Environ()
	envVars = append(envVars, "MONGO_URL=mongodb://host.docker.internal:27017/chatsAppDocker")
	runCmd := exec.Command("docker", "run", "--name", containerName, "-p", backendPort+":"+backendPort, "-e", "MONGO_URL=mongodb://host.docker.internal:27017/chatsAppDocker", imageName)
	runCmd.Env = envVars
	runCmd.Stdout = nil
	runCmd.Stderr = nil
	if err := runCmd.Start(); err != nil {
		fmt.Println("❌ Docker run failed:", err)
		return
	}
	fmt.Println("✅ Container started!")

	fmt.Println("\n----------------------------------------")
	fmt.Println("📋 Container Details:")
	fmt.Printf("   - Container Name: %s\n", containerName)
	fmt.Printf("   - Backend Port:   %s\n", backendPort)
	fmt.Printf("   - Entry File:     %s\n", entryFile)
	fmt.Printf("   - Database:       %s\n", dbType)
	fmt.Printf("   - Access Link:    http://localhost:%s\n", backendPort)
	fmt.Println("----------------------------------------")
	// Script exits immediately after starting container
}

func isValidPort(port string) bool {
	if forbiddenPorts[port] {
		return false
	}
	p, err := strconv.Atoi(port)
	return err == nil && p > 1024 && p < 65535
}

func isDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	return cmd.Run() == nil
}

func detectProjectType() string {
	if fileExists("package.json") || fileExists("package-lock.json") {
		return "node"
	}
	if fileExists("requirements.txt") {
		return "python"
	}
	return ""
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func generateDockerfile(projectType string, port string, entryFile string) {
	var dockerfile string
	if projectType == "node" {
		dockerfile = fmt.Sprintf(`FROM node:18
WORKDIR /app
COPY . .
RUN npm install
EXPOSE %s
CMD ["node", "%s"]
`, port, entryFile)
	} else if projectType == "python" {
		dockerfile = fmt.Sprintf(`FROM python:3.11
WORKDIR /app
COPY . .
RUN pip install -r requirements.txt
EXPOSE %s
CMD ["python", "%s"]
`, port, entryFile)
	}
	os.WriteFile("Dockerfile", []byte(dockerfile), 0644)
}

// Helper to prompt and validate entry file
func promptEntryFile(projectType string) (string, bool) {
	var entryFile string
	var entryPrompt string
	if projectType == "node" {
		entryPrompt = "📝 Enter your main server JS file (e.g. server.js, index.js, app.js): "
	} else if projectType == "python" {
		entryPrompt = "📝 Enter your main Python file (e.g. app.py, main.py): "
	} else {
		fmt.Println("❌ Project type not detected. Supported: Node.js, Python.")
		return "", false
	}
	fmt.Print(entryPrompt)
	fmt.Scanln(&entryFile)
	if entryFile == "" {
		fmt.Println("❌ Entry file cannot be empty. Please specify your main file.")
		return "", false
	}
	if fileExists(entryFile) {
		fmt.Printf("✅ Found entry file: %s\n", entryFile)
		return entryFile, true
	} else {
		if projectType == "node" {
			fmt.Println("❌ No entry JS file found. Please create server.js, index.js, or app.js, or specify the file name.")
		} else if projectType == "python" {
			fmt.Println("❌ No entry Python file found. Please create app.py, main.py, or specify the file name.")
		}
		return "", false
	}
}

// Show logs for the container
func showContainerLogs(containerName string) {
	if containerName == "" {
		fmt.Println("You did not provide a container name.\nUsage: docify.exe logs container_name")
		return
	}
	fmt.Printf("📜 Showing last 20 lines of logs for '%s'...\n", containerName)
	cmd := exec.Command("docker", "logs", "--tail", "20", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// Delete container by name (parameterized)
func docifyDeleteParam(containerName string) {
	if containerName == "" {
		fmt.Println("You did not provide a container name.\nUsage: docify.exe delete container_name")
		return
	}
	fmt.Printf("🛑 Stopping and removing container '%s'...\n", containerName)
	cmd := exec.Command("docker", "rm", "-f", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	fmt.Println("🗑️ Container removed.")
}

// Show all Docker containers
func showAllContainers() {
	fmt.Println("📋 Listing all Docker containers :")
	cmd := exec.Command("docker", "ps", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
