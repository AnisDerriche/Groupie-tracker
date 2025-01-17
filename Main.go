package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type album struct {
	Index interface{} "json:'index'"
}

func main() {
	url := "https://groupietrackers.herokuapp.com/api/relation" //lien de l'api

	// geter
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("Erreur lors de la requête : %v\n", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK { // verifie si il n'y a pas d'erreur dans le lien
		fmt.Printf("Code HTTP inattendu : %d\n", response.StatusCode)
		return
	}
	println(response)

	var data album
	err = json.NewDecoder(response.Body).Decode(&data) // permet de lire le le json et de le décoder
	if err != nil {                                    //erreur si il y a
		fmt.Printf("Erreur lors du décodage JSON : %v\n", err)
		return
	}

	fmt.Printf("Données reçues : %+v\n", data)
}
