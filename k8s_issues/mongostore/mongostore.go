package mongostore

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"k8s_issues/models"
)

func Connect(uri string) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func InsertIssues(client *mongo.Client, dbName, collectionName string, issues []models.Issue) error {
	collection := client.Database(dbName).Collection(collectionName)

	for _, issue := range issues {
		exists, err := IssueExists(client, dbName, collectionName, issue.URL)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		_, err2 := collection.InsertOne(context.Background(), issue)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func IssueExists(client *mongo.Client, dbName, collectionName, url string) (bool, error) {
	collection := client.Database(dbName).Collection(collectionName)

	filter := map[string]string{"url": url}
	var result models.Issue
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
