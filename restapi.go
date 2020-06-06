package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Books struct {
	Books []Book `json:"books"`
}

type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

type Author struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

// Local books store
// TODO: persistence
var books []Book

/**
 * Retrieves all books
 * TODO: paginate
 */
func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

/**
 * Retrieves a single book by its unique ID (from json param)
 */
func getBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, book := range books {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	json.NewEncoder(w).Encode(&Book{})
}

/**
 * Creates a new book (from json params)
 */
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book) // decode request body data directly into book struct
	book.ID = uuid.New().String()
	books = append(books, book)
	json.NewEncoder(w).Encode(book)
}

/**
 * Updates a book by its unique ID (from json param)
 */
func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newbook Book
	_ = json.NewDecoder(r.Body).Decode(&newbook) // decode request body data directly into book struct
	for i, book := range books {
		if book.ID == newbook.ID {
			books = append(books[:i], books[i+1:]...)
			books = append(books, newbook)
			break
		}
	}
	json.NewEncoder(w).Encode(newbook)
}

/**
 * Deletes a book by its unique ID (from json param)
 */
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, book := range books {
		if book.ID == params["id"] {
			books = append(books[:i], books[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(books)
}

func loadFixture(url string) {
	// temp data
	data, err := ioutil.ReadFile(url)
	if err != nil {
		fmt.Print(err)
	}
	var fixtureBooks Books
	err = json.Unmarshal(data, &fixtureBooks)
	if err != nil {
		fmt.Println("error:", err)
	}
	for _, book := range fixtureBooks.Books {
		books = append(books, book)
	}
}

func main() {
	loadFixture("./fixtures/books.json")

	// Init
	router := mux.NewRouter()

	// Endpoints
	// TODO: non-polymorphic endpoints?
	router.HandleFunc("/api/books", getBooks).Methods("GET")
	router.HandleFunc("/api/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/api/books", createBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", updateBook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", deleteBook).Methods("DELETE")

	// Start up
	log.Fatal(http.ListenAndServe(":7007", router))
	fmt.Println("Serving on localhost:7007")
}

// TODO: validation
