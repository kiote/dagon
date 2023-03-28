package main

import (
	"fmt"
	"log"
    "os"
    "context"

	"k8s_issues/httpclient"
	"k8s_issues/mongostore"
	"k8s_issues/telegram"
)

func main() {
	owner := "kubernetes"
	repo := "kubernetes"
	tag := "good first issue"

	issues, err := httpclient.FetchIssues(owner, repo, tag)
	if err != nil {
		log.Fatalf("Error fetching issues: %v", err)
	}

    mongoURI := os.Getenv("MONGO_URL")
	if mongoURI == "" {
		log.Fatalf("MONGO_URL environment variable not set")
	}
    mongoClient, err := mongostore.Connect(mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	dbName := "github_issues"
	collectionName := "issues"

	err = mongostore.InsertIssues(mongoClient, dbName, collectionName, issues)
	if err != nil {
		log.Fatalf("Error inserting issues: %v", err)
	}

	count, err2 := mongostore.CountIssues(mongoClient, dbName, collectionName)
	if err2 != nil {
		log.Fatalf("Error counting issues: %v", err2)
	}
	fmt.Printf("Total number of issues: %d\n", count)

	// if total number of issues > 21, send a message to telegram
	// since that means we have a new issue
	if count > 21 {
		// send message to telegram
		botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
		if botToken == "" {
			log.Fatalf("TELEGRAM_BOT_TOKEN environment variable not set")
		}
		chatID := os.Getenv("TELEGRAM_CHAT_ID")
		if chatID == "" {
			log.Fatalf("TELEGRAM_CHAT_ID environment variable not set")
		}
		err3 := telegram.SendMessage(botToken, chatID, "There is a new issue!")
		if err3 != nil {
			log.Fatalf("Error sending message: %v", err)
		} else {
			log.Println("Message sent successfully!")
		}
	}
}


