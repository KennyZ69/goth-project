package handles

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Import the PostgreSQL driver for it to work; ??? idk, it probably wont be used
)

func HandleSignUp(w http.ResponseWriter, r *http.Request) error {
	godotenv.Load()
	pwd := os.Getenv("PSQ_PWD")
	host := os.Getenv("PSQ_HOST")
	port := os.Getenv("PSQ_PORT")
	if r.Method == "POST" {
		db, err := sql.Open("postgres", fmt.Sprintf("user=postgres password=%s host=%s port=%s  dbname=testdb sslmode=disable", pwd, host, port))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// TODO: make validation for email, password, username and then start working on adding to the datatabase
		// newUsr :=
		fmt.Println("opened successfully")
		return nil
	} else if r.Method == "GET" {
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("unsupported method")
}
