package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	http.HandleFunc("/people", peopleFunc)
	http.HandleFunc("/people/", findPerson)
	http.ListenAndServe(":8080", nil)
}

type Person struct {
	Name       string `json: "name"`
	Age        int    `json: "age"`
	Profession string `json: "profession"`
	HairColor  string `json: "hairColor"`
}

var peopleMap = make(map[string]Person)

func peopleFunc(w http.ResponseWriter, req *http.Request) {

	/*
		var _, err = json.Marshal(thing)
		if err != nil {
			fmt.Fprintf(w, "we have an error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	*/
	if req.Method == http.MethodGet {
		peopleStr, err := json.Marshal(peopleMap)
		if err == nil {
			fmt.Fprintf(w, "%s", string(peopleStr))

			for _, person := range peopleMap {
				var personJSON, _ = json.Marshal(person)
				f, err := os.OpenFile("testOutput.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				f.Write(personJSON)
				//err := ioutil.WriteFile("testOutput.txt", personJSON, 0644)
				if err != nil {
					fmt.Println("Error occurred while attempting to write to file")
				}
			}

		} else {
			fmt.Fprintf(w, "Failed to Marshal the map: %s", err)
		}
	} else if req.Method == http.MethodPost {
		b, err := ioutil.ReadAll(req.Body)
		if err == nil {
			var newPeople []Person
			json.Unmarshal(b, &newPeople)
			for _, newPerson := range newPeople {
				peopleMap[newPerson.Name] = newPerson
			}
		} else {
			fmt.Fprintf(w, "Error occurred while attempting to read POST request body")
		}

	}
	//fmt.Fprintf(w, "%s", req.URL.Path[1:])
}

func findPerson(w http.ResponseWriter, req *http.Request) {
	personName := req.URL.Path[8:]
	person, prs := peopleMap[personName]
	var personJSON, _ = json.Marshal(person)

	if prs {
		fmt.Fprintf(w, "%s", personJSON)
	} else {
		fmt.Fprintf(w, "no such person exists: %s", personName)
	}
}
