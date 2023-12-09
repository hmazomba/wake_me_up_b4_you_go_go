package main

// Song represents the structure of our resource
type Song struct {
	ID     string `bson:"_id,omitempty" json:"id"`
	Title  string `bson:"title" json:"title"`
	Artist string `bson:"artist" json:"artist"`
	// Add other metadata fields as needed
}

type User struct {
	ID            string   `bson:"_id,omitempty" json:"id"`
	Username      string   `bson:"username" json:"username"`
	FavoriteSongs []string `bson:"favorite_songs" json:"favorite_songs"`
}
