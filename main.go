package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

var Movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(Movies)
	if err != nil {
		log.Fatalf("[ERROR] Can not encode movies. Error: %s\n", err)
	}
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, value := range Movies {
		if value.ID == params["id"] {
			err := json.NewEncoder(w).Encode(value)
			if err != nil {
				log.Fatalf("[ERROR] Can not enconde movie. Error: %s\n", err)
			}
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)

	lastID, _ := strconv.Atoi(Movies[len(Movies)-1].ID)
	movie.ID = strconv.Itoa(lastID + 1)

	Movies = append(Movies, movie)
	err := json.NewEncoder(w).Encode(movie)
	if err != nil {
		log.Fatalf("[ERROR] Can not encode new movie. Error: %s\n", err)
	}
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range Movies {
		if item.ID == params["id"] {
			Movies = append(Movies[:index], Movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.ID = params["id"]
			Movies = append(Movies, movie)
			err := json.NewEncoder(w).Encode(movie)
			if err != nil {
				log.Fatalf("[ERROR] Can not encode updated movie. Error: %s\n", err)
			}
			return
		}
	}
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, value := range Movies {
		if value.ID == params["id"] {
			Movies = append(Movies[:index], Movies[index+1:]...)
			break
		}
	}
	err := json.NewEncoder(w).Encode(Movies)
	if err != nil {
		log.Fatalf("[ERROR] Can not encode remaining movies. Error: %s\n", err)
	}
}

func main() {
	r := mux.NewRouter()

	Movies = append(Movies,
		Movie{
			ID:    "1",
			Isbn:  "438227",
			Title: "Dark Knight",
			Director: &Director{
				FirstName: "Kristopher",
				LastName:  "Nolan",
			},
		}, Movie{
			ID:    "2",
			Isbn:  "45455",
			Title: "Harry Potter",
			Director: &Director{
				FirstName: "Joanne",
				LastName:  "Rowling",
			},
		})

	r.HandleFunc("/movies", getMovies).Methods(http.MethodGet)
	r.HandleFunc("/movies/{id}", getMovie).Methods(http.MethodGet)
	r.HandleFunc("/movies", createMovie).Methods(http.MethodPost)
	r.HandleFunc("/movies/{id}", updateMovie).Methods(http.MethodPut)
	r.HandleFunc("/movies/{id}", deleteMovie).Methods(http.MethodDelete)

	log.Printf("Starting server at port 8080\n")
	err := http.ListenAndServe(":8080", r)

	if err != nil {
		log.Fatalf("[ERROR] Can not start server. Error: %s\n", err)
	}
}
