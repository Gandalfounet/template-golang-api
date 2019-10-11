
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

/**
* POST
* Create an event : http://localhost:3000/event.
* Params :
* {
*	"ID": "2",
*	"Title": "Article 2",
*	"Description": "This article is sent through postman"
* }
*/
func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body) //Using ioutil package to convert request data into a human readable format (converted into a slice)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	
	json.Unmarshal(reqBody, &newEvent) // Unmarshal it to fit into our event struct
	events = append(events, newEvent) // Append it to the events slice
	w.WriteHeader(http.StatusCreated) // Code 201 => Successfuly created ?
	fmt.Println(newEvent)
	fmt.Println(events)
	json.NewEncoder(w).Encode(newEvent) // Send back response ?
}

/**
* GET
* Get One event : http://localhost:3000/event/1.
*/
func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

/**
* GET
* Get all events : http://localhost:3000/events
*/
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

/**
* PATCH
* Params : id @ GET
* New event @ POST
* {
*	"ID": "2",
*	"Title": "Article 2",
*	"Description": "This article is sent through postman"
* }
* Update an event : http://localhost:3000/event/1/update
*/
func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
			break
		}
	}
}

/**
* DELETE
* Delete an event : http://localhost:3000/event/delete/1
*/
func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello World! Welcome Home ^_^")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/event/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/event/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/event/{id}", deleteEvent).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":3000", router))
}

