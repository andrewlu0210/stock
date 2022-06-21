package stock

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

const (
	STOCK_PRICE_COL = "stock_price"
)

var (
	db_host     string
	db_name     string
	db_account  string
	db_password string
	client      *mongo.Client
)

func SetMongo(dbHost, dbName, dbAccount, dbPassword string) {
	db_host = dbHost
	db_name = dbName
	db_account = dbAccount
	db_password = dbPassword
}

// Connect to Mongo DB
func ConnectDb() {
	if client != nil {
		log.Println("Already Connected to MongoDB")
		return
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s",
		db_account,
		db_password,
		db_host,
		db_name)
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	checkError(err)
	defer cancel()
	err = c.Ping(ctx, readpref.Primary())
	checkError(err)
	client = c
}

// Close Mongo Db connection
func DisconnectDb() {
	if client != nil {
		client.Disconnect(context.TODO())
		client = nil
	}
}

// Reset DB
func ResetDb(root, rootPasswd string) {
	DisconnectDb()
	uri := fmt.Sprintf("mongodb://%s:%s@%s/%s",
		root,
		rootPasswd,
		db_host,
		"admin")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	checkError(err)
	defer cancel()
	defer c.Disconnect(ctx)
	db := c.Database(db_name)
	db.Drop(ctx)
	log.Printf("Drop DB[%s] First/n", db_name)
	createIndex(c, STOCK_PRICE_COL, "code", "dateStr")
	deleteAccount(c)
	createAccount(c)
}

func deleteAccount(client *mongo.Client) {
	db := client.Database(db_name)
	err := db.RunCommand(
		context.TODO(),
		bson.D{{Key: "dropUser", Value: db_account}},
	).Err()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("account deleted!")
	}
}

func createAccount(client *mongo.Client) {
	db := client.Database(db_name)
	err := db.RunCommand(
		context.TODO(),
		bson.D{
			{Key: "createUser", Value: db_account},
			{Key: "pwd", Value: db_password},
			{
				Key: "roles",
				Value: bson.A{
					bson.D{
						{Key: "role", Value: "dbOwner"},
						{Key: "db", Value: db_name},
					},
				},
			},
		},
	).Err()
	if err != nil {
		//log.Fatal(err)
		log.Println(err)
	} else {
		log.Println("account created!")
	}
}

func createIndex(client *mongo.Client, colName, key1, key2 string) {
	col := client.Database(db_name).Collection(colName)
	opts := options.CreateIndexes().SetMaxTime(10 * time.Second)
	newTrue := true
	uniqueOpts := &options.IndexOptions{Unique: &newTrue}
	indexNames, err := col.Indexes().CreateMany(
		context.TODO(),
		[]mongo.IndexModel{
			{Keys: bsonx.Doc{{Key: key1, Value: bsonx.Int32(1)}}},
			{Keys: bsonx.Doc{{Key: key2, Value: bsonx.Int32(1)}}},
			{
				Keys: bsonx.Doc{
					{Key: key1, Value: bsonx.Int32(1)},
					{Key: key2, Value: bsonx.Int32(1)},
				},
				Options: uniqueOpts,
			},
		},
		opts,
	)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Index", indexNames, "Successfully createed at collection", colName)
	}
}
