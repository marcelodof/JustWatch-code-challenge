package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Movie Object
type Movie struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Director    string `json:"director"`
	Producer    string `json:"producer"`
	ReleaseDate string `json:"release_date"`
	RtScore     string `json:"rt_score"`
}

// Species Response Object
type SpeciesResponseObject struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Classification string   `json:"classification"`
	EyeColors      string   `json:"eye_colors"`
	HairColors     string   `json:"hair_colors"`
	URL            string   `json:"url"`
	People         []string `json:"people"`
	Films          []string `json:"films"`
}

// Movie Response Object
type MovieResponseObject struct {
	ID                     string   `json:"id"`
	Title                  string   `json:"title"`
	OriginalTitle          string   `json:"original_title"`
	OriginalTitleRomanised string   `json:"original_title_romanised"`
	Description            string   `json:"description"`
	Director               string   `json:"director"`
	Producer               string   `json:"producer"`
	ReleaseDate            string   `json:"release_date"`
	RunningTime            string   `json:"running_time"`
	RtScore                string   `json:"rt_score"`
	People                 []string `json:"people"`
	Species                []string `json:"species"`
	Locations              []string `json:"locations"`
	Vehicles               []string `json:"vehicles"`
	URL                    string   `json:"url"`
}

func queryAPI(url string) ([]byte, int) {
	/* Given an URL, queries it and returns the body. */

	res, err := http.Get(url)

	// Checking for errors in the GET request
	if err != nil {
		log.Println(err)
	}

	// Getting body from response
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	return data, res.StatusCode
}

func queryMovies(moviesURLs []string) []Movie {
	/* Given a slice of films, it queries the Ghibli API and return the slice of Movies. */

	var movies []Movie

	for _, movieURL := range moviesURLs {
		// Querying Ghibli API for each movie
		data, _ := queryAPI(movieURL)
		var movieResponseObject MovieResponseObject
		json.Unmarshal(data, &movieResponseObject)

		// Appeding movies as Movie struct
		movies = append(movies, Movie{
			ID:          movieResponseObject.ID,
			Title:       movieResponseObject.Title,
			Description: movieResponseObject.Description,
			Director:    movieResponseObject.Director,
			Producer:    movieResponseObject.Producer,
			ReleaseDate: movieResponseObject.ReleaseDate,
			RtScore:     movieResponseObject.RtScore,
		})
	}

	return movies
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	/* Given a species id as parameter, returns the list of movies that this species appears on. */
	w.Header().Set("Content-Type", "application/json")

	// Getting param species
	species, ok := r.URL.Query()["species"]
	if !ok || len(species) < 1 {
		msg := "Url Param 'species' is missing"
		log.Println(msg)
		json.NewEncoder(w).Encode(msg)
	}

	// Querying Ghibli API for species
	log.Println("Querying api for species", species[0])
	data, statusCode := queryAPI(fmt.Sprintf("https://ghibliapi.herokuapp.com/species/%s", species[0]))
	var speciesResponseObject SpeciesResponseObject
	json.Unmarshal(data, &speciesResponseObject)

	// Checking if there is the species
	if statusCode == 404 {
		msg := "Page not found"
		log.Println(msg)
		json.NewEncoder(w).Encode(msg)
	} else {
		// Querying all movies the species appears on.
		results := queryMovies(speciesResponseObject.Films)
		log.Println("Request was successfull")
		json.NewEncoder(w).Encode(results)
	}
}

func handleRequests() {
	/* Request Handler */
	http.HandleFunc("/movies", getMovies)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	/* Main Entry Point */
	handleRequests()
}
