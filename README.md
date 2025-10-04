# 🐳 Docify

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

- 🚀 **Zero-config setup**: `./docify.exe`
- 🔄 **Quick rebuild**: `./docify.exe rebuild`
- 📜 **Smart logs**: `./docify.exe logs` (auto-detects your container)
- ❌ **Easy cleanup**: `./docify.exe delete` (auto-detects your container)
- 📋 **Project info**: `./docify.exe info`
- 🗂️ **All containers**: `./docify.exe show`

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

```powershell
# PowerShell
Invoke-WebRequest -Uri "https://github.com/yash-gautam9953/docifyCode/releases/latest/download/docify.exe" -OutFile "docify.exe"
```

```bash
# Git Bash / WSL
curl -L -o docify.exe https://github.com/yash-gautam9953/docifyCode/releases/latest/download/docify.exe
```

---

## 🚀 Usage

### **Basic Usage**

```bash
# Navigate to your project folder
cd my-node-app

# Download docify.exe (one time)
# ... (use download command above)

# Run your app in Docker (that's it!)
./docify.exe
```

### **Advanced Commands**

```bash
# Rebuild after code changes
./docify.exe rebuild

# View logs (auto-detects your container)
./docify.exe logs

# Delete container (auto-detects your container)
./docify.exe delete

# Show current project's container info
./docify.exe info

# List all Docker containers
./docify.exe show
```

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
└── docify.exe        # Place here
```

### **Python Project**

```
my-api/
├── app.py            # Auto-detected entry file
├── requirements.txt  # Project type detection
├── .env             # PORT=5000 (optional)
├── models/
└── docify.exe       # Place here
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
./docify.exe

# Code changes? Just rebuild
./docify.exe rebuild

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

| Command                | Function                       | Example            |
| ---------------------- | ------------------------------ | ------------------ |
| `./docify.exe`         | Complete Docker setup + run    | Initial deployment |
| `./docify.exe rebuild` | Rebuild with latest code       | After code changes |
| `./docify.exe logs`    | Show container logs            | Debug issues       |
| `./docify.exe delete`  | Remove container + cleanup     | Clean shutdown     |
| `./docify.exe info`    | Show project container details | Check status       |
| `./docify.exe show`    | List all Docker containers     | System overview    |

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

----


## 👨‍💻 Author

**Built with ❤️ & 🐳 by [Yash Gautam](https://github.com/yash-gautam9953)**

⭐ **Star this repo if Docify saved your time!**

