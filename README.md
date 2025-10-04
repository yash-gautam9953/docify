# 🐳 Docify

**Docify** is a lightweight **Go-based CLI tool** that automatically Dockerizes your **Node.js** or **Python** projects.  
It detects ports, databases, generates a `Dockerfile`, build an image, and runs your app in a container — with almost **zero setup**. 🚀  

---

## ✨ Features
- 🔍 **Auto-detects backend port** (from `.env` or `.js` files)  
- 🐳 **Generates Dockerfile automatically** (Node.js / Python supported)  
- 📦 **Builds & runs Docker containers** with a single command  
- 🛑 **Stops & removes conflicting containers** (same port or name)  
- 📜 **View logs of running containers** (`docify logs <container_name>`)  
- ❌ **Delete containers easily** (`docify delete <container_name>`)  
- 📋 **List all containers** (`docify show`)  
- 🗃 **MongoDB support** — injects `MONGO_URL` automatically  
- ⚡ Works seamlessly for **Node.js** & **Python** backends  

---

## 📋 Requirements
- [Docker Desktop](https://www.docker.com/products/docker-desktop/) running ✅  
- Node.js or Python project with:  
  - `package.json` (for Node.js)  
  - `requirements.txt` (for Python)  
- Optional: `.env` file with `PORT=XXXX`  

---

## 📥 Download

You don’t need to clone the repo just for the `.exe`. Simply download it using PowerShell inside your project folder and run it:


    Invoke-WebRequest -Uri "https://github.com/yash-gautam9953/docify/raw/main/docify.exe" -OutFile "docify.exe"



## 📂 Project Structure Example

Keep your project like this for smooth usage:

    myapp/
      │── server.js
      │── models/
      │── routes/
      │── package.json
      │── docify.exe 👈 keep the exe file here


Now just open a terminal in `myapp/` and run:

    ./docify.exe


## 🚀 Usage

     ./docify.exe
    
    ./docify.exe logs <container_name>
    
    ./docify.exe show

    ./docify.exe delete <container_name>



## 🗄 MongoDB Setup

Inside your Node.js project, always connect like this:

    const mongoUrl = process.env.MONGO_URL || "mongodb://127.0.0.1:27017/YOUR-DB-NAME";
    await mongoose.connect(mongoUrl);

Docify will inject the correct MONGO_URL into your container automatically.

## 👨‍💻 Author

### Built with ❤️ & 🐳 by Yash Gautam


