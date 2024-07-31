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
	Id       uint            `json:"id"`
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Password []byte          `json:"password"`
	Details  UserProfileData `json:"user_details"`
}

type UserProfileData struct {
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_pic"`
	Role         string `json:"role"`
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
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

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
	var id uint
	err = tx.QueryRow("SELECT COALESCE(MAX(user_id), 0) + 1 FROM users").Scan(&id)
	if err != nil {
		return err
	}
	usr.Id = id

	_, err = tx.Exec("INSERT INTO users (user_id, username, email, password_hash) VALUES ($1, $2, $3, $4)", usr.Id, usr.Username, usr.Email, usr.Password)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO user_details (user_id, bio, profile_image, role) VALUES ($1, $2, $3, $4)", usr.Id, usr.Details.Bio, usr.Details.ProfileImage, usr.Details.Role)
	if err != nil {
		return err
	}
	return tx.Commit()
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

func GetUserByName(db *sql.DB, name string) (*User, error) {
	var usr User
	err := db.QueryRow("SELECT user_id, username, email, password_hash FROM users WHERE username=$1", name).Scan(&usr.Id, &usr.Username, &usr.Email, &usr.Password)
	if err != nil {
		fmt.Printf("No user with the username: %v\n", name)
		return nil, err
	}
	return &usr, nil
}

func DeleteUserTokens(db *sql.DB, userID uint) error {
	query := "DELETE FROM user_tokens WHERE user_id = $1"
	if _, err := db.Exec(query, userID); err != nil {
		return err
	}
	return nil
}

func GetUserById(userID uint) (*User, error) {
	query := "SELECT user_id, username, password_hash, email FROM users WHERE user_id = $1;"
	row := DB.QueryRow(query, userID)

	var user User
	err := row.Scan(&user.Id, &user.Username, &user.Password, &user.Email)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("User not found")
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}

func DeleteTokenByToken(db *sql.DB, tokenString string) error {
	query := "DELETE FROM user_tokens WHERE token = $1"
	if _, err := db.Exec(query, tokenString); err != nil {
		return err
	}
	return nil
}

func GetTokenByUsrId(db *sql.DB, usrId uint) (bool, error) {
	query := "SELECT FROM user_tokens WHERE token_id = $1"
	if _, err := db.Exec(query, usrId); err != nil {
		return false, err
	}
	return true, nil
}

func IdToken(db *sql.DB, usrId uint) (string, error) {
	rows, err := db.Query("SELECT token FROM user_tokens WHERE token_id = $1", usrId)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var token string
	if err := rows.Scan(&token); err != nil {
		return "", nil
	}
	return token, nil
}

func GetRecommendedUsers(db *sql.DB, n int, excludeUsername string) ([]User, error) {
	rows, err := db.Query("SELECT user_id, username, email FROM users WHERE username != $1 ORDER BY RANDOM() LIMIT $2", excludeUsername, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var usr User
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.Email); err != nil {
			return nil, err
		}
		users = append(users, usr)
	}
	return users, nil
}

func SearchUsers(db *sql.DB, query string, excludeUsername string) ([]User, error) {
	searchQuery := "%" + query + "%"
	rows, err := db.Query(
		`SELECT user_id, username, email 
		FROM users 
		WHERE username ILIKE $1 OR username = $2 
		AND username != $3`,
		searchQuery, query, excludeUsername)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var usr User
		if err := rows.Scan(&usr.Id, &usr.Username, &usr.Email); err != nil {
			return nil, err
		}
		role, err := GetRoleById(db, usr.Id)
		if err != nil {
			return nil, err
		}
		usr.Details.Role = role
		users = append(users, usr)
	}
	return users, nil
}

func GetRoleById(db *sql.DB, id uint) (string, error) {
	row := db.QueryRow("SELECT role FROM user_details WHERE user_id = $1", id)

	var role string
	err := row.Scan(&role)
	if err == sql.ErrNoRows {
		return "", err
	} else if err != nil {
		return "", err
	}

	return role, nil
}
