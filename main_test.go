package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	type Response struct {
		Hello string
	}

	var expectedRes = Response{
		Hello: "world",
	}

	r := mux.NewRouter()
    r.HandleFunc("/", handleMain)
	
	testServer := httptest.NewServer(r)
	defer testServer.Close()
	
	res, err := http.Get(testServer.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var resp Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Error("Error decoding response body:", err)
	}

	if resp.Hello != expectedRes.Hello {
		t.Errorf("Expected '%s', got '%s'", expectedRes.Hello, resp.Hello)
	}
}

func TestCreateItem_NoMock(t *testing.T) {
	type Response struct {
		Name string
	}
	expectedRes := Response{
		Name: "New",
	}

	var storage Storage = &LocalStorage{} 

	r := mux.NewRouter()
    r.HandleFunc("/items", storage.handleCreateItemAndLogin).Methods("POST")

	testServer := httptest.NewServer(r)
	defer testServer.Close()

	body := []byte(`{"name": "New"}`)
	
	res, err := http.Post(testServer.URL + "/items", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var resp Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Error("Error decoding response body:", err)
	}

	if res.Header.Get("Authorization") == "" {
		t.Errorf("Expected Authorization Header")
	}

	if resp.Name != expectedRes.Name {
		t.Errorf("Expected '%s', got '%s'", expectedRes.Name, resp.Name)
	}

}

func TestGetAllItemsWithJWT_NoMock(t *testing.T) {
	expectedRes := []User{
			{
				Id: "1",
				Name: "New",
			},
			{
				Id: "2",
				Name: "New",
			},
		}
	

	var storage Storage = &LocalStorage{} 

	r := mux.NewRouter()
    r.HandleFunc("/items", withJWTAuth(storage.handleGetAllItems)).Methods("GET")

	testServer := httptest.NewServer(r)
	defer testServer.Close()

	token, err := CreateAndSignToken("")

	if err != nil {
		t.Fatal(err)
	}
	
	client := &http.Client{}
	req, err := http.NewRequest("GET", testServer.URL + "/items", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	
	if err != nil {
		t.Fatal(err)
	}

	res, err := client.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var resp []User
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		t.Error("Error decoding response body:", err)
	}

	assert.Equal(t, expectedRes, resp)  
}
