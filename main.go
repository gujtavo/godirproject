package main

import (
	"log"
	"net/http"

	"encoding/json"
	
	"io/ioutil"
	"os"
	"fmt"

	"github.com/gorilla/mux"
	"crypto/md5"
	"io"
	"encoding/hex"
	"time"
)

var filesPath = "./folder"

type Document struct {
	Id   string
	Name string
	Size int64
	Modified time.Time
}

func getFiles(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return files
}

func addmd5(filesInfo []os.FileInfo) []Document {
	var files []Document
	
	for _, f := range filesInfo {
		MD5, err := hashFileMD5(filesPath + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, Document{Id: MD5, Name: f.Name(), Size: f.Size(), Modified: f.ModTime()})
	}
	return files
}

// https://www.mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
func hashFileMD5(filePath string) (string, error) {
	// Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	// Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	// Tell the program to call the following function when the current function returns
	defer file.Close()

	// Open a new hash interface to write to
	hash := md5.New()

	// Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	// Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	// Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil
}

func getDocuments(w http.ResponseWriter, r *http.Request) {
	info := getFiles(filesPath)
	md5 := addmd5(info)
	json.NewEncoder(w).Encode(md5)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	log.Fatal(http.ListenAndServe(":9000", router))
}
