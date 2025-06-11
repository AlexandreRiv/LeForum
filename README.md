# LeForum

LeForum is a web forum application built with Go, designed to provide a modern and efficient platform for discussions. This project implements a complete forum system with user authentication, post management, categories, and an interactive UI with both light and dark modes.

[![Go Version](https://img.shields.io/badge/go-1.23+-blue.svg)](https://golang.org)


## 📋 Contents

- [Description](#-description)
- [Features](#-features)
- [Technologies](#-technologies)
- [Project Structure](#-project-structure)
- [Connection](#-connection)
- [Configuration](#-configuration)
- [Database](#-database)
- [Usage](#-usage)


## 🎯 Description

LeForum is a modern web forum application built with Go, offering an efficient platform for online discussions. The project includes a complete forum system with user authentication, post management, categorization, and an interactive user interface with light and dark mode support.

## ✨ Features

### 🔐  Authentication System
- ✅ Traditional sign-up and login (email/password)
- ✅ OAuth integration with GitHub and Google
- ✅ Session management
- ✅ User profile management

### 📝 Content Management
- ✅ Create and view posts
- ✅ Post categorization
- ✅ Like/dislike system
- ✅ Comments and replies
- ✅ Post filtering by category or recency

### 🎨 User Interface
- ✅ Responsive design
- ✅ Light/dark mode toggle with preference persistence
- ✅ User-friendly navigation

## 🛠️ Technologies

### Backend
![go-badge] ![mysql-badge] 
- **Go (Golang)** - Main programming language  
- **HTTP Server** - Go standard library
- **MySQL** - Database
- **Session Management** - Custom session handling
### Frontend
![tailwind-badge] ![javascript-badge] ![html-badge]

- **HTML Templates** - Server-side rendering
- **Tailwind CSS** - CSS framework
- **JavaScript** - Client-side interactivity

### Authentification
![oauth-badge]
- **Custom authentication system**
- **OAuth 2.0** (Google, GitHub)
- **Secure password hashing** with bcrypt

## 📁 project Structure

```
LeForum/
├── cmd/
│   └── server/          # Point d'entrée de l'application
├── internal/
│   ├── api/             # Gestionnaires HTTP et routage
│   │   ├── handlers/    # Gestionnaires de requêtes
│   │   └── middleware/  # Middleware HTTP
│   ├── auth/            # Composants d'authentification
│   │   ├── oauth/       # Fournisseurs OAuth
│   │   └── session/     # Gestion des sessions
│   ├── config/          # Configuration de l'application
│   ├── domain/          # Modèles de domaine
│   ├── service/         # Logique métier
│   └── storage/         # Accès aux données
│       └── repositories/ # Implémentations des repositories
└── web/
    ├── static/          # Ressources statiques (CSS, JS)
    └── templates/       # Templates HTML
```

## ⚙️ Configuration

### Variables d'environnement requises

| Variable | Description | Example              |
|----------|-------------|----------------------|
| `DB_USER` | Database username | `root`               |
| `DB_PASSWD_USER` | Database password | `password123`        |
| `DB_HOST` | 	Database host | `localhost`          |
| `DB_PORT` | Database port | `3306`               |
| `DB_NAME` | Database name | `leforum`            |
| `GITHUB_CLIENT_ID` | GitHub OAuth client ID | `your_github_id`     |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth client secret | `your_github_secret` |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID| `your_google_id`     |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret | `your_google_secret` |

## 🗄️ Database

The application requires a MySQL database with the following main tables:

### Main Tables

- **`users`** - User information
- **`sessions`** - Active user sessions
- **`posts`** - Forum posts
- **`comments`** - Post comments
- **`categories`** - Post categories
- **`affectation`** - Relations between posts and categories
- **`liked_posts`** -  Tracking likes and dislikes

### Basic Schema

```sql
-- Exemple de structure (simplifié)
TABLE USER ( 
id int [primary key, not null, increment],
username varchar(20) [not null],
mail varchar(255) [not null],
password varchar(255) [not null],
newsletters bool [not null, default: false] 
);

CREATE TABLE postsTABLE POST {   
    id int [primary key, not null, increment],
    name varchar(50) [not null] ,
    description text [not null],
    image image  ,
    like int [default: 0],
    dislike int [default: 0],
    id_user int [not null] 
);
```

## 🚀 Connection

To connect to the site, enter the following address in your browser:
```forum.ynov.zeteox.fr```

To connect to the VPS, use the following credentials: \
```ssh -p 50022 guest@server-app.zeteox.fr```\
Use```guest``` as the password.

## 💻 Usage

### Quick Start

1. **Create an account** - Sign up with email/password or via OAuth
2. **Browse categories** -Explore various discussion topics
3. **Create a post** -Share your thoughts and start a conversation
4. **Engage** -Like, dislike, and comment on posts
5. **Customize** -Switch between light and dark modes as you prefer

### OAuth Features

To use OAuth authentication:

1. **GitHub** : Set up an OAuth App in your GitHub settings
2. **Google** : Create a project in Google Cloud Console and enable the OAuth API


## 📝 Contributors

## 📝 Contributors

- **Loic** - [![GitHub][github-badge]](https://github.com/Zeteox/Zeteox)
- **Nathan** - [![GitHub][github-badge]](https://github.com/zoom26042604)
- **Matteo** - [![GitHub][github-badge]](https://github.com/Enoxboo)
- **Alexandre** - [![GitHub][github-badge]](https://github.com/AlexandreRiv)

---

[github-badge]: https://img.shields.io/badge/-GitHub-181717?style=flat&logo=github&logoColor=white
[javascript-badge]: https://img.shields.io/badge/-JavaScript-F7DF1E?style=flat&logo=javascript&logoColor=black
[go-badge]: https://img.shields.io/badge/-Go-00ADD8?style=flat&logo=go&logoColor=white
[tailwind-badge]: https://img.shields.io/badge/-TailwindCSS-06B6D4?style=flat&logo=tailwindcss&logoColor=white
[mysql-badge]: https://img.shields.io/badge/-MySQL-4479A1?style=flat&logo=mysql&logoColor=white
[oauth-badge]: https://img.shields.io/badge/-OAuth-000000?style=flat&logo=oauth&logoColor=white
[html-badge]: https://img.shields.io/badge/-HTML-E34F26?style=flat&logo=html5&logoColor=white
