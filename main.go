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
	"path/filepath"
	
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(md5)
}

func getDocumentsById(w http.ResponseWriter, r *http.Request) { //for example http://localhost:9000/documents/1593371f4e183af281e34375
	vars := mux.Vars(r)
	id := vars["id"]	
	
	info := getFiles(filesPath)
	md5 := addmd5(info)
	Result := Document{}
	for _, n := range md5 {
		if id == n.Id {
			Result = n
		}
	}

	if Result.Id != "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Result)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found!"))
	}


}

func deleteDocumentById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	status := 0
	info := getFiles(filesPath)
	var files []Document
	
	for _, f := range info {
		MD5, err := hashFileMD5(filesPath + "/" + f.Name())
		if err != nil {
			log.Fatal(err)
		}
		files = append(files, Document{Id: MD5, Name: f.Name(), Size: f.Size(), Modified: f.ModTime()})
		if MD5 == id {
			os.Remove(filesPath + "/" + f.Name())
			status = 200
		}
	}

	if status != 200 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Delete was not succesful"))

	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("deleted!"))
	}

}

	
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func createDocument(w http.ResponseWriter, r *http.Request) {
	fileType := r.PostFormValue("type")	//make sure a text key named "key" is sent
	file, _, err := r.FormFile("file")  //make sure a file key named "file" must be sent and a file must be attached (it will be the value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("no file in form-data"))
		return
	}


	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileName := "fdfdfd3343d" //just a sample name
	
	newPath := filepath.Join(filesPath, fileName+"."+fileType)
	var _, error = os.Stat(newPath)

	if os.IsNotExist(error) { 
		newFile, err := os.Create(newPath)
		check(err)
		if _, err := newFile.Write(fileBytes); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("error when creating the file"))
			return
		}
	  	defer newFile.Close()
	} else
	{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("file already exists"))
		return
	}
	

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File created!"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/documents", getDocuments).Methods("GET")
	router.HandleFunc("/documents/{id}", getDocumentsById).Methods("GET")
	router.HandleFunc("/documents", createDocument).Methods("POST")
	router.HandleFunc("/documents/{id}", deleteDocumentById).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":9000", router))
}
