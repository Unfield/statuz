# 🟢 Statuz

> **Self‑hosted status and incident management for everything that matters.**

Statuz is a lightweight, open‑source platform for monitoring your services,
tracking incidents, and keeping your users informed — all in real‑time,
with a simple, efficient Go backend and a clean modern UI.

---

## ✨ Features

**Statuz** is built with performance, simplicity, and transparency in mind.

### 🚀 Core Features
- Modern **Go backend** using lightweight concurrency
- Multiple check types — HTTP(S), TCP, and more coming soon
- Real‑time updates via channels and event broadcasting
- Persistent incident tracking (automatic and manual)
- Planned maintenance scheduling
- Simple configuration (YAML or environment variables)

### 🧠 Architecture Highlights
- **Goroutine‑based monitors** for efficient parallel checks
- **Channel‑driven events** for predictable and decoupled updates
- **API‑first design** — everything accessible via REST endpoints
- **Extensible** — easily add custom checks or notification integrations

### 💬 Notifications & Integrations
- User subscription system for service updates
- Notification channels (planned): Email, Webhooks, Slack, Discord
- Custom webhook support for external tools

### 🧭 Frontend & UI
- Minimal and responsive web dashboard
- Live visual status indicators and history
- Built with either **Svelte** or **htmx** (TBD)
- Authentication for administrative view

### 🐳 Deployment
- Single static Go binary — no dependency mess
- Official Docker and Compose setup planned
- Portable configuration for any environment

---

## 🧰 Tech Stack
- **Language:** Go
- **Frontend (planned):** Svelte or htmx
- **Database:** SQLite or Postgres (through Interface based adapters)
- **Architecture:** Concurrent monitors, channel‑based messaging, minimal REST API
- **License:** MIT

---

## ❤️ Philosophy
Statuz was built out of a need for a simple, fast, and self‑hosted alternative to heavy JS‑based monitoring dashboards.
It aims to be **clear, reliable, and hackable** — a tool that just runs and works the way you expect.

---

## 📡 Project Links
- Website: [https://statuz.sh](https://statuz.sh) *(coming soon)*
- GitHub: [github.com/statuzlabs/statuz](https://github.com/statuzlabs/statuz)
- License: MIT

---

> _Statuz is under active development — ideas, feedback, and contributions are always welcome._
