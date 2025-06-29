package main

// Creating this all using struct and slices only not database
import (
	"encoding/json" //Encode data into json when sending to postman
	"fmt"
	"log"
	"math/rand" //For id generator for movies
	"net/http"  //For allowing to create server
	"strconv"   //For converting string to int

	"github.com/gorilla/mux" //For creating router
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}
type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// createing slice of movie
var movies []Movie

// Get ALL MOVIES
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// Delete a movie
// This function will delete a movie from the slice
// It will take the id from the request and find the movie in the slice
// If found, it will remove the movie from the slice and return the updated slice
// It uses the mux package to get the URL parameters and sets the content type to application/json
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
		}
	}
	json.NewEncoder(w).Encode(movies)
}

// Get A Movie
// This function will return a single movie based on the id passed in the URL
// It will search through the movies slice and return the movie if found
// If not found, it will return an empty response
// It uses the mux package to get the URL parameters
// It sets the content type to application/json before sending the response
func getmovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// Create a Movie
// This function will create a new movie based on the data passed in the request body
// It will generate a random ID for the movie and append it to the movies slice
// It uses the json package to decode the request body into a Movie struct
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

// Update a movie
// This function will update an existing movie based on the ID passed in the URL
// It will search through the movies slice, find the movie with the matching ID, and update it with the new data from the request body
// It will remove the old movie from the slice and append the updated movie
// It uses the mux package to get the URL parameters and the json package to decode the request body
func updateMovie(w http.ResponseWriter, r *http.Request) {
	//set json content type
	w.Header().Set("Content-Type", "application/json")
	//params
	params := mux.Vars(r)

	for index, item := range movies {
		if item.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func main() {
	r := mux.NewRouter()

	//adding some movies prior for testing and observing
	movies = append(movies, Movie{ID: "1", Isbn: "11111", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "22222", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

	r.HandleFunc("/movies", getMovies).Methods("GET")           //Get ALL movies func
	r.HandleFunc("/movies/{id}", getmovie).Methods("GET")       //Get A Movie func
	r.HandleFunc("/movies", createMovie).Methods("POST")        //Create a Movie func
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")    //Update a movie func
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE") //Delete a movie func

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
