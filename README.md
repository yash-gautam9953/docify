# 🐳 Docify

![Release](https://img.shields.io/github/v/release/yash-gautam9953/docify?label=Latest%20Release&color=blue)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-green)

**Docify** is a smart **Go-based CLI tool** that auto-generates a production-ready **Docker image** for your project — with **zero configuration**. 🚀

Just run `./docify` in your project folder. It auto-detects everything and builds a Docker image. **No questions asked.**

---

## ✨ Supported Projects

| Type                        | Detected Via                                 | Dockerfile                          |
| --------------------------- | -------------------------------------------- | ----------------------------------- |
| **Node.js (Backend)**       | `package.json` + server code (express, etc.) | `node:18-alpine` + `npm ci`         |
| **React / Vite (Frontend)** | `react-scripts` or `vite` in package.json    | Multi-stage: build → `nginx:alpine` |
| **Next.js (Fullstack)**     | `next` in dependencies or `next.config.*`    | Multi-stage: build → `npm start`    |
| **Python**                  | `requirements.txt`                           | `python:3.11-slim` + `pip install`  |

---

## 📋 Requirements

- [Docker Desktop](https://www.docker.com/products/docker-desktop/) running ✅
- A supported project (see table above)

---

## 📥 Quick Install

### 🪟 Windows

```powershell
Invoke-WebRequest -Uri "https://github.com/yash-gautam9953/docify/releases/latest/download/docify.exe" -OutFile "docify.exe"
```

### 🐧 Linux

```bash
curl -L -o docify https://github.com/yash-gautam9953/docify/releases/latest/download/docify
chmod +x docify
```

---

## 🚀 Usage

```bash
cd my-project
./docify
```

**That's it.** Docify will:

1. ✅ Verify Docker is running
2. 📦 Auto-detect project type (Node backend / React / Next.js / Python)
3. 📄 Auto-detect entry file (backend/Python only)
4. 🌐 Auto-detect port
5. 📝 Generate optimized Dockerfile (multi-stage for frontend)
6. 🔨 Build the Docker image

Then run your image:

```bash
docker run -p 3000:3000 docify-my-project
```

> User input is **only** asked if auto-detection fails (rare).

---

## 🎯 What Gets Auto-Detected

| Detection      | How                                                                        |
| -------------- | -------------------------------------------------------------------------- |
| **Project**    | `package.json` deps → React/Next/Vite/Backend, `requirements.txt` → Python |
| **Entry File** | `package.json` main/start script → `server.js`, `app.js`...                |
| **Port**       | `.env` → code scan → `package.json` scripts → framework defaults           |
| **Image Name** | Your project folder name → `docify-<folder>`                               |

### Port Detection Priority

1. `.env` file → `PORT=3000`
2. Code scan → `app.listen(3000)`, `const port = 3000`, `app.run(port=5000)`
3. `package.json` → `--port 5173` in scripts
4. Framework defaults → Vite: `5173`, React/Next.js: `3000`
5. **If all fail** → asks user (only input ever needed)

---

## 📂 Examples

### React / Vite Frontend

```
my-react-app/
├── src/
├── package.json       ← has "react-scripts" or "vite"
├── public/
└── docify             ← run here
```

Generates a **multi-stage Dockerfile**: builds with Node, serves with Nginx.

### Next.js App

```
my-next-app/
├── pages/ or app/
├── next.config.mjs
├── package.json       ← has "next"
└── docify             ← run here
```

Generates a **multi-stage Dockerfile**: builds with Node, runs with `npm start`.

### Node.js Backend

```
my-api/
├── server.js          ← auto-detected entry
├── package.json
├── .env               ← PORT=4000 (optional)
└── docify             ← run here
```

### Python Backend

```
my-flask-app/
├── app.py             ← auto-detected entry
├── requirements.txt
└── docify             ← run here
```

---

## 🏆 Why Docify?

| Without Docify 😫                | With Docify 🎉          |
| -------------------------------- | ----------------------- |
| Write Dockerfile manually        | Auto-generated          |
| Research multi-stage builds      | Built-in for React/Next |
| Figure out entry file & port     | Auto-detected           |
| Different setup per project type | One command for all     |
| **~15 min per project**          | **~10 seconds**         |

---

## 📦 Releases

### Latest: [v1.0.1](https://github.com/yash-gautam9953/docify/releases/tag/v1.0.1)

| Asset        | Platform        | Size    |
| ------------ | --------------- | ------- |
| `docify`     | Linux (amd64)   | 1.91 MB |
| `docify.exe` | Windows (amd64) | 2.91 MB |

Download from the [Releases page](https://github.com/yash-gautam9953/docify/releases).

---

## 👨‍💻 Author

**Built with ❤️ & 🐳 by [Yash Gautam](https://github.com/yash-gautam9953)**

⭐ **Star this repo if Docify saved your time!**
