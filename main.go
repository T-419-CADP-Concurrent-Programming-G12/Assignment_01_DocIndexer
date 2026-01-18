package main

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

func main() {
}
