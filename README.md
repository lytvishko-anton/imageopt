# 🚀 Distributed Image Optimizer

A high-performance image optimization pipeline using **Symfony**, **Go**, and **RabbitMQ**.

## 🛠 Tech Stack
- **PHP 8.4 (Symfony)**: Handles file uploads and dispatches tasks.
- **Go**: High-speed worker that converts images to WebP.
- **RabbitMQ**: The message broker connecting the two.
- **Docker**: Fully containerized environment.

## ⚡ Quick Start
1. `git clone <your-repo-url>`
2. `docker-compose up -d --build`
3. Open `http://localhost:8000/index.html`