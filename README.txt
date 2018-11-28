This little program implements documents services:

GET	/documents	List of documents
GET	/documents/:id	Get documents with ID
POST	/documents	Creates a document,having the file as body data
DELETE	/documents/:id	Deletes a document

If you have issues when trying to compile it, please run the following commands:

go get -u github.com/golang/dep/cmd/dep
go get github.com/gorilla/mux

(I am assuming that GO is installed on your OS)

To test POST AND DELETE, I recommend Postman (https://www.getpostman.com/)

