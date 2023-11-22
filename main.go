package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json: "id"`
	Isbn     string    `json: "isbn"`
	Title    string    `json: "title`
	Director *Director `json: "director"`
}
type Director struct {
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
}

// movie slice of type Movie
var movies []Movie

// get all movies
func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

// delete function
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	for index, item := range movies {
		if item.ID == params["id"] {
			// here we get the movie with id at index and using the append func
			// we remove the movie and other movies would shifted into the place the previous movie held
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(movies)
}

// get a movie
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}

// create movie
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)    // we send moovie object from frontend and want to decode it
	movie.ID = strconv.Itoa(rand.Intn(100000000)) // give the ID a random number
	movies = append(movies, movie)
	// return movie
	json.NewEncoder(w).Encode(movie)
}

// update movie
func updateMovie(w http.ResponseWriter, r *http.Request) {
	// set json content type
	w.Header().Set("Content-type", "application/json")
	// get params
	params := mux.Vars(r)
	// loop over movie
	for index, item := range movies {
		if item.ID == params["id"] {
			// delete movie with the id sent
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			// add movie that was sent in the body from postman or frontend
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			movies = append(movies, movie)
			// return movie
			json.NewEncoder(w).Encode(movie)
		}
	}

	// add movie that was sent in the body from postman or frontend
}

// main function
func main() {
	r := mux.NewRouter()

	movies = append(movies, Movie{
		ID: "1", Isbn: "423617", Title: "Movie one", Director: &Director{
			FirstName: "John", LastName: "James",
		},
	})
	movies = append(movies, Movie{
		ID: "2", Isbn: "423618", Title: "Movie two", Director: &Director{
			FirstName: "Jane", LastName: "Won",
		},
	})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
