package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
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

	artists, err := getartist()
	if err != nil {
		http.Error(w, "Erreur de récupération des artistes", http.StatusInternalServerError)
		return
	}

	// Récupérer le paramètre de recherche (si existant)
	filter := strings.ToLower(r.URL.Query().Get("filter"))
	var filtered []relation

	fmt.Println("Filtrage demandé :", filter)
	fmt.Println("Résultats trouvés :", filtered)
	for _, artist := range artists {
		if filter == "" || strings.Contains(strings.ToLower(artist.Name), filter) {
			for _, relation := range relations {
				if relation.ID == artist.ID {
					relation.ArtistName = artist.Name
					relation.ArtistImage = artist.Image
					filtered = append(filtered, relation)
				}
			}
		}
	}

	// Vérifie si c'est une requête AJAX (fetch)
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(filtered)
		return
	}

	// Sinon, on sert la page HTML complète
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Erreur lors du chargement de la page", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, filtered)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/filter", handler) // Ajoute cette ligne pour gérer le filtrage
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
