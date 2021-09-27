package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const DBName string = "jokes-api"
const JokesCollectionName string = "jokes"

var JokeNotFound = errors.New("joke not found")

type Database struct {
	Client          *mongo.Client
	JokesCollection *mongo.Collection
}

func NewDatabase() (*Database, error) {
	var db Database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db.Client = client
	collection := client.Database(DBName).Collection(JokesCollectionName)
	db.JokesCollection = collection
	return &db, err
}

func (d *Database) GetJokes(ctx context.Context) ([]models.Joke, error) {
	cur, err := d.JokesCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	result := []models.Joke{}
	if err := cur.All(ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (d *Database) AddJoke(ctx context.Context, title, body string) (models.Joke, error) {
	id := primitive.NewObjectID().Hex()
	joke := models.NewJoke(id, title, body, 0)

	_, err := d.JokesCollection.InsertOne(ctx, joke)
	return joke, err
}

func (d *Database) GetJokeByText(ctx context.Context, text string) ([]models.Joke, error) {
	filter := bson.M{"$or": []interface{}{bson.D{{Key: "body", Value: primitive.Regex{Pattern: text, Options: ""}}}, bson.D{{Key: "title", Value: primitive.Regex{Pattern: text, Options: ""}}}}}
	cur, err := d.JokesCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	result := []models.Joke{}
	if err := cur.All(ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (d *Database) GetJokeByID(ctx context.Context, id string) (models.Joke, error) {
	filter := bson.D{{Key: "id", Value: id}}
	var joke models.Joke
	err := d.JokesCollection.FindOne(ctx, filter).Decode(&joke)
	if err == mongo.ErrNoDocuments {
		return joke, JokeNotFound
	} else if err != nil {
		log.Fatal(err)
	}
	return joke, nil
}

func (d *Database) GetRandomJokes(ctx context.Context) ([]models.Joke, error) {
	pipeline := []bson.M{{"$sample": bson.D{{Key: "size", Value: 300}}}}
	cur, err := d.JokesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	result := []models.Joke{}
	if err := cur.All(ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}

func (d *Database) GetFunniestJokes(ctx context.Context) ([]models.Joke, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "score", Value: -1}})

	cur, err := d.JokesCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)
	result := []models.Joke{}
	if err := cur.All(ctx, &result); err != nil {
		return result, err
	}
	return result, nil
}
