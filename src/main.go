package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	e.Static("/", "public")
	// Route for song search
	e.GET("/search", searchSongs)

	// Route for streaming songs
	e.GET("/stream/:songID", streamSong)

	e.POST("/user/:userID/favorites", saveFavoriteSong)

	e.POST("/upload", uploadSong)

	e.Logger.Fatal(e.Start(":1323"))
}
func uploadSong(c echo.Context) error {
	/*
		uploadSong handles the file upload functionality in a web application.
		It receives a file from the request, opens it, creates a destination file,
		and then copies the contents of the source file to the destination file.
		Finally, it returns a JSON response indicating the success or failure of the file upload.

		Example Usage:
			// Upload a file using the /upload endpoint
			curl -X POST -F "file=@/path/to/file" http://localhost:1323/upload

		Inputs:
			- c (echo.Context): The context object representing the HTTP request and response.

		Outputs:
			- If the file upload is successful, it returns a JSON response with a status code of 200 (OK)
			  and the message "File uploaded successfully".
			- If there is an error during the file upload process, it returns a JSON response with a status code
			  of 400 (Bad Request) or 500 (Internal Server Error) depending on the type of error, along with the corresponding error message.
	*/
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()

	dst, err := os.Create(filepath.Join("uploads", file.Filename))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "File uploaded successfully")
}
func saveFavoriteSong(c echo.Context) error {
	userID := c.Param("userID")     // Get the user ID from the request
	songID := c.FormValue("songID") // Get the song ID from the form data

	// Add the song ID to the user's list of favorite songs
	err := updateUserFavorites(userID, songID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, "Song saved to favorites")
}

// Handler function for the streaming endpoint
func streamSong(c echo.Context) error {
	songID := c.Param("songID") // Get the song ID from the request

	// Retrieve the song file based on the song ID
	// This is a placeholder for your file retrieval logic
	filePath := getSongFilePath(songID)
	if filePath == "" {
		return c.String(http.StatusNotFound, "Song not found")
	}

	// Open the song file
	file, err := os.Open(filePath)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer file.Close()

	// Stream the song file to the user
	c.Response().Header().Set(echo.HeaderContentType, "audio/mpeg")
	c.Response().WriteHeader(http.StatusOK)
	io.Copy(c.Response().Writer, file)

	return nil
}

// getSongFilePath retrieves the file path of the song based on the song ID
func getSongFilePath(songID string) string {
	// Define the directory where your song files are stored
	songsDirectory := "path/to/your/songs/directory"

	// Construct the file name using the song ID
	// This example assumes that the song files are named with their ID and have a .mp3 extension
	fileName := songID + ".mp3"

	// Construct the full file path
	fullFilePath := filepath.Join(songsDirectory, fileName)

	// Check if the file exists
	if _, err := os.Stat(fullFilePath); os.IsNotExist(err) {
		// The file does not exist
		return ""
	}

	// Return the full file path
	return fullFilePath
}

// Handler function for the search endpoint
// searchSongs handles the search endpoint of a web application.
// It retrieves the search query from the request, calls a function called 'searchDatabaseForSongs' to search for songs in a database based on the query, and returns the search results as a JSON response.
//
// Parameters:
// - c (echo.Context): The context object representing the HTTP request and response.
//
// Returns:
// - error: If an error occurs during the search, it returns a JSON response with a status code of 500 (Internal Server Error) and the error message.
// - nil: If the search is successful, it returns a JSON response with a status code of 200 (OK) and the search results.
func searchSongs(c echo.Context) error {
	query := c.QueryParam("query") // Get the search query from the request
	songs, err := searchDatabaseForSongs(query)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, songs)
}
