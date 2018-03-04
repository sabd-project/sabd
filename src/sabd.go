package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	//_ "github.com/mattn/go-sqlite3"
)

const databasePath = "/sabdDb/"
const dbName = "iGurbani.sqlite"
const dbURL = "https://www.dropbox.com/sh/qazoft8own8u1na/AADXugSSY0IjqPw2qP1tkL3oa?dl=1" // TODO move to s3

func main() {
	l := log.New(os.Stdout, "", 0)

	dbPath := os.TempDir() + databasePath + dbName

	// does db hosting dir exist?
	// does db exist locally?
	// download db if not
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err2 := os.MkdirAll(os.TempDir()+databasePath, 0644)
		if err != nil {
			l.Fatal(err2)
		}

		// now download file
		err3 := downloadFile(dbPath, dbURL)
		if err3 != nil {
			l.Fatal(err3)
		}
	}

	// database, _ := sql.Open("sqlite3", "./nraboy.db")
	// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	// statement.Exec()
	// statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	// statement.Exec("Nic", "Raboy")
	// rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
	// var id int
	// var firstname string
	// var lastname string
	// for rows.Next() {
	// 	rows.Scan(&id, &firstname, &lastname)
	// 	fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
	// }
}

// https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
func downloadFile(filepath string, url string) (err error) {
	fmt.Println("what")
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data into RAM - TODO stream straight to disk
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
