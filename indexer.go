package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
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
type SearchEngine map[string]DocTermFrequency

// DocumentFrequencyMapping is a tuple of a DocumentID and the DocTermFrequency associated with that document.
type DocumentFrequencyMapping struct {
	document  DocumentID
	frequency DocTermFrequency
}

var WordRegex = regexp.MustCompile("[a-zA-Z]+(['-][a-zA-Z]+)*")

// DocumentsContaining returns the set of documents that contain the given term.
func (se *SearchEngine) DocumentsContaining(term string) []DocumentID {
	panic("not implemented")
}

// RelevantDocuments returns a list of documents relevant for the given term, ordered from most relevant to least relevant.
func (se *SearchEngine) RelevantDocuments(term string) []DocumentID {
	panic("not implemented")
}

// InverseDocumentFrequency implements the inverse document frequency idf(t, D) for a given term and set of Documents.
func InverseDocumentFrequency(term string, document []Document) float64 {
	panic("not implemented")
}

// TermFrequency implements the term frequency td(t, d) for a given term and document.
func TermFrequency(term string, document string) float64 {
	panic("not implemented")
}

// ReduceDocuments is the reducer function of the implemented reducer pattern.
// It reads documents from a channel and creates a search engine, that is passed back through another channel.
func ReduceDocuments(documents chan DocumentFrequencyMapping, output chan SearchEngine) {
	searchEngine := SearchEngine{}
	defer func() { output <- searchEngine }()
	for document := range documents {
		searchEngine[string(document.document)] = document.frequency
	}
}

// Frequencies calculates the term frequency for a given document.
// Reads the file from disk using the given DocumentID (= file path), performs all text processing operations, and finally writes the result to the channel.
// Errors are printed to STDERR, but not communicated to any other part of the program. The WaitGroup (orchestrated by the caller) will ensure that there are no deadlocks.
func Frequencies(document DocumentID, ch chan DocumentFrequencyMapping) {
	// LABEL ReadFile
	// Read file into array of lines.
	file, err := os.Open(string(document))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading %s: %s", document, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// LABEL Split into words and count them.
	var wordCounts DocTermFrequency = make(map[string]int)
	for _, line := range lines {
		words := WordRegex.FindAllString(line, -1)
		for _, word := range words {
			lowercaseWord := strings.ToLower(word)
			if count, ok := wordCounts[lowercaseWord]; ok {
				wordCounts[lowercaseWord] = count + 1
			} else {
				wordCounts[lowercaseWord] = 1
			}
		}
	}

	// LABEL PublishDocumentFrequencyMapping
	// Push the computed mapping to the channel.
	ch <- DocumentFrequencyMapping{document, wordCounts}
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

func ReadDirectory(directory string) (SearchEngine, error) {
	files, err := FindFiles(directory)
	if err != nil {
		return nil, err
	}
	// XXX: Right now, we read all files synchronously and then dump them all in a channel.
	// We could also write files into a channel directly and start processing while we're still searching for files,
	// but we're skipping that for now because finding the files should be very fast (there aren't a lot I guess?)
	// and it reduces complexity a bit.

	// LABEL CreateDocTermFreqencyChannelAndGoroutines
	// Create a channel and launch a goroutine for each file, writing results to the channel.
	chFrequencies := make(chan DocumentFrequencyMapping)
	wgFrequencies := new(sync.WaitGroup)
	for _, file := range files {
		wgFrequencies.Go(func() {
			Frequencies(file, chFrequencies)
		})
	}

	// LABEL InitReducer
	// Initialize a goroutine to read from the channel and aggregate everything into a SearchEngine object.
	chSearchEngine := make(chan SearchEngine)
	go ReduceDocuments(chFrequencies, chSearchEngine)

	wgFrequencies.Wait()
	close(chFrequencies)

	searchEngine := <-chSearchEngine
	return searchEngine, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing required command line argument DIRECTORY. Aborting.")
		os.Exit(1)
	}

	directory := os.Args[1]

	searchEngine, err := ReadDirectory(directory)
	if err != nil {
		fmt.Println("Failed to read the directory: ", err)
		os.Exit(1)
	}

	fmt.Println(searchEngine)

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
