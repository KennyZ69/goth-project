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
	// ! Just commenting out the pwd validation for easier testing.
	// if valid := PasswordValidator(string(usr.Password)); !valid {
	// 	return fmt.Errorf("invalid password")
	// }
	// var id uint
	// err := DB.QueryRow("SELECT COALESCE(MAX(user_id), 0) + 1 FROM users").Scan(&id)
	// if err != nil {
	// 	return err
	// }
	// usr.Id = id

	_, err := DB.Exec("INSERT INTO users (user_id, username, email, password_hash) VALUES ($1, $2, $3, $4)", usr.Id, usr.Username, usr.Email, usr.Password)
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

func UsrId() (uint, error) {
	var id uint
	err := DB.QueryRow("SELECT COALESCE(MAX(user_id), 0) + 1 FROM users").Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
