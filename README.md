# 🐳 Docify

![Release](https://img.shields.io/github/v/release/yash-gautam9953/docify?label=Latest%20Release&color=blue)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-green)

**Docify** is an intelligent **Go-based CLI tool** that automatically Dockerizes your **Node.js** or **Python** projects.  
It intelligently detects project type, entry files, ports, databases, and manages the complete Docker lifecycle — with **90% automation**. 🚀

**Stop wasting time on Docker configuration. Focus on coding, let Docify handle containerization!**

---

## ✨ Features

### 🧠 **Smart Auto-Detection**

- 🎯 **Auto-detects project type** (Node.js via `package.json`, Python via `requirements.txt`)
- 📝 **Auto-detects entry files** (`server.js`, `app.js`, `index.js`, `app.py`, `main.py`)
- 🔍 **Auto-detects backend port** (from `.env`, code patterns, `package.json`)
- 🗃 **Auto-detects database usage** (MongoDB patterns in code)
- 📦 **Smart container naming** (based on project folder)

### 🐳 **Docker Management**

- 🏗️ **Generates optimized Dockerfile** automatically
- 📦 **Builds & runs containers** with single command
- 🛑 **Handles port conflicts** intelligently
- 🔄 **Rebuild containers** with latest code changes
- 📊 **Project-specific container tracking**

### 🎛️ **Developer Commands**

- 🚀 **Zero-config setup**: `./docify` (Linux) / `./docify.exe` (Windows)
- 🔄 **Quick rebuild**: `./docify rebuild`
- 📜 **Smart logs**: `./docify logs` (auto-detects your container)
- ❌ **Easy cleanup**: `./docify delete` (auto-detects your container)
- 📋 **Project info**: `./docify info`
- 🗂️ **All containers**: `./docify show`

---

## 📋 Requirements

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) running ✅
- Node.js or Python project with:
  - `package.json` (for Node.js)
  - `requirements.txt` (for Python)
- Optional: `.env` file with `PORT=XXXX`

---

## 📥 Quick Install

Download directly into your project folder:

### 🪟 **Windows**

```powershell
# PowerShell
Invoke-WebRequest -Uri "https://github.com/yash-gautam9953/docify/releases/latest/download/docify.exe" -OutFile "docify.exe"
```

```bash
# Git Bash / WSL
curl -L -o docify.exe https://github.com/yash-gautam9953/docify/releases/latest/download/docify.exe
```

### 🐧 **Linux**

```bash
# Download the Linux binary
curl -L -o docify https://github.com/yash-gautam9953/docify/releases/latest/download/docify

# Make it executable
chmod +x docify
```

---

## 🚀 Usage

### **Basic Usage**

```bash
# Navigate to your project folder
cd my-node-app

# Download docify (one time) — see Quick Install above

# Run your app in Docker (that's it!)
./docify          # Linux
./docify.exe      # Windows
```

### **Advanced Commands**

```bash
# Rebuild after code changes
./docify rebuild

# View logs (auto-detects your container)
./docify logs

# Delete container (auto-detects your container)
./docify delete

# Show current project's container info
./docify info

# List all Docker containers
./docify show
```

> 💡 On Windows, replace `./docify` with `./docify.exe` in all commands above.

---

## 📂 Project Structure Example

### **Node.js Project**

```
my-blog-app/
├── server.js          # Auto-detected entry file
├── package.json       # Project type detection
├── .env              # PORT=3000 (optional)
├── routes/
├── models/
└── docify            # Linux binary (or docify.exe on Windows)
```

### **Python Project**

```
my-api/
├── app.py            # Auto-detected entry file
├── requirements.txt  # Project type detection
├── .env             # PORT=5000 (optional)
├── models/
└── docify           # Linux binary (or docify.exe on Windows)
```

---

## 🔄 Development Workflow

### **Traditional Docker Workflow** ❌

```bash
# Write Dockerfile manually
# Build image
docker build -t my-app .
# Handle port conflicts
docker stop $(docker ps -q --filter publish=3000)
# Run container
docker run -p 3000:3000 --name my-container my-app
# Code changes? Repeat everything...
```

### **Docify Workflow** ✅

```bash
# Initial setup
./docify

# Code changes? Just rebuild
./docify rebuild

# That's it! 🎉
```

---

## 🗄️ Database Support

### **MongoDB Auto-Configuration**

Docify automatically detects MongoDB usage and configures connection:

```javascript
// Your Node.js code
const mongoUrl = process.env.MONGO_URL || "mongodb://127.0.0.1:27017/mydb";
await mongoose.connect(mongoUrl);
```

```python
# Your Python code
import os
mongo_url = os.getenv('MONGO_URL', 'mongodb://127.0.0.1:27017/mydb')
client = MongoClient(mongo_url)
```

**Docify automatically injects**: `MONGO_URL=mongodb://host.docker.internal:27017/chatsAppDocker`

---

## 🎯 Smart Features

### **Auto-Detection Examples**

#### **Port Detection**

```javascript
// Detects from multiple patterns:
const port = 3000;                    ✅
const PORT = process.env.PORT || 5000; ✅
app.listen(8080, () => {});           ✅
```

#### **Entry File Detection**

```bash
# Priority order:
1. package.json "main" field          ✅
2. package.json "start" script        ✅
3. server.js, app.js, index.js        ✅
4. Files with app.listen() patterns   ✅
```

#### **Container Naming**

```bash
# Smart naming based on folder:
/my-blog-app     → docify-my-blog-app
/chat-server     → docify-chat-server
/generic-folder  → docify-app (fallback)
```

---

## 🔧 Command Reference

| Command            | Function                       | Example            |
| ------------------ | ------------------------------ | ------------------ |
| `./docify`         | Complete Docker setup + run    | Initial deployment |
| `./docify rebuild` | Rebuild with latest code       | After code changes |
| `./docify logs`    | Show container logs            | Debug issues       |
| `./docify delete`  | Remove container + cleanup     | Clean shutdown     |
| `./docify info`    | Show project container details | Check status       |
| `./docify show`    | List all Docker containers     | System overview    |

> 💡 On Windows, use `docify.exe` instead of `docify`.

---

## 💡 Use Cases

- 🏗️ **Rapid Prototyping**: Get ideas running in containers instantly
- 👥 **Team Collaboration**: Share consistent Docker environments
- 🚀 **Client Demos**: Quick containerized deployments
- 🧪 **Testing**: Isolated environment testing
- 📚 **Learning**: Docker without Docker complexity

---

## 🏆 Why Docify?

### **Before Docify** 😫

- Manual Dockerfile writing
- Port conflict management
- Docker commands memorization
- Environment setup complexity
- 30+ minutes per project

### **After Docify** 🎉

- Zero Docker knowledge required
- One command deployment
- Automatic conflict resolution
- Smart environment detection
- 30 seconds per project

---

## � Releases

### Latest: [v1.0.1](https://github.com/yash-gautam9953/docify/releases/tag/v1.0.1)

| Asset        | Platform        | Size    |
| ------------ | --------------- | ------- |
| `docify`     | Linux (amd64)   | 1.91 MB |
| `docify.exe` | Windows (amd64) | 2.91 MB |

Download from the [Releases page](https://github.com/yash-gautam9953/docify/releases).

---

## �👨‍💻 Author

**Built with ❤️ & 🐳 by [Yash Gautam](https://github.com/yash-gautam9953)**

⭐ **Star this repo if Docify saved your time!**
