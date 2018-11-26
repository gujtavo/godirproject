package main

import (
	"log"
	"net/http"

	"encoding/json"
	
	"io/ioutil"
	"os"

	"github.com/gorilla/mux"
)

var filesPath = "./folder"

type Document struct {
	Id   string
	Name string
	Size int
}

func getFiles(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	info := getFiles(filesPath)
	json.NewEncoder(w).Encode(info)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}
