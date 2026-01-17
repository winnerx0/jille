# Jille

A modern, minimal **Voting System** built with **Go** and the **Fiber** web framework.

Jille follows **Clean Architecture** principles to ensure maintainability, testability, and separation of concerns.

---

## Features

- **User Authentication**: Secure JWT-based authentication (Access & Refresh tokens).
- **Poll Management**: Create, view, and manage polls and their options.
- **Voting System**: Secure and reliable voting mechanism.
- **Clean Architecture**: Domain-driven design with Hexagonal layers.
- **Auto-generated Documentation**: Swagger integration for API exploration.
- **Data Persistence**: Robust PostgreSQL integration with GORM.

---

## 🛠️ Tech Stack

- **Language**: [Go](https://go.dev/) (1.25+)
- **Web Framework**: [Fiber v3](https://docs.gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Security**: [JWT](https://github.com/golang-jwt/jwt) & [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- **Validation**: [Go Playground Validator](https://github.com/go-playground/validator)

---

## 📁 Project Structure

```text

├── api/                
│   └── middleware    
├── app/                  
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

### 2. Environment Setup

Copy the example environment file and fill in your credentials:

```.env
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
### 3. Running the Application

You can use the provided `Makefile`:

```bash
# To run the application in development mode
make run

# To build the binary
make build
```

The server will start at `http://localhost:9000`.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
