package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type DateConcert struct {
	ID   int                 `json:"id"`
	Date map[string][]string `json:"datesLocations"`
}

type Resultats struct {
	Position string
	date     []string
}

func main() {
	baseURL := "https://groupietrackers.herokuapp.com/api/relation/"

	maxID := 52 // Nombres maximum d'IDs à tester

	for id := 1; id <= maxID; id++ {
		url := fmt.Sprintf(baseURL+"%d", id) // Construire l'URL avec l'ID
		response, err := http.Get(url)       // Effectuer la requête GET
		if err != nil {
			fmt.Printf("Erreur lors de la requête : %v\n", err)
			return
		}
		defer response.Body.Close()

		// Vérifier le code HTTP
		if response.StatusCode != http.StatusOK {
			fmt.Printf("Code HTTP inattendu pour l'ID %d : %d\n", id, response.StatusCode)
			continue
		}

		// Décoder la réponse JSON
		var data DateConcert
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			fmt.Printf("Erreur lors du décodage JSON pour l'ID %d : %v\n", id, err)
			continue
		}

		// Afficher les données de manière lisible
		fmt.Printf("Données reçues pour l'ID %d :\n", id)
		fmt.Printf("- ID : %d\n", data.ID)
		fmt.Println()

		var resultats []Resultats
		for location, dates := range data.Date {
			resultat := Resultats{
				Position: location,
				date:     dates,
			}
			resultats = append(resultats, resultat)
		}
		for _, res := range resultats {
			fmt.Printf("%s %v\n", res.Position, res.date)
		}
	}
}
