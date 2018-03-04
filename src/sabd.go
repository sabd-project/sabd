package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

const databasePath = "/sabdDb/"
const dbName = "iGurbani.sqlite"
const dbURL = "https://www.dropbox.com/sh/qazoft8own8u1na/AADXugSSY0IjqPw2qP1tkL3oa?dl=1" // TODO move to s3

func main() {
	l := log.New(os.Stdout, "", 0)

	dbPath := os.TempDir() + databasePath + dbName
	l.Println("Checking " + dbPath + " to see if it exists or not")
	/*
		does db hosting dir exist?
		does db exist locally?
		download db if not
	*/
	_, err := os.Stat(dbPath)
	if err != nil {

		err2 := os.MkdirAll(os.TempDir()+databasePath, 0777)
		if err2 != nil {
			l.Fatal(err2)
		}

		// now download file
		l.Println("Downloading Gurbani DB from Internet...")
		err3 := downloadFile(dbPath, dbURL)
		if err3 != nil {
			l.Fatal(err3)
		}
	}

	database, _ := sql.Open("sqlite3", dbPath)
	// statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	// statement.Exec()
	// statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
	// statement.Exec("Nic", "Raboy")
	rows, _ := database.Query("select _id,gurmukhi,english_ssk,transliteration from shabad where first_ltr_start like '107,107,107,107,107,107%' order by _id,ang_id,line_id")
	var id int
	var sabd, english, transliteration string
	for rows.Next() {
		rows.Scan(&id, &sabd, &english, &transliteration)
		fmt.Println(strconv.Itoa(id) + ": " + sabd + " " + english + " " + transliteration)
	}
}

/*
	download db from internet - TODO replace with something more efficient and verbose
	@see https://stackoverflow.com/questions/11692860/how-can-i-efficiently-download-a-large-file-using-go
*/
func downloadFile(filepath string, url string) (err error) {

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
