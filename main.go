package main

import (
	"bufio"
	"fmt"
	"os"
)

// DocumentID uniquely identifies a document by its filepath.
type DocumentID string

// DocumentIDs is a set of document IDs.
type DocumentIDs map[string]struct{}

// Document represents a document that was read from a file.
type Document string

// DocTermFrequency captures how often each term occurs in a given document.
type DocTermFrequency map[string]int

// CollectionTermFrequency captures how often each term occurs in each document of a given collection.
// It maps filenames to DocTermFrequency maps.
type CollectionTermFrequency map[string]DocTermFrequency

// SearchEngine represents the search index.
type SearchEngine map[string]DocumentIDs

func AddDocument(engine SearchEngine, document Document) {
	panic("not implemented")
}

func IndexLookup(engine SearchEngine, term string) []DocumentID {
	panic("not implemented")
}

// TermFrequency implements the term frequency td(t, d) for a given term and document.
func TermFrequency(term string, document string) float64 {
	panic("not implemented")
}

// InverseDocumentFrequency implements the inverse document frequency idf(t, D) for a given term and set of Documents.
func InverseDocumentFrequency(term string, document []Document) float64 {
	panic("not implemented")
}

func RelevanceLookup(term string, engine SearchEngine) []DocumentID {
	panic("not implemented")
}

func FindFiles(directory string) ([]DocumentID, error) {
	entries, err := os.ReadDir(directory)

	documents := make([]DocumentID, 0)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		entryPath := directory + string(os.PathSeparator) + entry.Name()
		if entry.IsDir() {
			subdirEntries, err := FindFiles(entryPath)
			if err != nil {
				return nil, err
			}
			documents = append(documents, subdirEntries...)
		} else {
			documents = append(documents, DocumentID(entryPath))
		}
	}

	return documents, nil
}

func ReadDirectory(directory string) error {
	files, err := FindFiles(directory)
	if err != nil {
		return err
	}
	fmt.Println(files)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing required command line argument DIRECTORY. Aborting.")
		os.Exit(1)
	}

	directory := os.Args[1]

	err := ReadDirectory(directory)
	if err != nil {
		fmt.Println("Failed to read the directory: ", err)
		os.Exit(1)
	}

	// Input read loop by the example of https://stackoverflow.com/a/49715256.
	cliReader := bufio.NewScanner(os.Stdin)
	for {
		cliReader.Scan()
		term := cliReader.Text()
		if len(term) > 0 {
			fmt.Println("== " + term)
			// TODO: Search for the term and output documents.
		} else {
			break
		}
	}
}
