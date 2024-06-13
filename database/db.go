package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq" // Import the PostgreSQL driver for it to work; ??? idk, it probably wont be used
)

type User struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

func InitDB() *sql.DB {
	godotenv.Load()
	pwd := os.Getenv("PSQ_PWD")
	host := os.Getenv("PSQ_HOST")
	port := os.Getenv("PSQ_PORT")

	db, err := sql.Open("postgres", fmt.Sprintf("user=postgres password=%s host=%s port=%s  dbname=testdb sslmode=disable", pwd, host, port))
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return db
}

func CreateTables(db *sql.DB) error {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL
		password VARCHAR(100) NOT NULL
    );
    `
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	// Example of creating the users_tokens table
	createTableQuery = `
    CREATE TABLE IF NOT EXISTS user_tokens (
        id SERIAL PRIMARY KEY,
        user_id INTEGER REFERENCES users(id),
        token TEXT NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );
    `
	_, err = db.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

var DB = InitDB()

func HashPwd(pwd string) []byte {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return nil
	}
	return hashedPwd
}

func SaveUser(usr User) error {
	if valid := EmailValidator(usr.Email); !valid {
		return fmt.Errorf("invalid email")
	}
	if valid := UsernameValidator(usr.Username); !valid {
		return fmt.Errorf("invalid username")
	}
	if valid := PasswordValidator(string(usr.Password)); !valid {
		return fmt.Errorf("invalid password")
	}
	var id uint
	err := DB.QueryRow("SELECT COALESCE(MAX(id), 0) + 1 FROM users").Scan(&id)
	if err != nil {
		return err
	}
	usr.Id = id

	_, err = DB.Exec("INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4)", usr.Id, usr.Username, usr.Email, usr.Password)
	// _, err := DB.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", usr.Username, usr.Email, usr.Password)
	if err != nil {
		return err
	}
	return nil
}

func UserExists(usr User) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1 OR email=$2", usr.Username, usr.Email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u *User) ComparePassword(password string) error {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	return err
}
