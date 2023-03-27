package main

import (
	"fmt"
	"log"
    "os"
    "context"

	"k8s_issues/httpclient"
	"k8s_issues/mongostore"
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

	fmt.Println("Issues successfully inserted into MongoDB.")
}


