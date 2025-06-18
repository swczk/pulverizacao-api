package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect(uri string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Configurações específicas para MongoDB Atlas
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(10).
		SetMinPoolSize(1).
		SetMaxConnIdleTime(30 * time.Second).
		SetServerSelectionTimeout(30 * time.Second)

	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, err
	}

	// Verificar conexão
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
