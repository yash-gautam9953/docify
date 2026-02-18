package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// ─── Project Types ──────────────────────────────────────────────────
const (
	TypeNodeBackend  = "node-backend"
	TypeNodeFrontend = "node-frontend"
	TypeNextJS       = "nextjs"
	TypePython       = "python"
)

var forbiddenPorts = map[string]bool{
	"20": true, "21": true, "22": true, "23": true, "25": true,
	"53": true, "80": true, "443": true, "3306": true, "5432": true, "27017": true,
}

// ─── Main ───────────────────────────────────────────────────────────
func main() {
	fmt.Println("\n🚀 Docify — Auto Docker Image Builder")
	fmt.Println("========================================")

	// 1. Docker check
	if !isDockerRunning() {
		fatal("Docker is not running. Please start Docker Desktop and try again.")
	}
	info("Docker is running.")

	// 2. Detect project type
	projectType := detectProjectType()
	if projectType == "" {
		fatal("Could not detect project type.\n   Supported: Node.js, React, Next.js, Vite, Python")
	}
	info("Detected project: " + friendlyName(projectType))

	// 3. Auto-detect entry file (not needed for frontend/nextjs builds)
	entryFile := ""
	if projectType == TypeNodeBackend || projectType == TypePython {
		entryFile = autoDetectEntryFile(projectType)
		if entryFile == "" {
			entryFile = askUser("📝 Could not auto-detect entry file. Enter it manually")
			if entryFile == "" || !fileExists(entryFile) {
				fatal("Invalid or missing entry file.")
			}
		}
		info("Entry file: " + entryFile)
	}

	// 4. Auto-detect port
	port := detectPort(projectType)
	if port == "" {
		port = askUser("🔢 Could not detect port. Enter the port your app runs on")
		if !isValidPort(port) {
			fatal("Invalid or forbidden port.")
		}
	}
	info("Port: " + port)

	// 5. Image name from folder
	imageName := generateImageName()
	info("Image name: " + imageName)

	// 6. Generate Dockerfile
	fmt.Println("\n----------------------------------------")
	os.Remove("Dockerfile") // remove old if exists
	generateDockerfile(projectType, port, entryFile)
	info("Dockerfile generated.")

	// 7. Build
	fmt.Printf("🔨 Building Docker image '%s'...\n", imageName)
	buildCmd := exec.Command("docker", "build", "-t", imageName, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		fatal("Docker build failed: " + err.Error())
	}

	// Done
	fmt.Println("\n========================================")
	fmt.Println("✅ Docker image built successfully!")
	fmt.Println("========================================")
	fmt.Printf("   Image:   %s\n", imageName)
	fmt.Printf("   Type:    %s\n", friendlyName(projectType))
	if entryFile != "" {
		fmt.Printf("   Entry:   %s\n", entryFile)
	}
	fmt.Printf("   Port:    %s\n", port)
	fmt.Println("----------------------------------------")
	fmt.Println("▶️  Run it with:")
	fmt.Printf("   docker run -p %s:%s %s\n", port, port, imageName)
	fmt.Println("========================================")
}

// ─── Helpers ────────────────────────────────────────────────────────
func fatal(msg string) {
	fmt.Println("❌ " + msg)
	os.Exit(1)
}

func info(msg string) {
	fmt.Println("✅ " + msg)
}

func askUser(prompt string) string {
	fmt.Print(prompt + ": ")
	var input string
	fmt.Scanln(&input)
	return strings.TrimSpace(input)
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func readFile(name string) string {
	data, err := os.ReadFile(name)
	if err != nil {
		return ""
	}
	return string(data)
}

func isDockerRunning() bool {
	cmd := exec.Command("docker", "info")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run() == nil
}

func isValidPort(port string) bool {
	if forbiddenPorts[port] {
		return false
	}
	p, err := strconv.Atoi(port)
	return err == nil && p > 1024 && p < 65535
}

func friendlyName(t string) string {
	names := map[string]string{
		TypeNodeBackend:  "Node.js (Backend)",
		TypeNodeFrontend: "Node.js (Frontend — React/Vite)",
		TypeNextJS:       "Next.js (Fullstack)",
		TypePython:       "Python",
	}
	if n, ok := names[t]; ok {
		return n
	}
	return t
}

// ─── Detection: Project Type ────────────────────────────────────────
func detectProjectType() string {
	// Python first (no package.json)
	if !fileExists("package.json") && fileExists("requirements.txt") {
		return TypePython
	}

	if !fileExists("package.json") {
		return ""
	}

	pkg := readFile("package.json")

	// Next.js — has "next" in dependencies or next.config file
	if strings.Contains(pkg, `"next"`) || fileExists("next.config.js") || fileExists("next.config.mjs") || fileExists("next.config.ts") {
		return TypeNextJS
	}

	// Frontend (React / Vite / CRA) — has react-scripts or vite, no server patterns
	if strings.Contains(pkg, `"react-scripts"`) ||
		strings.Contains(pkg, `"vite"`) ||
		strings.Contains(pkg, `"@vitejs/`) ||
		(strings.Contains(pkg, `"react"`) && !hasServerCode()) {
		return TypeNodeFrontend
	}

	// Node.js backend
	return TypeNodeBackend
}

// Check if project has typical backend server code
func hasServerCode() bool {
	patterns := []string{"express()", "app.listen(", "server.listen(", "createServer", "fastify(", "koa()", "hapi.server("}
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".js") || strings.HasSuffix(f.Name(), ".ts") {
			content := readFile(f.Name())
			for _, p := range patterns {
				if strings.Contains(content, p) {
					return true
				}
			}
		}
	}
	return false
}

