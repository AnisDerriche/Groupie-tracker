package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type artist struct {
	ID    int          `json:"id"`
	Name  string       `json:"name"`
	Image template.URL `json:"image"`
}

type relation struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
	ArtistName     string
	ArtistImage    template.URL
}

// J'ai ajoouté une partie pour filtré les données reçues car sinon le JS ne marchait pas et ajouté ça dans le handler
// (nouveau type struct pour transmettre l'ID correct)
type FilteredRelation struct {
	ID             int                 `json:"id"`
	ArtistID       int                 `json:"artist_id"`
	ArtistName     string              `json:"artist_name"`
	ArtistImage    template.URL        `json:"artist_image"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

type ApiResponse struct {
	Index []relation `json:"index"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	relations, err := getrelation()
	if err != nil {
		http.Error(w, "Erreur de récupération des données", http.StatusInternalServerError)
		return
	}

	artists, err := getartist()
	if err != nil {
		http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
		return
	}

	filter := strings.ToLower(r.URL.Query().Get("filter"))
	var filtered []FilteredRelation

	fmt.Println("Filtrage demandé :", filter)
	for _, artist := range artists {
		if filter == "" || strings.Contains(strings.ToLower(artist.Name), filter) {
			for _, relation := range relations {
				if relation.ID == artist.ID {
					filtered = append(filtered, FilteredRelation{
						ID:             relation.ID,
						ArtistID:       artist.ID,
						ArtistName:     artist.Name,
						ArtistImage:    artist.Image,
						DatesLocations: relation.DatesLocations,
					})
				}
			}
		}
	}

	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filtered)
		return
	}

	tmpl, err := template.ParseGlob("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement des pages", http.StatusInternalServerError)
		return
	}
	err = tmpl.ExecuteTemplate(w, "index.html", filtered)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	artists, err := getartist()
	if err != nil {
		http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
		return
	}

	relations, err := getrelation()
	if err != nil {
		http.Error(w, "Erreur de récupération des données", http.StatusInternalServerError)
		return
	}

	var selectedArtist *artist
	var selectedRelation *relation

	for _, a := range artists {
		if a.ID == id {
			selectedArtist = &a
			break
		}
	}

	for _, r := range relations {
		if r.ID == id {
			selectedRelation = &r
			break
		}
	}

	if selectedArtist == nil {
		http.Error(w, "Artiste non trouvé", http.StatusNotFound)
		return
	}

	data := struct {
		Name           string
		Image          template.URL
		DatesLocations map[string][]string
	}{
		Name:           selectedArtist.Name,
		Image:          selectedArtist.Image,
		DatesLocations: selectedRelation.DatesLocations,
	}

	tmpl, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page artiste", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}

func artistConcertsHandler(w http.ResponseWriter, r *http.Request) {
	artistID := r.URL.Query().Get("id")
	if artistID == "" {
		http.Error(w, "ID de l'artiste manquant", http.StatusBadRequest)
		return
	}

	// Récupérer les relations des concerts et filtrer par artisteID
	relations, err := getrelation()
	if err != nil {
		http.Error(w, "Erreur de récupération des données des concerts", http.StatusInternalServerError)
		return
	}

	var artistConcerts []map[string][]string

	for _, relation := range relations {
		if strconv.Itoa(relation.ID) == artistID {
			// Récupérer les villes pour cet artiste
			for date, locations := range relation.DatesLocations {
				artistConcerts = append(artistConcerts, map[string][]string{
					"date":      {date},
					"locations": locations,
				})
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(artistConcerts)
}

func getCityCoordinates(city string) (float64, float64) {
	// Utilise l'API Nominatim pour obtenir les coordonnées
	url := "https://nominatim.openstreetmap.org/search?q=" + city + "&format=json&addressdetails=1"
	response, err := http.Get(url)
	if err != nil {
		return 48.8566, 2.3522 // Retourne les coordonnées de Paris en cas d'erreur
	}
	defer response.Body.Close()

	var data []struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil || len(data) == 0 {
		return 48.8566, 2.3522 // Retourne les coordonnées de Paris si la ville n'est pas trouvée
	}

	lat, err := strconv.ParseFloat(data[0].Lat, 64)
	if err != nil {
		lat = 48.8566
	}

	lon, err := strconv.ParseFloat(data[0].Lon, 64)
	if err != nil {
		lon = 2.3522
	}

	return lat, lon
}
func googleHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	date := r.URL.Query().Get("date")

	if city == "" || date == "" {
		http.Error(w, "Paramètres manquants", http.StatusBadRequest)
		return
	}

	// Utiliser une API pour récupérer les coordonnées de la ville
	lat, lon := getCityCoordinates(city)

	// Renvoyer les paramètres à google.html
	data := struct {
		City string
		Date string
		Lat  float64
		Lon  float64
	}{
		City: city,
		Date: date,
		Lat:  lat,
		Lon:  lon,
	}

	tmpl, err := template.ParseFiles("templates/google.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page Google", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur lors de l'exécution du template", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/artist.html", artistHandler)
	http.HandleFunc("/filter", handler)
	http.HandleFunc("/artistConcerts", artistConcertsHandler)
	http.HandleFunc("/google.html", googleHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	port := ":8098"
	fmt.Println("Serveur démarré sur http://localhost" + port)
	http.ListenAndServe(port, nil)
}

func getartist() ([]artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists"

	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erreur lors de la requête : %v\n", err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var artists []artist
	err = json.Unmarshal(body, &artists)
	if err != nil {
		return nil, err
	}

	return artists, nil
}

func getrelation() ([]relation, error) {
	url := "https://groupietrackers.herokuapp.com/api/relation"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse ApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}
	return apiResponse.Index, nil
}
