package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/dchooyc/book"
)

func main() {
	allBooks := retrieveFile("output.json")
	fmt.Println("Length of all books: ", len(allBooks))

	genreToBooks := sortByGenre(allBooks)
	createJsons(genreToBooks)
}

func createJsons(genreToBooks map[string][]book.Book) {
	for genre := range genreToBooks {
		books := book.Books{Books: genreToBooks[genre]}
		filename := "jsons/" + genre + ".json"

		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating json for: ", genre)
			continue
		}

		sort.Slice(books.Books, func(i, j int) bool {
			return books.Books[i].Ratings > books.Books[j].Ratings
		})

		if len(books.Books) > 1000 {
			books.Books = books.Books[:1000]
		}

		jsonData, err := json.Marshal(books)
		if err != nil {
			fmt.Println("Error marshalling json for: ", genre)
			continue
		}

		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing to json for: ", genre)
			continue
		}

		fmt.Println("Created: ", filename, " with length: ", len(books.Books))
	}
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
