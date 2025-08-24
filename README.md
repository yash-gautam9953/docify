# 🐳 Docify

Docify is a lightweight **Go-based CLI tool** that automatically Dockerizes your **Node.js** or **Python** projects. It detects ports, databases, generates a Dockerfile, builds an image, and runs your app in a container — with almost zero setup.

---

## ✨ Features
- 🔍 Auto-detects backend port (from `.env` or `.js` files)  
- 🐳 Generates a Dockerfile automatically (Node.js / Python)  
- 📦 Builds & runs Docker containers with a single command  
- 🛑 Stops and removes conflicting containers (same port or name)  
- 🗑 Cleans up containers gracefully on `Ctrl + C`  
- 🗃 MongoDB support — injects `MONGO_URL` automatically  
- ⚡ Works for both Node.js and Python backends  

---

## 📋 Requirements
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) running  
- Node.js or Python project with:  
  - `package.json` (Node.js)  
  - `requirements.txt` (Python)  
- Optional: `.env` file with `PORT=XXXX`  

---

## 📂 Project Structure Example

When using Docify, keep it like this:

myapp/
│── server.js
│── models/
│── routes/
│── package.json
│── docify.exe 👈 keep the exe file here


Now just open a terminal in `myapp/` and run:

./docify.exe


## 🚀 Usage

./docify



🗄 MongoDB Setup

Inside your Node.js project, always connect like this:

const mongoUrl = process.env.MONGO_URL || "mongodb://127.0.0.1:27017/YOUR-DB-NAME";
await mongoose.connect(mongoUrl);

