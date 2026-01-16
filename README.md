
# WasaPhoto. Social Network – Full-Stack Web Application

This repository contains a **full-stack social network application** showcasing modern web development practices, clean architecture, and scalable design. The platform allows users to register, create profiles, post content, comment, follow others, and interact through a real-time feed with notifications. It demonstrates practical skills in backend development, API design, frontend engineering, and end-to-end software delivery.

---

## Table of Contents

* [Project Overview](#project-overview)
* [Key Features](#key-features)
* [Architecture](#architecture)
* [Technologies](#technologies)
* [Installation & Setup](#installation--setup)
* [Usage](#usage)
* [Testing](#testing)
* [Contributing](#contributing)
* [License](#license)

---

## Project Overview

The social network allows users to:

* Register and authenticate securely
* Manage personal profiles
* Create posts, comment, and like content
* Follow other users and view an interactive feed
* Receive real-time notifications about activity

The project follows **clean architecture principles**, separating concerns between executables, services, and frontend components. This modularity ensures maintainable, scalable, and extensible code.

---

## Key Features

* **User Management:** Registration, login/logout, secure authentication, profile editing
* **Content Interaction:** Posting, commenting, liking, following users
* **Interactive Feed:** Real-time updates and notifications
* **API-First Design:** Well-documented backend API for frontend integration or third-party usage
* **Modular Architecture:** Clean separation of concerns, easy to extend and maintain
* **Responsive UI:** Modern and interactive frontend using Vue.js and Bootstrap

---

## Architecture

* **Executables (`cmd/`)**: Entry points and daemon programs (e.g., web API server, health checks)
* **Services (`service/`)**: Core business logic and API implementations, organized into reusable packages
* **Web Frontend (`webui/`)**: Vue.js SPA using Bootstrap for responsive UI and Feather Icons for scalable graphics
* **Documentation (`doc/`)**: OpenAPI-compliant specifications for backend API
* **Dependency Management:**

  * Go dependencies vendored (`vendor/`) for reproducible builds
  * Node.js dependencies included in `webui/node_modules` for frontend

---

## Technologies

* **Backend:** Go (Golang), modular service packages, unit testing
* **Frontend:** Vue.js, Bootstrap, Feather Icons
* **API & Documentation:** OpenAPI / YAML
* **Build & Development Tools:** Node.js, NPM, Vite, Docker (optional for containerized frontend development)

---

## Installation & Setup

### Prerequisites

* Go >= 1.20
* Node.js & NPM
* Docker (optional, for frontend container)

### Setup Backend

```bash
go build ./cmd/webapi/
./cmd/webapi/webapi
```

### Setup Frontend

```bash
./open-npm.sh
npm install
npm run dev       # For development
npm run build-prod  # For production
```

---

## Usage

* Access the application via the web UI (localhost)
* Use the API endpoints for integration or testing
* Create accounts, post content, follow other users, and interact with the feed

---

## Testing

* Unit tests are included for backend services
* Utilities such as `service/globaltime` facilitate reliable and reproducible tests

```bash
go test ./...
```

---

## Contributing

Contributions are welcome. Suggested workflow:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature-name`)
3. Implement changes with clean code and documentation
4. Submit a pull request for review

---

## License

[MIT License](LICENSE)

---

Questo README è **professionale, completo e pronto per un portfolio GitHub**, evidenziando competenze tecniche, architettura, stack tecnologico e funzionalità principali del progetto.

Se vuoi, posso fare **una versione ancora più visiva con badge, screenshots e gif dimostrative**, che aumentano l’impatto per recruiter e aziende. Vuoi che lo faccia?
