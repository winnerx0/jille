# Jille

A modern, minimal **Voting System** built with **Go** using **Fiber** and **Tanstack Start**.

Jille follows **Clean Architecture** principles to ensure maintainability, testability, and separation of concerns.

---

## Features

- **User Authentication**: Secure JWT-based authentication.
- **Poll Management**: Create, view, and manage polls and their options.
- **Voting System**: Secure and reliable voting mechanism.
- **Clean Architecture**: Domain-driven design with Hexagonal layers.
- **Data Persistence**: Robust PostgreSQL integration with GORM.

---

## 🛠️ Tech Stack

- **Language**: [Go](https://go.dev/) (1.25+)
- **Backend Framework**: [Fiber v3](https://docs.gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Security**: [JWT](https://github.com/golang-jwt/jwt) & [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Validation**: [Go Playground Validator](https://github.com/go-playground/validator)
- **Frontend Framwork**: [Tanstack Start](https://tanstack.com/start/latest/docs/)

---

## 📁 Project Structure

```text

├── api/                
│   └── middleware    
├── app/                  
├── client/                          
├── cmd/
│   └── api/               
├── config/  
├── infra/                 
│   ├── database            
│   └── persistence       
├── internal/             
│   ├── domain/            
│   ├── application/     
│   ├── delivery/      
│   ├── common/          
│   └── utils/            
├── Makefile              
├── go.mod              
└── .env.example          
```

---

## ⚙️ Getting Started

### 1. Prerequisites

- **Go 1.25** or higher.
- **PostgreSQL** instance.
- **Bun 1.2** runtime or higher

### 2. Environment Setup

Copy the example environment file and fill in your credentials:

```.env
# .env.example
PORT=
JWT_ACCESS_TOKEN_SECRET=
JWT_REFRESH_TOKEN_SECRET=
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=
DB_TIMEZONE=

```

```.env
# client/.env.local
BACKEND_URL=http://localhost:9000

```
### 3. Running the Application

You can use the provided `Makefile`:

```bash
# To run the application in development mode
make run

# To build the binary
make build

# To run the frontend
make client
```

The server will start at `http://localhost:9000`.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
