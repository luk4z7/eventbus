package adapters

import (
	"context"
	"eventhandler/entities"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type (
	TransactionRepository struct {
		client *mongo.Client
	}

	TransactionModel struct {
		ID            string    `bson:"_id,omitempty"`
		TransactionID string    `bson:"transaction_id"`
		Origin        string    `bson:"origin"`
		UserID        string    `bson:"user_id"`
		Amount        float64   `bson:"amount"`
		Kind          string    `bson:"kind"`
		CreatedAt     time.Time `bson:"created_at"`
	}
)

func NewTransactionRepository() *TransactionRepository {
	client, err := NewMongoConnection()
	if err != nil {
		panic(err)
	}

	return &TransactionRepository{
		client: client,
	}
}

func (t *TransactionRepository) Save(entry entities.Transaction) error {
	coll := t.client.Database("admin").Collection("transactions")

	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Starts a session on the client
	session, err := t.client.StartSession()
	if err != nil {
		return err
	}

	// Defers ending the session after the transaction is committed or ended
	defer session.EndSession(context.TODO())

	model := TransactionModel{
		TransactionID: entry.ID,
		Origin:        entry.Origin,
		UserID:        entry.UserID,
		Amount:        entry.Amount,
		Kind:          entry.Kind,
		CreatedAt:     time.Now(),
	}

	// Inserts multiple documents into a collection within a transaction,
	// then commits or ends the transaction
	_, err = session.WithTransaction(context.TODO(), func(ctx mongo.SessionContext) (interface{}, error) {
		res, err := coll.InsertOne(context.TODO(), model)
		if err != nil {
			log.Fatal(err)
		}

		return res, err
	}, txnOptions)

	return err
}

func NewMongoConnection() (*mongo.Client, error) {
	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s`,
		os.Getenv("MONGO_USERNAME"),
		os.Getenv("MONGO_PASSWORD"),
		os.Getenv("MONGO_HOST"),
		os.Getenv("MONGO_DATABASE"),
	)

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
