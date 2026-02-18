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
		} else if cmd == "info" {
			showProjectContainerInfo()
			return
		} else if cmd == "rebuild" {
			rebuildContainer()
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
					// Node.js: const/let/var port = 3000; (case insensitive)
					lineTrim := strings.TrimSpace(line)
					lowerLine := strings.ToLower(lineTrim)

					// Check for port variable declarations (both uppercase and lowercase)
					portPatterns := []string{
						"const port =", "let port =", "var port =",
						"const PORT =", "let PORT =", "var PORT =",
					}

					for _, pattern := range portPatterns {
						if strings.HasPrefix(lowerLine, pattern) {
							parts := strings.Split(lineTrim, "=")
							if len(parts) >= 2 {
								portValue := strings.TrimSpace(parts[1])
								portValue = strings.Trim(portValue, ";")
								portValue = strings.Trim(portValue, " ")
								if _, err := strconv.Atoi(portValue); err == nil {
									return portValue
								}
							}
							break
						}
					}

					// Also check for app.listen(port, callback) pattern where port is a variable
					if strings.Contains(lowerLine, "app.listen(") || strings.Contains(lowerLine, "server.listen(") {
						// Extract port from listen(port, ...)
						if strings.Contains(line, "listen(") {
							start := strings.Index(line, "listen(") + 7
							if start < len(line) {
								rest := line[start:]
								end := strings.Index(rest, ",")
								if end == -1 {
									end = strings.Index(rest, ")")
								}
								if end > 0 {
									portVar := strings.TrimSpace(rest[:end])
									// If it's a number directly
									if _, err := strconv.Atoi(portVar); err == nil {
										return portVar
									}
									// If it's a variable, try to find its value in same file
									if portValue := findVariableValue(string(data), portVar); portValue != "" {
										return portValue
									}
								}
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

// Helper function to find variable value in code
func findVariableValue(content, varName string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		lowerLine := strings.ToLower(trimmed)

		// Look for variable declarations
		patterns := []string{
			"const " + strings.ToLower(varName) + " =",
			"let " + strings.ToLower(varName) + " =",
			"var " + strings.ToLower(varName) + " =",
		}

		for _, pattern := range patterns {
			if strings.HasPrefix(lowerLine, pattern) {
				parts := strings.Split(trimmed, "=")
				if len(parts) >= 2 {
					value := strings.TrimSpace(parts[1])
					value = strings.Trim(value, ";")
					value = strings.Trim(value, " ")
					if _, err := strconv.Atoi(value); err == nil {
						return value
					}
				}
				break
			}
		}
	}
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
	defaultContainerName := generateContainerName()
	fmt.Println("\n----------------------------------------")
	fmt.Printf("📦 Enter container name (default: %s): ", defaultContainerName)
	var containerName string
	fmt.Scanln(&containerName)
	if containerName == "" {
		containerName = defaultContainerName
		fmt.Printf("🎯 Using default container name: %s\n", containerName)
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

	// Save container info for future reference
	saveContainerInfo(containerName, backendPort, entryFile)

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

// Generate smart container name based on current directory
func generateContainerName() string {
	// Get current directory name
	currentDir, err := os.Getwd()
	if err != nil {
		return "docify-app" // fallback
	}

	// Extract just the folder name
	folderName := strings.ToLower(strings.ReplaceAll(currentDir[strings.LastIndex(currentDir, string(os.PathSeparator))+1:], " ", "-"))

	// Clean up the name (remove special characters)
	folderName = strings.ReplaceAll(folderName, "_", "-")
	folderName = strings.ReplaceAll(folderName, ".", "-")

	// If folder name is too short or generic, use docify-app
	if len(folderName) < 3 || folderName == "src" || folderName == "app" {
		return "docify-app"
	}

	return "docify-" + folderName
}

// Save container info to a local file for tracking
func saveContainerInfo(containerName, port, entryFile string) {
	info := fmt.Sprintf("CONTAINER_NAME=%s\nPORT=%s\nENTRY_FILE=%s\n", containerName, port, entryFile)
	os.WriteFile(".docify-info", []byte(info), 0644)
}

// Load container info from local file
func loadContainerInfo() (string, string, string) {
	if !fileExists(".docify-info") {
		return "", "", ""
	}

	data, err := os.ReadFile(".docify-info")
	if err != nil {
		return "", "", ""
	}

	var containerName, port, entryFile string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "CONTAINER_NAME=") {
			containerName = strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
		} else if strings.HasPrefix(line, "PORT=") {
			port = strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
		} else if strings.HasPrefix(line, "ENTRY_FILE=") {
			entryFile = strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
		}
	}

	return containerName, port, entryFile
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

// Auto-detect and validate entry file
func promptEntryFile(projectType string) (string, bool) {
	// First try to auto-detect the entry file
	autoDetectedFile := autoDetectEntryFile(projectType)
	if autoDetectedFile != "" {
		fmt.Printf("🎯 Auto-detected entry file: %s\n", autoDetectedFile)
		fmt.Print("✅ Use this file? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" || response == "" {
			return autoDetectedFile, true
		}
	}

	// If auto-detection failed or user rejected, ask manually
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

// Auto-detect entry file based on project type and common patterns
func autoDetectEntryFile(projectType string) string {
	if projectType == "node" {
		// Priority order for Node.js files
		nodeFiles := []string{"server.js", "app.js", "index.js", "main.js"}

		// First check package.json for main field
		if fileExists("package.json") {
			data, err := os.ReadFile("package.json")
			if err == nil {
				content := string(data)
				// Look for "main": "filename.js"
				if strings.Contains(content, `"main"`) {
					lines := strings.Split(content, "\n")
					for _, line := range lines {
						if strings.Contains(line, `"main"`) && strings.Contains(line, ":") {
							parts := strings.Split(line, ":")
							if len(parts) > 1 {
								mainFile := strings.TrimSpace(parts[1])
								mainFile = strings.Trim(mainFile, `"`)
								mainFile = strings.Trim(mainFile, ",")
								if fileExists(mainFile) {
									return mainFile
								}
							}
						}
					}
				}
				// Look for "start" script
				if strings.Contains(content, `"start"`) && strings.Contains(content, "node ") {
					lines := strings.Split(content, "\n")
					for _, line := range lines {
						if strings.Contains(line, `"start"`) && strings.Contains(line, "node ") {
							if strings.Contains(line, "node ") {
								parts := strings.Split(line, "node ")
								if len(parts) > 1 {
									scriptFile := strings.TrimSpace(parts[1])
									scriptFile = strings.Trim(scriptFile, `"`)
									scriptFile = strings.Trim(scriptFile, ",")
									if fileExists(scriptFile) {
										return scriptFile
									}
								}
							}
						}
					}
				}
			}
		}

		// Check common Node.js entry files
		for _, file := range nodeFiles {
			if fileExists(file) {
				return file
			}
		}

		// Look for files with server patterns
		files, _ := os.ReadDir(".")
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".js") && !f.IsDir() {
				data, err := os.ReadFile(f.Name())
				if err == nil {
					content := string(data)
					// Check for server patterns
					if strings.Contains(content, "app.listen(") ||
						strings.Contains(content, "server.listen(") ||
						strings.Contains(content, "express()") ||
						strings.Contains(content, "createServer") {
						return f.Name()
					}
				}
			}
		}

	} else if projectType == "python" {
		// Priority order for Python files
		pythonFiles := []string{"app.py", "main.py", "server.py", "run.py", "wsgi.py"}

		// Check common Python entry files
		for _, file := range pythonFiles {
			if fileExists(file) {
				return file
			}
		}

		// Look for files with Flask/Django/FastAPI patterns
		files, _ := os.ReadDir(".")
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".py") && !f.IsDir() {
				data, err := os.ReadFile(f.Name())
				if err == nil {
					content := string(data)
					// Check for web framework patterns
					if strings.Contains(content, "from flask import") ||
						strings.Contains(content, "Flask(__name__)") ||
						strings.Contains(content, "from django") ||
						strings.Contains(content, "from fastapi import") ||
						strings.Contains(content, "FastAPI()") ||
						strings.Contains(content, "app.run(") {
						return f.Name()
					}
				}
			}
		}
	}

	return "" // No suitable file found
}

// Show logs for the container
func showContainerLogs(containerName string) {
	// If no container name provided, try to load from saved info first
	if containerName == "" {
		savedContainer, _, _ := loadContainerInfo()
		if savedContainer != "" {
			containerName = savedContainer
			fmt.Printf("🎯 Using saved container name: %s\n", containerName)
		} else {
			containerName = generateContainerName()
			fmt.Printf("🎯 Using default container name: %s\n", containerName)
		}
	}
	fmt.Printf("📜 Showing last 20 lines of logs for '%s'...\n", containerName)
	cmd := exec.Command("docker", "logs", "--tail", "20", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// Delete container by name (parameterized)
func docifyDeleteParam(containerName string) {
	// If no container name provided, try to load from saved info first
	if containerName == "" {
		savedContainer, _, _ := loadContainerInfo()
		if savedContainer != "" {
			containerName = savedContainer
			fmt.Printf("🎯 Using saved container name: %s\n", containerName)
		} else {
			containerName = generateContainerName()
			fmt.Printf("🎯 Using default container name: %s\n", containerName)
		}
	}
	fmt.Printf("🛑 Stopping and removing container '%s'...\n", containerName)
	cmd := exec.Command("docker", "rm", "-f", containerName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
	fmt.Println("🗑️ Container removed.")

	// Remove the saved info file since container is deleted
	if fileExists(".docify-info") {
		os.Remove(".docify-info")
		fmt.Println("🧹 Cleaned up container info.")
	}
}

// Show all Docker containers
func showAllContainers() {
	fmt.Println("📋 Listing all Docker containers :")
	cmd := exec.Command("docker", "ps", "-a")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}

// Show current project's container info
func showProjectContainerInfo() {
	savedContainer, port, entryFile := loadContainerInfo()

	if savedContainer == "" {
		fmt.Println("❌ No container info found for this project.")
		fmt.Println("💡 Run 'docify' first to create a container.")
		return
	}

	fmt.Println("📋 Current Project Container Info:")
	fmt.Printf("   - Container Name: %s\n", savedContainer)
	fmt.Printf("   - Port:          %s\n", port)
	fmt.Printf("   - Entry File:    %s\n", entryFile)

	// Check if container is running
	checkCmd := exec.Command("docker", "ps", "--filter", "name="+savedContainer, "--format", "{{.Status}}")
	output, err := checkCmd.Output()
	if err == nil && len(output) > 0 {
		fmt.Printf("   - Status:        🟢 Running\n")
		fmt.Printf("   - Access Link:   http://localhost:%s\n", port)
	} else {
		fmt.Printf("   - Status:        🔴 Stopped/Not Found\n")
	}

	fmt.Println("\n💡 Quick Commands:")
	fmt.Println("   - docify logs     (Show logs)")
	fmt.Println("   - docify delete   (Remove container)")
	fmt.Println("   - docify rebuild  (Rebuild with latest code)")
	fmt.Println("   - docify show     (Show all containers)")
}

// Rebuild container with latest code changes
func rebuildContainer() {
	fmt.Println("\n🔄 Rebuilding Container with Latest Code...")
	fmt.Println("==========================================")

	// Load existing container info
	savedContainer, port, entryFile := loadContainerInfo()

	if savedContainer == "" {
		fmt.Println("❌ No existing container found for this project.")
		fmt.Println("💡 Run 'docify' first to create a container, then use 'rebuild'.")
		return
	}

	fmt.Printf("🎯 Found existing container: %s\n", savedContainer)

	// Detect project type
	projectType := detectProjectType()
	if projectType == "" {
		fmt.Println("❌ Project type not detected. Supported: Node.js, Python.")
		return
	}

	fmt.Printf("📦 Project Type: %s\n", strings.Title(projectType))

	// Step 1: Stop and remove existing container
	fmt.Println("\n📛 Step 1: Stopping existing container...")
	stopCmd := exec.Command("docker", "stop", savedContainer)
	stopCmd.Stdout = nil
	stopCmd.Stderr = nil
	_ = stopCmd.Run()

	removeCmd := exec.Command("docker", "rm", "-f", savedContainer)
	removeCmd.Stdout = nil
	removeCmd.Stderr = nil
	_ = removeCmd.Run()
	fmt.Printf("✅ Removed container: %s\n", savedContainer)

	// Step 2: Remove old Docker image
	fmt.Println("\n🗑️  Step 2: Removing old Docker image...")
	imageName := "docify-app"
	removeImageCmd := exec.Command("docker", "rmi", "-f", imageName)
	removeImageCmd.Stdout = nil
	removeImageCmd.Stderr = nil
	_ = removeImageCmd.Run()
	fmt.Printf("✅ Removed image: %s\n", imageName)

	// Step 3: Regenerate Dockerfile (in case entry file changed)
	fmt.Println("\n📝 Step 3: Updating Dockerfile...")
	if fileExists("Dockerfile") {
		os.Remove("Dockerfile")
	}
	generateDockerfile(projectType, port, entryFile)
	fmt.Println("✅ Dockerfile updated with latest configuration")

	// Step 4: Build new image
	fmt.Println("\n🔨 Step 4: Building new Docker image...")
	buildCmd := exec.Command("docker", "build", "-t", imageName, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fmt.Println("❌ Docker build failed:", err)
		return
	}
	fmt.Println("✅ New Docker image built successfully!")

	// Step 5: Run new container
	fmt.Println("\n🚀 Step 5: Starting new container...")

	// Check database type for MongoDB setup
	dbType := detectDBUsage()
	var runCmd *exec.Cmd

	if dbType == "mongodb" {
		runCmd = exec.Command("docker", "run", "--name", savedContainer, "-p", port+":"+port, "-e", "MONGO_URL=mongodb://host.docker.internal:27017/chatsAppDocker", "-d", imageName)
	} else {
		runCmd = exec.Command("docker", "run", "--name", savedContainer, "-p", port+":"+port, "-d", imageName)
	}

	runCmd.Stdout = nil
	runCmd.Stderr = nil
	if err := runCmd.Run(); err != nil {
		fmt.Println("❌ Container start failed:", err)
		return
	}

	fmt.Println("✅ New container started successfully!")

	// Final summary
	fmt.Println("\n🎉 Rebuild Complete!")
	fmt.Println("==========================================")
	fmt.Printf("   - Container:    %s\n", savedContainer)
	fmt.Printf("   - Port:         %s\n", port)
	fmt.Printf("   - Entry File:   %s\n", entryFile)
	fmt.Printf("   - Database:     %s\n", dbType)
	fmt.Printf("   - Access Link:  http://localhost:%s\n", port)
	fmt.Println("\n💡 Your latest code changes are now running!")
}
