package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection
var client *mongo.Client
var err error

// Set up MongoDB connection
func init() {
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("musicdb").Collection("songs")
}

// searchDatabaseForSongs searches the MongoDB database for songs matching the query
func searchDatabaseForSongs(query string) ([]Song, error) {
	// searchDatabaseForSongs searches a MongoDB database for songs matching a given query.
	// It uses the `$or` operator to search for songs with a matching title or artist.
	// The function returns a slice of `Song` objects that match the query.
	//
	// Inputs:
	//   - query (string): The search query to match against song titles or artists.
	//
	// Outputs:
	//   - []Song: A slice of `Song` objects that match the search query.
	//   - error: An error object if any error occurred during the search process.

	var songs []Song
	cursor, err := collection.Find(context.TODO(), bson.M{
		"$or": []bson.M{
			{"title": bson.M{"$regex": query, "$options": "i"}},
			{"artist": bson.M{"$regex": query, "$options": "i"}},
		},
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var song Song
		if err = cursor.Decode(&song); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}

	if err = cursor.Err(); err != nil {
		return nil, err
	}

	return songs, nil
}
func updateUserFavorites(userID, songID string) error {
	// Connect to your NoSQL database
	// Update the user's document to include the new favorite song
	// This will vary depending on which NoSQL database you are using

	// Placeholder for database update logic
	// Replace with actual update operation to your NoSQL database
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"favorite_songs": songID}},
	)
	return err
}

// saveSongMetadataToDatabase saves the metadata of a song to the database
func saveSongMetadataToDatabase(song Song) error {
	_, err := collection.InsertOne(context.TODO(), song)
	return err
}
