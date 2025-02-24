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

type ApiResponse struct {
	Index []relation `json:"index"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	relations, err := getrelation()

	if err != nil {
		http.Error(w, "Erreur de récupération des données", http.StatusInternalServerError)
		return
	}

	artist, err := getartist()
	if err != nil {
		http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
		return
	}

	filter := strings.ToLower(r.URL.Query().Get("filter"))
	var filtered []relation

	for _, artist := range artist {
		if strings.Contains(strings.ToLower(artist.Name), filter) {
			for _, relation := range relations {
				if relation.ID == artist.ID {
					relation.ArtistName = artist.Name
					relation.ArtistImage = artist.Image
					filtered = append(filtered, relation)
				}
			}
		}
	}
	tmpl, err := template.ParseGlob("templates/*.html")
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
	// Récupère l'ID de l'artiste dans l'URL
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

	// trouver l'artiste et ses information de concert
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

	// struct pour transmettre a artist.html
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

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/artist.html", artistHandler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	port := ":8080"
	fmt.Println("Serveur démarré sur http://localhost" + port)
	http.ListenAndServe(port, nil)
}

func getartist() ([]artist, error) {
	url := "https://groupietrackers.herokuapp.com/api/artists" //lien de l'api

	// geter
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
