package search

import "github.com/webx-top/echo"

type Searcher interface {
	Add(index string, primaryKey *string, docs ...interface{}) error
	Update(index string, primaryKey *string, docs ...interface{}) error
	Delete(index string, ids ...string) error
	Flush() error
	InitIndex(cfg *IndexConfig) error
	Search(index string, keywords string, options *SearchRequest) (int64, []echo.H, error)
}

// DefaultSearch is the default search engine implementation.
var DefaultSearch = &NopSearch{}

// NopSearch is a no-op implementation of the Search interface.
type NopSearch struct{}

// Add implements the search engine interface by doing nothing (no-op pattern)
// index: the name of the index to add documents to (ignored)
// primaryKey: the primary key identifier (ignored)
// docs: documents to add (ignored)
// Returns nil as this is a no-op implementation
func (n *NopSearch) Add(index string, primaryKey *string, docs ...interface{}) error {
	return nil
}

// Update is a no-op implementation of the search engine update method.
// It takes an index name, optional primary key, and documents to update,
// but performs no actual operation and always returns nil.
func (m *NopSearch) Update(index string, primaryKey *string, docs ...interface{}) error {
	return nil
}

// Delete is a no-op implementation that does nothing and always returns nil.
// It satisfies the Search interface requirement for deleting documents from an index.
func (m *NopSearch) Delete(index string, ids ...string) error {
	return nil
}

// InitIndex initializes the search index with the given configuration.
// This is a no-op implementation that always returns nil.
func (m *NopSearch) InitIndex(cfg *IndexConfig) error {
	return nil
}

// Flush is a no-op implementation of the Search interface's Flush method.
func (n *NopSearch) Flush() error {
	return nil
}

// Search implements a no-operation search that returns empty results.
// It returns zero count, nil results, and nil error for any input.
func (n *NopSearch) Search(index string, keywords string, options *SearchRequest) (int64, []echo.H, error) {
	return 0, nil, nil
}