// ─── Detection: Port ────────────────────────────────────────────────
func detectPort(projectType string) string {
	// Frontend defaults
	switch projectType {
	case TypeNodeFrontend:
		// Vite = 5173, CRA = 3000
		pkg := readFile("package.json")
		if strings.Contains(pkg, `"vite"`) {
			return "5173"
		}
		return "3000"
	case TypeNextJS:
		return "3000"
	}

	// Backend: try .env → code scan → package.json
	if port := portFromEnv(); port != "" {
		return port
	}
	if port := portFromCode(); port != "" {
		return port
	}
	if port := portFromPackageJSON(); port != "" {
		return port
	}
	return ""
}

func portFromEnv() string {
	content := readFile(".env")
	if content == "" {
		return ""
	}
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "PORT=") {
			port := strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
			if isValidPort(port) {
				return port
			}
		}
	}
	return ""
}

func portFromCode() string {
	files, _ := os.ReadDir(".")
	for _, f := range files {
		ext := filepath.Ext(f.Name())
		if ext != ".js" && ext != ".ts" && ext != ".py" {
			continue
		}
		content := readFile(f.Name())
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			// listen(3000) or listen(PORT, ...)
			if strings.Contains(line, "listen(") {
				if port := extractPortFromListen(line, content); port != "" {
					return port
				}
			}
			// Port variable: const port = 3000
			if port := extractPortVariable(line); port != "" {
				return port
			}
			// Python: app.run(port=5000)
			if strings.Contains(line, "port=") {
				if port := extractPythonPort(line); port != "" {
					return port
				}
			}
		}
	}
	return ""
}

func portFromPackageJSON() string {
	content := readFile("package.json")
	// Look for --port in scripts
	for _, line := range strings.Split(content, "\n") {
		if strings.Contains(line, "--port") {
			parts := strings.Split(line, "--port")
			if len(parts) > 1 {
				rest := strings.TrimLeft(parts[1], " =")
				port := strings.Fields(rest)[0]
				port = strings.Trim(port, `"',`)
				if isValidPort(port) {
					return port
				}
			}
		}
	}
	return ""
}

func extractPortFromListen(line, fullContent string) string {
	idx := strings.Index(line, "listen(")
	if idx == -1 {
		return ""
	}
	rest := line[idx+7:]
	end := strings.IndexAny(rest, ",)")
	if end <= 0 {
		return ""
	}
	val := strings.TrimSpace(rest[:end])
	if p, err := strconv.Atoi(val); err == nil && p > 1024 && p < 65535 {
		return val
	}
	// Maybe it's a variable name — resolve it
	return findVarValue(fullContent, val)
}

func extractPortVariable(line string) string {
	lower := strings.ToLower(strings.TrimSpace(line))
	prefixes := []string{"const port =", "let port =", "var port =", "port ="}
	for _, prefix := range prefixes {
		if strings.HasPrefix(lower, prefix) {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) < 2 {
				continue
			}
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, "; ")
			if isValidPort(val) {
				return val
			}
		}
	}
	return ""
}

func extractPythonPort(line string) string {
	idx := strings.Index(line, "port=")
	if idx == -1 {
		return ""
	}
	rest := line[idx+5:]
	end := strings.IndexAny(rest, ",) ")
	if end == -1 {
		end = len(rest)
	}
	val := strings.TrimSpace(rest[:end])
	if isValidPort(val) {
		return val
	}
	return ""
}

func findVarValue(content, varName string) string {
	varName = strings.TrimSpace(varName)
	if varName == "" {
		return ""
	}
	lowerVar := strings.ToLower(varName)
	for _, line := range strings.Split(content, "\n") {
		lower := strings.ToLower(strings.TrimSpace(line))
		for _, kw := range []string{"const ", "let ", "var "} {
			if strings.HasPrefix(lower, kw+lowerVar+" =") || strings.HasPrefix(lower, kw+lowerVar+"=") {
				parts := strings.SplitN(line, "=", 2)
				if len(parts) > 1 {
					val := strings.Trim(strings.TrimSpace(parts[1]), "; ")
					if isValidPort(val) {
						return val
					}
				}
			}
		}
	}
	return ""
}

