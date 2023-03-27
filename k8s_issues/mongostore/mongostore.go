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
		_, err := collection.InsertOne(context.Background(), issue)
		if err != nil {
			return err
		}
	}

	return nil
}
