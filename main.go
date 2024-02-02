package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/dchooyc/book"
)

type Genres struct {
	Genres []string `json:"genres"`
}

func main() {
	allBooks := retrieveFile("output.json")
	allBooks = removeDuplicates(allBooks)
	fmt.Println("Length of all books: ", len(allBooks))

	genreToBooks := sortByGenre(allBooks)
	createJsons(genreToBooks)
}

func removeDuplicates(books []book.Book) []book.Book {
	idToBook := make(map[string]book.Book)

	for i := 0; i < len(books); i++ {
		b := books[i]
		idToBook[b.ID] = b
	}

	books = []book.Book{}

	for id := range idToBook {
		books = append(books, idToBook[id])
	}

	return books
}

func createJsons(genreToBooks map[string][]book.Book) {
	genres := Genres{}
	err := os.Mkdir("./jsons", 0755)
	if err != nil {
		fmt.Println("failed creating dir: ", err)
		return
	}

	for genre := range genreToBooks {
		if len(genre) == 0 || strings.Contains(genre, "%") {
			continue
		}

		books := book.Books{Books: genreToBooks[genre]}
		filename := "jsons/" + genre + ".json"

		err := createJsonBooks(filename, books)
		if err != nil {
			fmt.Println(genre, err)
			continue
		}

		genres.Genres = append(genres.Genres, genre)
	}

	err = createJsonGenres("genres.json", genres)
	if err != nil {
		fmt.Println("failed creating genres json: ", err)
	}
}

func createJsonGenres(filename string, genres Genres) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed creating json: %w", err)
	}

	sort.Strings(genres.Genres)

	jsonData, err := json.Marshal(genres)
	if err != nil {
		return fmt.Errorf("failed marshalling json: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed writing json: %w", err)
	}

	return nil
}

func createJsonBooks(filename string, books book.Books) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed creating json: %w", err)
	}

	sort.Slice(books.Books, func(i, j int) bool {
		return books.Books[i].Ratings > books.Books[j].Ratings
	})

	if len(books.Books) > 1000 {
		books.Books = books.Books[:1000]
	}

	jsonData, err := json.Marshal(books)
	if err != nil {
		return fmt.Errorf("failed marshalling json: %w", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("failed writing json: %w", err)
	}

	fmt.Println("Created: ", filename, " with length: ", len(books.Books))
	return nil
}

func sortByGenre(books []book.Book) map[string][]book.Book {
	genreToBooks := make(map[string][]book.Book)

	for i := 0; i < len(books); i++ {
		book := books[i]

		for _, genre := range book.Genres {
			genreToBooks[genre] = append(genreToBooks[genre], book)
		}
	}

	return genreToBooks
}

func retrieveFile(target string) []book.Book {
	file, err := os.Open(target)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil
	}

	var books book.Books
	err = json.Unmarshal(bytes, &books)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil
	}

	return books.Books
}
