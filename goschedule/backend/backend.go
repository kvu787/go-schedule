package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kvu787/goschedule"
)

// Scrape will begin a full time schedule scrape and store results in a database.
// Parameter url must be a the time schedule page listing departments and colleges.
func Scrape(url string, db *sql.DB) {
	if err := db.Ping(); err != nil {
		panic("Bad db connection")
	}
	body, err := get(url)
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch time schedule root at %s: %s", url, err))
	}
	colleges, err := goschedule.ExtractColleges(body)
	if err != nil {
		log.Println(err)
	}
	uniqueDepts := make(map[string]int)
	for _, college := range colleges {
		if err := goschedule.Insert(db, college); err != nil {
			log.Println(err)
		}
		depts, err := goschedule.ExtractDepts(body[college.Start:college.End], college.Abbreviation, url, &uniqueDepts)
		if err != nil {
			log.Println(err)
		}
		for _, dept := range depts {
			if err := goschedule.Insert(db, dept); err != nil {
				log.Println(err)
			}
		}
	}
}

// parseConfig reads a JSON format byte slice into a map.
func parseConfig(config []byte) (result map[string]interface{}) {
	json.Unmarshal(config, &result)
	return
}

// get requests a url with the given client and returns the bytes
// of the response body if successful.
// A response with a non-2XX/3XX status code is considered an error.
func get(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Fetch.Get error while using client.Get: %#v\nLink: %v", err, url)
	}
	defer resp.Body.Close()
	if resp.StatusCode > 399 || resp.StatusCode < 200 {
		return "", fmt.Errorf("Fetch.Get returned with a non-2XX/3XX status code: %d\nLink: %v", resp.StatusCode, url)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Fetch.Get error in reading response body: %v", err)
	}
	return string(body), nil
}
