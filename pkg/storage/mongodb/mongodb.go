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

// DBName used to set database name.
const DBName string = "jokes-api"

// JokesCollectionName used to set collection name.
const JokesCollectionName string = "jokes"

// ErrJokeNotFound describes the error when the joke is not found.
var ErrJokeNotFound = errors.New("joke not found")

// Database struct.
type Database struct {
	Client          *mongo.Client
	JokesCollection *mongo.Collection
	DataSize        int
}

// NewDatabase creating a new Database object.
func NewDatabase() (*Database, error) {
	var db Database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	db.Client = client
	collection := client.Database(DBName).Collection(JokesCollectionName)
	db.JokesCollection = collection
	dataSize, err := db.JokesCollection.CountDocuments(ctx, bson.D{})
	db.DataSize = int(dataSize)
	return &db, err
}

// GetJokes method returns all jokes.
func (d *Database) GetJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error) {
	cur, err := d.JokesCollection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(ctx)

	result := []models.Joke{}

	if err := cur.All(ctx, &result); err != nil {
		return result, 0, err
	}
	if skip > d.DataSize {
		return []models.Joke{}, 0, nil
	}
	if seed > d.DataSize {
		return result[skip:d.DataSize], d.DataSize, nil
	}
	return result[skip : skip+seed], d.DataSize, nil
}

// AddJoke method creating new joke.
func (d *Database) AddJoke(ctx context.Context, title, body string) (models.Joke, error) {
	id := primitive.NewObjectID().Hex()
	joke := models.NewJoke(id, title, body, 0)

	_, err := d.JokesCollection.InsertOne(ctx, joke)
	d.DataSize++

	return joke, err
}

// GetJokesByText returns jokes which contain the desired text.
func (d *Database) GetJokesByText(ctx context.Context, skip, seed int, text string) ([]models.Joke, int, error) {
	filter := bson.M{"$or": []interface{}{
		bson.D{{Key: "body", Value: primitive.Regex{Pattern: text, Options: "i"}}},
		bson.D{{Key: "title", Value: primitive.Regex{Pattern: text, Options: "i"}}},
	}}

	result := []models.Joke{}

	cur, err := d.JokesCollection.Find(ctx, filter)
	if err != nil {
		return result, 0, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &result); err != nil {
		return result, 0, err
	}
	if skip > len(result) {
		return []models.Joke{}, 0, nil
	}
	if seed > len(result) {
		return result[skip:], len(result), nil
	}
	return result[skip : skip+seed], len(result), nil
}

// GetJokeByID returns joke that has the same id.
func (d *Database) GetJokeByID(ctx context.Context, id string) (models.Joke, error) {
	filter := bson.D{{Key: "id", Value: id}}

	var joke models.Joke

	err := d.JokesCollection.FindOne(ctx, filter).Decode(&joke)
	if err == mongo.ErrNoDocuments {
		return joke, ErrJokeNotFound
	} else if err != nil {
		return joke, err
	}
	return joke, nil
}

// GetRandomJokes returns random jokes.
func (d *Database) GetRandomJokes(ctx context.Context, seed int) ([]models.Joke, int, error) {
	pipeline := []bson.M{{"$sample": bson.D{{Key: "size", Value: seed}}}}

	result := []models.Joke{}

	cur, err := d.JokesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return result, 0, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &result); err != nil {
		return result, 0, err
	}
	if seed > d.DataSize {
		return result[:d.DataSize], d.DataSize, nil
	}
	return result, d.DataSize, nil
}

// GetFunniestJokes returns jokes, sorted by score.
func (d *Database) GetFunniestJokes(ctx context.Context, skip, seed int) ([]models.Joke, int, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"score": -1})

	result := []models.Joke{}

	cur, err := d.JokesCollection.Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return result, 0, err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &result); err != nil {
		return result, 0, err
	}
	if skip > d.DataSize {
		return []models.Joke{}, 0, nil
	}
	if seed > d.DataSize {
		return result[skip:d.DataSize], d.DataSize, nil
	}
	return result[skip:seed], d.DataSize, nil
}
