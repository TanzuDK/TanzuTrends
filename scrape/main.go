package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
	// Get the PostgreSQL database connection parameters from environment variables
	content, err := ioutil.ReadFile("/bindings/tanzutrends-db/username")
	user := string(content)
	content, err = ioutil.ReadFile("/bindings/tanzutrends-db/password")
	password := string(content)
	content, err = ioutil.ReadFile("/bindings/tanzutrends-db/instancename")
	host := string(content)
	content, err = ioutil.ReadFile("/bindings/tanzutrends-db/dbname")
	dbname := string(content)

	fmt.Println("Username : " + user)
	fmt.Println("Password : " + password)
	fmt.Println("Host : " + host)
	fmt.Println("DB Name : " + dbname)
	//user := os.Getenv("username")
	//password := os.Getenv("password")
	//dbname := os.Getenv("dbname")
	//host := os.Getenv("instancename")
	port := os.Getenv("port")

	if port == "" {
		port = "5432"
	}

	// List secret binding
	//entries, err := os.ReadDir("/bindings/tanzutrends-db")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, e := range entries {
	//	fmt.Println(e.Name())
	//}

	// Construct the PostgreSQL database connection string
	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Create Database
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS trends(id serial PRIMARY KEY, create_dtm VARCHAR, order_id VARCHAR, phone VARCHAR, name VARCHAR, address VARCHAR, menu VARCHAR, total_item VARCHAR, pay VARCHAR)")
	if err != nil {
		log.Println("Database Already exist")
	}

	scraper := twitterscraper.New()
	//scraper = scraper.SetSearchMode(twitterscraper.SearchLatest)
	//scraper = scraper.WithDelay(5)

	for {
		for tweet := range scraper.SearchTweets(context.Background(),
			"#tanzu OR #vmware OR #tanzuvanguard OR #tmc -filter:retweets", 500) {
			if tweet.Error != nil {
				panic(tweet.Error)
			}

			// Join the hashtag slice into a single comma-separated string
			hashtags := strings.Join(tweet.Hashtags, ",")
			fmt.Println("Hashtag variable :" + hashtags)

			// Insert the data into the database
			_, err = db.Exec("INSERT INTO "+dbname+" (id, time, username, text, hashtags) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING", tweet.ID, tweet.TimeParsed, tweet.Username, tweet.Text, hashtags)
			fmt.Println("ID 		: " + tweet.ID)
			fmt.Println("Username 	: " + tweet.Username)
			fmt.Println("Text 		: " + tweet.Text)
			fmt.Println("hashtags 	: " + hashtags)

			if err != nil {
				log.Println("!!!")
				log.Println("!!! Error inserting tweet into the database: %v", err)
				log.Println("!!!")

			}
			fmt.Println(tweet.Text)

		}
		time.Sleep(60 * time.Minute)
	}
}
