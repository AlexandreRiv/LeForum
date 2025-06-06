package storage

import (
    "database/sql"
    "fmt"
    _ "github.com/go-sql-driver/mysql"
    "os"
)

var DB *sql.DB

func InitDB() error {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWD_USER"),
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_NAME"),
    )

    var err error
    DB, err = sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    return DB.Ping()
}

func SaveUserIfNotExists(email, username string) error {
    var exists bool
    err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE mail=?)", email).Scan(&exists)
    if err != nil {
        return err
    }

    if !exists {
        _, err := DB.Exec("INSERT INTO users (username, mail, password) VALUES (?, ?, '')", username, email)
        if err != nil {
            return err
        }
    }
    return nil
}

func CreateUser(username, email, password string) error {
    _, err := DB.Exec(
        "INSERT INTO users (username, mail, password, darkmode) VALUES (?, ?, ?, ?)",
        username, email, password, 0,
    )
    return err
}

func GetUserByEmail(email string) (*User, error) {
    user := &User{}
    err := DB.QueryRow(
        "SELECT id, username, mail, password FROM users WHERE mail = ?",
        email,
    ).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
    return user, err
}

type User struct {
    ID       int
    Username string
    Email    string
    Password string
}
