package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func main() {
	// Get the PostgreSQL database connection parameters from environment variables
	user := os.Getenv("username")
	password := os.Getenv("password")
	dbname := os.Getenv("dbname")
	host := os.Getenv("instancename")
	port := os.Getenv("port")

	if port == "" {
		port = "5432"
	}

	// Log env variables of rtrubleshooting purpose
	fmt.Println(os.Environ())

	// Construct the PostgreSQL database connection string
	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	scraper := twitterscraper.New()
	//scraper = scraper.SetSearchMode(twitterscraper.SearchLatest)
	//scraper = scraper.WithDelay(5)

	for {
		for tweet := range scraper.SearchTweets(context.Background(),
			"#tanzu OR #vmware OR #tanzuvanguard -filter:retweets", 500) {
			if tweet.Error != nil {
				panic(tweet.Error)
			}

			// Join the hashtag slice into a single comma-separated string
			hashtags := strings.Join(tweet.Hashtags, ",")

			// Insert the data into the database
			_, err = db.Exec("INSERT INTO tweets (id, time, username, text, hashtags) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING", tweet.ID, tweet.TimeParsed, tweet.Username, tweet.Text, hashtags)
			if err != nil {
				log.Fatalf("Error inserting tweet into the database: %v", err)

			}
			fmt.Println(tweet.Text)

		}
		time.Sleep(60 * time.Minute)
	}
}
