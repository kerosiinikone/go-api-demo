package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/gorilla/mux"
)

type Storage interface {
	handleCreateItemAndLogin(http.ResponseWriter, *http.Request)
	handleGetAllItems(http.ResponseWriter, *http.Request)
}

type Db struct {
	db *dbx.DB
}

type User struct {
	Id, Name string
}

type LocalStorage struct{}

type Item struct {
	Name string `json:"name"`
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var response = map[string]string{
		"hello": "world",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (ls *LocalStorage) handleCreateItemAndLogin(w http.ResponseWriter, r *http.Request) {
	var item Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Validation Error", http.StatusUnprocessableEntity)
	}
	defer r.Body.Close()

	// Persist to a local storage

	token, err := CreateAndSignToken(item.Name)

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
	}

	w.Header().Add("Authorization", "Bearer " + token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)

}

func (ls *LocalStorage) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{
			Id: "1",
			Name: "New",
		},
		{
			Id: "2",
			Name: "New",
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (db *Db) handleCreateItemAndLogin(w http.ResponseWriter, r *http.Request) {
	var item Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Validation Error", http.StatusUnprocessableEntity)
	}
	defer r.Body.Close()

	if _, err := db.db.Insert("items", dbx.Params{"name": item.Name}).Execute(); err != nil {
		http.Error(w, "Error creating item", http.StatusInternalServerError)
	}

	token, err := CreateAndSignToken(item.Name)

	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
	}

	w.Header().Add("Authorization", "Bearer " + token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(item)

}

func (db *Db) handleGetAllItems(w http.ResponseWriter, r *http.Request) {
	q := db.db.Select().
		From("items")

	var users []User
	if err := q.All(&users); err != nil {
		log.Fatal(err.Error())
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func withJWTAuth(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if len(tokenString) == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"Error": "No token",
			})
			return 
		}

		splitToken := strings.Split(tokenString, " ")
		
		if _, err := ParseToken(splitToken[1]); err != nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"Error": "Invalid Token",
			})
			return 
		}

		f(w, r)
	}
}


func main() {
	os.Setenv("SECRET", "secret")
	db := CreateDBConnection()

	// Simply, db implements the functions defined in Storage.
	// Therefore also every single other db-like struct with same
	// implementations can be used as a "Storage" provider
	
	var storage Storage = &Db{
		db: db,
	}

	r := mux.NewRouter()
    r.HandleFunc("/items", storage.handleCreateItemAndLogin).Methods("POST")
    r.HandleFunc("/items", withJWTAuth(storage.handleGetAllItems)).Methods("GET")
    r.HandleFunc("/", handleMain)
	
	log.Fatal(http.ListenAndServe(":8000", r))
}
