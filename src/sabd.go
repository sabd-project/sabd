package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli" // imports as package "cli"
)

const databasePath = "/sabdDb/"
const dbName = "iGurbani.sqlite"
const dbURL = "https://www.dropbox.com/sh/qazoft8own8u1na/AADXugSSY0IjqPw2qP1tkL3oa?dl=1" // TODO move to s3

func main() {
	var searchString, sabdId string;
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "search, s",
			Value:       "",
			Usage:       "First letter search chars",
			Destination: &searchString,
		},
		cli.StringFlag{
			Name:        "display, d",
			Value:       "",
			Usage:       "which Sabd Id to render",
			Destination: &sabdId,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	dbPath := os.TempDir() + databasePath + dbName
	log.Println("Checking " + dbPath + " to see if it exists or not")

	/**
		does db hosting dir exist?
		does db exist locally?
		download db if not
	 */
	fileStats, errStat := os.Stat(dbPath)
	if errStat != nil {

		err2 := os.MkdirAll(os.TempDir()+databasePath, 0777)
		if err2 != nil {
			log.Fatal(err2)
		}

		// now download file
		log.Println("Downloading Gurbani DB from Internet...")
		err3 := downloadFile(dbPath, dbURL)
		if err3 != nil {
			log.Fatal(err3)
		}
	} else {
		log.Printf("db alredy exists at %s %d bytes", dbPath, fileStats.Size)
	}

	/*
		search
	 */
	if searchString != "" {
		log.Printf("Searching for %s", searchString)

		database, _ := sql.Open("sqlite3", dbPath)
		rows, _ := database.Query("select _id,gurmukhi,english_ssk,transliteration from shabad where first_ltr_start like '" + stringToFirstLetterSearch(searchString) + "%' order by _id,ang_id,line_id")

		var id int
		var sabd, english, transliteration string
		for rows.Next() {
			rows.Scan(&id, &sabd, &english, &transliteration)
			log.Printf(strconv.Itoa(id) + ": " + sabd + " " + english + " " + transliteration)
		}
	}

	if sabdId != "" {
		//TODO render gurbani
	}
}

/*
 TODO
*/
func stringToFirstLetterSearch(input string) (output string) {
	//convert strings to rune array then to int arrray
	runes := []int32([]rune(input))

	// https://stackoverflow.com/questions/37532255/one-liner-to-transform-int-into-string/37533144
	return strings.Trim(strings.Join(strings.Split(fmt.Sprint(runes), " "), ","), "[]")
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
