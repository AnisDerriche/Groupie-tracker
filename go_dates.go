package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Dates struct {
	ID   int                 `json:"id"`
	Date map[string][]string `json:"dates"`
}

type Result struct {
	date []string
}

func go_artist() {
	baseURL := "https://groupietrackers.herokuapp.com/api//dates/"

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
		var data Dates
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			fmt.Printf("Erreur lors du décodage JSON pour l'ID %d : %v\n", id, err)
			continue
		}

		// Afficher les données de manière lisible
		fmt.Printf("Données reçues pour l'ID %d :\n", id)
		fmt.Printf("- ID : %d\n", data.ID)
		fmt.Println()

		var Resultats []Result
		for dates := range data.Date {
			Resultats := Result{
				date: dates,
			}
			Resultats = append(Resultats, Result)
		}
		for _, res := range Resultats {
			fmt.Printf("%s %v\n", res.date)
		}
	}
}
