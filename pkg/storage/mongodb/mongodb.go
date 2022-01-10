package mongodb

import (
	"context"
	"time"

	"github.com/DanilLagunov/jokes-api/pkg/models"
	"github.com/DanilLagunov/jokes-api/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBName used to set database name.
const DBName string = "jokes-api"

// JokesCollectionName used to set collection name.
const JokesCollectionName string = "jokes"

// Database struct.
type Database struct {
	client          *mongo.Client
	jokesCollection *mongo.Collection
}

// NewDatabase creating a new Database object.
func NewDatabase(URI string) (*Database, error) {
	var db Database

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	db.client = client
	collection := client.Database(DBName).Collection(JokesCollectionName)
	db.jokesCollection = collection
	return &db, err
}

// GetJokes method returns a number of jokes given by skip and limit parameters and total amount of jokes.
func (d *Database) GetJokes(ctx context.Context, skip, limit int) ([]models.Joke, int, error) {
	amount, err := d.jokesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return []models.Joke{}, int(amount), err
	}

	options := options.Find()
	options.SetSkip(int64(skip))
	options.SetLimit(int64(limit))

	cur, err := d.jokesCollection.Find(ctx, bson.M{}, options)
	if err != nil {
		return []models.Joke{}, int(amount), err
	}
	defer cur.Close(ctx)

	result := []models.Joke{}

	if err := cur.All(ctx, &result); err != nil {
		return result, int(amount), err
	}

	return result, int(amount), nil
}

// AddJoke method creating new joke.
func (d *Database) AddJoke(ctx context.Context, title, body string, score int) (models.Joke, error) {
	id := primitive.NewObjectID().Hex()
	joke := models.NewJoke(id, title, body, score)

	_, err := d.jokesCollection.InsertOne(ctx, joke)

	return joke, err
}

// GetJokesByText returns a number of jokes, which contain the desired text, given by skip and limit parameters and total amount of found jokes.
func (d *Database) GetJokesByText(ctx context.Context, skip, limit int, text string) ([]models.Joke, int, error) {
	filter := bson.M{"$or": []interface{}{
		bson.M{"body": primitive.Regex{Pattern: text, Options: "i"}},
		bson.M{"title": primitive.Regex{Pattern: text, Options: "i"}},
	}}

	amount, err := d.jokesCollection.CountDocuments(ctx, filter)
	if err != nil {
		return []models.Joke{}, int(amount), err
	}

	options := options.Find()
	options.SetSkip(int64(skip))
	options.SetLimit(int64(limit))

	result := []models.Joke{}

	cur, err := d.jokesCollection.Find(ctx, filter, options)
	if err != nil {
		return result, int(amount), err
	}
	defer cur.Close(ctx)

	if err := cur.All(ctx, &result); err != nil {
		return result, int(amount), err
	}

	return result, int(amount), nil
}

// GetJokeByID returns joke that has the same id.
func (d *Database) GetJokeByID(ctx context.Context, id string) (models.Joke, error) {
	filter := bson.D{{Key: "_id", Value: id}}

	var joke models.Joke

	err := d.jokesCollection.FindOne(ctx, filter).Decode(&joke)
	if err == mongo.ErrNoDocuments {
		return joke, storage.ErrJokeNotFound
	}

	return joke, err
}

// GetRandomJokes returns number of random jokes given by limit parameter and total amount of jokes.
func (d *Database) GetRandomJokes(ctx context.Context, limit int) ([]models.Joke, int, error) {
	amount, err := d.jokesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return []models.Joke{}, int(amount), err
	}

	pipeline := []bson.M{{"$sample": bson.M{"size": limit}}}

	cur, err := d.jokesCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return []models.Joke{}, int(amount), err
	}

	defer cur.Close(ctx)

	result := []models.Joke{}
	if err := cur.All(ctx, &result); err != nil {
		return result, int(amount), err
	}

	return result, int(amount), nil
}

// GetFunniestJokes returns number of sorted jokes given by skip and limit parameters and total amount of jokes.
func (d *Database) GetFunniestJokes(ctx context.Context, skip, limit int) ([]models.Joke, int, error) {
	amount, err := d.jokesCollection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return []models.Joke{}, int(amount), err
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"score": -1})
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	cur, err := d.jokesCollection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, int(amount), err
	}
	defer cur.Close(ctx)

	result := []models.Joke{}

	if err := cur.All(ctx, &result); err != nil {
		return result, int(amount), err
	}

	return result, int(amount), nil
}