// ─── Detection: Entry File ─────────────────────────────────────────
func autoDetectEntryFile(projectType string) string {
	if projectType == TypeNodeBackend {
		return detectNodeEntry()
	}
	if projectType == TypePython {
		return detectPythonEntry()
	}
	return "" // frontend/nextjs don't need entry file
}

func detectNodeEntry() string {
	// 1. package.json "main" field
	if fileExists("package.json") {
		pkg := readFile("package.json")
		if f := extractJSONValue(pkg, "main"); f != "" && fileExists(f) {
			return f
		}
		// 2. "start" script → "node server.js"
		if strings.Contains(pkg, `"start"`) {
			for _, line := range strings.Split(pkg, "\n") {
				if strings.Contains(line, `"start"`) && strings.Contains(line, "node ") {
					parts := strings.Split(line, "node ")
					if len(parts) > 1 {
						f := strings.Trim(strings.TrimSpace(parts[1]), `",`)
						if fileExists(f) {
							return f
						}
					}
				}
			}
		}
	}

	// 3. Common entry files
	for _, f := range []string{"server.js", "app.js", "index.js", "main.js", "server.ts", "app.ts", "index.ts"} {
		if fileExists(f) {
			return f
		}
	}

	// 4. Scan for server patterns in .js/.ts files
	files, _ := os.ReadDir(".")
	for _, f := range files {
		name := f.Name()
		if (strings.HasSuffix(name, ".js") || strings.HasSuffix(name, ".ts")) && !f.IsDir() {
			content := readFile(name)
			if strings.Contains(content, "app.listen(") || strings.Contains(content, "server.listen(") ||
				strings.Contains(content, "express()") || strings.Contains(content, "createServer") {
				return name
			}
		}
	}
	return ""
}

func detectPythonEntry() string {
	// 1. Common entry files
	for _, f := range []string{"app.py", "main.py", "server.py", "run.py", "wsgi.py", "manage.py"} {
		if fileExists(f) {
			return f
		}
	}

	// 2. Scan for framework patterns
	files, _ := os.ReadDir(".")
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".py") && !f.IsDir() {
			content := readFile(f.Name())
			if strings.Contains(content, "Flask(") || strings.Contains(content, "FastAPI(") ||
				strings.Contains(content, "from django") || strings.Contains(content, "app.run(") ||
				strings.Contains(content, "uvicorn.run(") {
				return f.Name()
			}
		}
	}
	return ""
}

// Simple JSON value extractor (no external deps)
func extractJSONValue(json, key string) string {
	search := `"` + key + `"`
	idx := strings.Index(json, search)
	if idx == -1 {
		return ""
	}
	rest := json[idx+len(search):]
	rest = strings.TrimLeft(rest, ": ")
	if len(rest) == 0 || rest[0] != '"' {
		return ""
	}
	rest = rest[1:]
	end := strings.Index(rest, `"`)
	if end == -1 {
		return ""
	}
	return rest[:end]
}

// ─── Image Name ─────────────────────────────────────────────────────
func generateImageName() string {
	dir, err := os.Getwd()
	if err != nil {
		return "docify-app"
	}
	name := filepath.Base(dir)
	name = strings.ToLower(name)
	name = strings.NewReplacer(" ", "-", "_", "-", ".", "-").Replace(name)
	if len(name) < 3 || name == "src" || name == "app" {
		return "docify-app"
	}
	return "docify-" + name
}

// ─── Dockerfile Generation ──────────────────────────────────────────
func generateDockerfile(projectType, port, entryFile string) {
	var df string

	switch projectType {
	case TypeNodeBackend:
		df = fmt.Sprintf(`FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE %s
CMD ["node", "%s"]
`, port, entryFile)

	case TypeNodeFrontend:
		df = fmt.Sprintf(`FROM node:18-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
COPY --from=build /app/build /usr/share/nginx/html
EXPOSE %s
CMD ["nginx", "-g", "daemon off;"]
`, port)

	case TypeNextJS:
		df = fmt.Sprintf(`FROM node:18-alpine AS build
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=build /app/.next ./.next
COPY --from=build /app/public ./public
COPY --from=build /app/package*.json ./
COPY --from=build /app/node_modules ./node_modules
EXPOSE %s
CMD ["npm", "start"]
`, port)

	case TypePython:
		df = fmt.Sprintf(`FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE %s
CMD ["python", "%s"]
`, port, entryFile)
	}

	os.WriteFile("Dockerfile", []byte(df), 0644)
}
