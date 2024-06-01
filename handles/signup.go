package handles

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
		return nil
	} else if r.Method == "GET" {
		HandleComponents(w, r)
		return nil
	}
	return fmt.Errorf("unsupported method")
}
