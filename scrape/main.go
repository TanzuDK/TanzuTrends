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
	//Define port Staticly
	port := "5432"

	// Get the PostgreSQL database connection parameters from files
	content, err := ioutil.ReadFile("/bindings/tanzutrends-db/username")
	user := string(content)
	content, err = ioutil.ReadFile("/bindings/tanzutrends-db/password")
	password := string(content)
	content, err = ioutil.ReadFile("/bindings/tanzutrends-db/instancename")
	host := string(content)
	content, err = ioutil.ReadFile("/bindings/tanzutrends-db/dbname")
	dbname := string(content)

	// Get the PostgreSQL database connection parameters from ENV Variables
	if os.Getenv("POSTGRES_USER") != "" {
		user = os.Getenv("POSTGRES_USER")
	}
	if os.Getenv("POSTGRES_PASSWORD") != "" {
		password = os.Getenv("POSTGRES_PASSWORD")
	}
	if os.Getenv("POSTGRES_DB") != "" {
		dbname = os.Getenv("POSTGRES_DB")
	}
	if os.Getenv("POSTGRES_HOST") != "" {
		host = os.Getenv("POSTGRES_HOST")
	}
	if os.Getenv("POSTGRES_PORT") != "" {
		port = os.Getenv("POSTGRES_PORT")
	}

	fmt.Println("Username : " + user)
	fmt.Println("Password : " + password)
	fmt.Println("Host : " + host)
	fmt.Println("DB Name : " + dbname)

	// Construct the PostgreSQL database connection string
	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbname + "?sslmode=disable"

	// Connect to the PostgreSQL database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalln("Error connecting to the database: %v", err)
	}
	defer db.Close()

	// Create Database
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS tweets (id BIGINT PRIMARY KEY,	time TEXT, username TEXT, text TEXT, hashtags TEXT)")
	if err != nil {
		log.Println("Database Already exist")
	}

	scraper := twitterscraper.New()
	//scraper = scraper.SetSearchMode(twitterscraper.SearchLatest)
	//scraper = scraper.WithDelay(5)
	scraper.Login("tanzutrend1900", "7@3aQZiV8Sc#EX")

	for {

		for tweet := range scraper.SearchTweets(context.Background(),
			"#tanzu OR #vmware OR #tanzuvanguard OR #tmc OR #tap OR #tkg -filter:retweets", 500) {
			if tweet.Error != nil {
				panic(tweet.Error)
			}

			// Join the hashtag slice into a single comma-separated string
			hashtags := strings.Join(tweet.Hashtags, ",")
			log.Println("Hashtag variable :" + hashtags)

			// Insert the data into the database
			_, err = db.Exec("INSERT INTO tweets (id, time, username, text, hashtags) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING", tweet.ID, tweet.TimeParsed, tweet.Username, tweet.Text, hashtags)
			log.Println("ID 		: " + tweet.ID)
			log.Println("Username 	: " + tweet.Username)
			log.Println("Text 		: " + tweet.Text)
			log.Println("hashtags 	: " + hashtags)

			if err != nil {
				log.Fatalln("!!!")
				log.Fatalln("!!! Error inserting tweet into the database: %v", err)
				log.Fatalln("!!!")

			}

		}
		time.Sleep(60 * time.Minute)
	}
}
