package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
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
	Country      string `json:"country"`
	City         string `json:"city"`
}

type RequestStatus string

const (
	StatusPending  RequestStatus = "pending"
	StatusAccepted RequestStatus = "accepted"
	StatusDenied   RequestStatus = "denied"
	StatusReceived RequestStatus = "received"
	StatusSent     RequestStatus = "sent"
)

type Connection_req struct {
	Id          uint          `json:"id"`
	Sender_id   uint          `json:"sender_id"`
	Receiver_id uint          `json:"receiver_id"`
	Status      RequestStatus `json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
}

type Friend struct {
	Id          uint      `json:"id"`
	User_id     uint      `json:"user_id"`
	User_name   string    `json:"user_name"`
	Friend_name string    `json:"friend_name"`
	Friend_id   uint      `json:"friend_id"`
	Created_at  time.Time `json:"created_at"`
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

	_, err = tx.Exec("INSERT INTO user_details (user_id, bio, profile_image, role, country, city) VALUES ($1, $2, $3, $4, $5, $6)", usr.Id, usr.Details.Bio, usr.Details.ProfileImage, usr.Details.Role, usr.Details.Country, usr.Details.City)
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

func UsernameExists(username string) bool {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1", username).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
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

func GetDetailsById(db *sql.DB, id uint) (*UserProfileData, error) {
	var role, bio, profile_image, city, country string
	row := db.QueryRow("SELECT bio, profile_image, role, country, city FROM user_details WHERE user_id = $1", id)
	err := row.Scan(&bio, &profile_image, &role, &country, &city)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		return nil, err
	}
	data := UserProfileData{
		Role:         role,
		Bio:          bio,
		ProfileImage: profile_image,
		Country:      country,
		City:         city,
	}
	return &data, nil
}

func GetConnectionReq(currentUser User, receiver User) (RequestStatus, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM connection_requests WHERE sender_id=$1 AND receiver_id=$2 AND status=$3", currentUser.Id, receiver.Id, StatusPending).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("there was an error getting the count of requests for the button: %s", err)
	}
	if count > 0 {
		return StatusSent, nil
	}
	err = DB.QueryRow("SELECT COUNT(*) FROM connection_requests WHERE sender_id=$1 AND receiver_id=$2 AND status=$3", receiver.Id, currentUser.Id, StatusPending).Scan(&count)
	if count > 0 {
		return StatusReceived, nil
	}
	err = DB.QueryRow("SELECT COUNT(*) FROM connection_requests WHERE sender_id=$1 AND receiver_id=$2 AND status=$3 OR sender_id=$4 AND receiver_id=$5 AND status=$6", currentUser.Id, receiver.Id, StatusAccepted, receiver.Id, currentUser.Id, StatusAccepted).Scan(&count)
	if count > 0 {
		return StatusAccepted, nil
	}
	err = DB.QueryRow("SELECT COUNT(*) FROM connection_requests WHERE sender_id=$1 AND receiver_id=$2 AND status=$3 OR sender_id=$4 AND receiver_id=$5 AND status=$6", currentUser.Id, receiver.Id, StatusDenied, receiver.Id, currentUser.Id, StatusDenied).Scan(&count)
	if count > 0 {
		return StatusDenied, nil
	}
	return "", nil
}

func GetIfHasRequests(user *User) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM connection_requests WHERE receiver_id=$1 AND status=$2", user.Id, StatusPending).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("there was a problem getting the requests of the user: %v", err)
	}
	return count > 0, nil
}

func GetRequestsOnUser(user *User) ([]Connection_req, error) {
	rows, err := DB.Query("SELECT status, sender_id, created_at FROM connection_requests WHERE receiver_id=$1 AND status=$2", user.Id, StatusPending)
	if err != nil {
		return nil, fmt.Errorf("there was an error getting the pending requests on the user: %v", err)
	}
	defer rows.Close()

	var requests []Connection_req
	for rows.Next() {
		var req Connection_req
		if err := rows.Scan(&req.Status, &req.Sender_id, &req.CreatedAt); err != nil {
			return nil, fmt.Errorf("problem scanning the rows on db while getting requests for the user: %v", err)
		}
		requests = append(requests, req)
	}
	return requests, nil
}

func UpdateReqStatus(sender_id, receiver_id uint, newStatus RequestStatus) error {
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("problem beginning transaction: %v", err)
	}
	defer tx.Rollback()

	if newStatus == StatusDenied {
		_, err = tx.Exec("UPDATE connection_requests SET status=$1 expires_at=$2 WHERE sender_id=$3 AND receiver_id=$4", newStatus, time.Now().Add(5*24*time.Hour), sender_id, receiver_id)
		if err != nil {
			return fmt.Errorf("problem updating in the db: %v", err)
		}
	}
	if newStatus == StatusPending {
		_, err = tx.Exec("UPDATE connection_requests SET status=$1 expires_at=$2 WHERE sender_id=$3 AND receiver_id=$4", newStatus, time.Now().Add(30*24*time.Hour), sender_id, receiver_id)
		if err != nil {
			return fmt.Errorf("problem updating in the db: %v", err)
		}
	} else {

		_, err = tx.Exec("UPDATE connection_requests SET status=$1 WHERE sender_id=$2 AND receiver_id=$3", newStatus, sender_id, receiver_id)
		if err != nil {
			return fmt.Errorf("problem updating in the db: %v", err)
		}
	}
	return tx.Commit()
}

func DeleteExpiredRequests() error {
	_, err := DB.Exec(`
		DELETE FROM connection_requests 
		WHERE expires_at < NOW()`)
	if err != nil {
		return fmt.Errorf("problem deleting expired requests: %v", err)
	}
	return nil
}

func GetFriendNameById(id uint) (error, string) {
	query := "SELECT username FROM users WHERE user_id=$1"
	row := DB.QueryRow(query, id)
	var username string
	err := row.Scan(&username)
	if err != nil {
		return err, ""
	}
	return err, username
}

// func GetMessagesByChatId(chatId uint) ([]*ws.Message, error) {
// 	var msgArr []*ws.Message
// 	var msg *ws.Message
// 	var count int
// 	tx, err := DB.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	err = tx.QueryRow("SELECT COUNT(*) FROM messages WHERE chat_id=$1", chatId).Scan(&count)
// 	if err != nil {
// 		return nil, fmt.Errorf("Error counting messages for given chat: %v", err)
// 	}
// 	for _ = range count {
// 		// err = tx.QueryRow("SELECT *")
// 	}
//
// 	return msgArr, tx.Commit()
// }

func SaveMsgsToDb(senderId uint, chatId uint, content string) error {
	tx, err := DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the sender's username from the users table
	var senderName string
	err = tx.QueryRow("SELECT username FROM users WHERE user_id=$1", senderId).Scan(&senderName)
	if err != nil {
		return fmt.Errorf("error getting sender name: %v", err)
	}

	// Insert the message into the messages table
	_, err = tx.Exec(`
        INSERT INTO messages (chat_id, sender_id, sender_name, content, created_at) 
        VALUES ($1, $2, $3, $4, $5)`,
		chatId, senderId, senderName, content, time.Now())
	if err != nil {
		return fmt.Errorf("error inserting message: %v", err)
	}

	// Update the last_message timestamp in the chats table
	_, err = tx.Exec(`
        UPDATE chats 
        SET last_message = $1
        WHERE chat_id = $2`,
		time.Now(), chatId)
	if err != nil {
		return fmt.Errorf("error updating last message: %v", err)
	}

	return tx.Commit()
}
