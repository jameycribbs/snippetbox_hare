package models

import (
	"sort"
	"strconv"
	"time"

	"github.com/jameycribbs/hare"
)

type Snippet struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Created time.Time `json:"created"`
	Expires time.Time `json:"expires"`
}

func (snippet *Snippet) GetID() int {
	return snippet.ID
}

func (snippet *Snippet) SetID(id int) {
	snippet.ID = id
}

func (snippet *Snippet) AfterFind() {
	*snippet = Snippet(*snippet)
}

type Snippets struct {
	*hare.Table
}

func NewSnippets(db *hare.Database) (*Snippets, error) {
	tbl, err := db.GetTable("snippets")
	if err != nil {
		return nil, err
	}

	return &Snippets{Table: tbl}, nil
}

func (snippets *Snippets) Query(queryFn func(snippet Snippet) bool, limit int) ([]Snippet, error) {
	var results []Snippet
	var err error

	for _, id := range snippets.IDs() {
		snippet := Snippet{}

		if err = snippets.Find(id, &snippet); err != nil {
			return nil, err
		}

		if queryFn(snippet) {
			results = append(results, snippet)
		}

		if limit != 0 && limit == len(results) {
			break
		}
	}

	return results, err
}

// This will return the 10 most recently created snippets.
func (snippets *Snippets) Latest() ([]*Snippet, error) {
	now := time.Now()

	results, err := snippets.Query(func(r Snippet) bool {
		return now.Before(r.Expires)
	}, 0)
	if err != nil {
		return nil, err
	}

	latest := make([]*Snippet, 0, 10)

	// Don't use for range because the fact that you are appending a pointer
	// of the rec in results to the latest array would mean that you would
	// be populating all entries in latest with the last looped rec from
	// results.
	for i := 0; i < len(results); i++ {
		if len(latest) < 10 {
			latest = append(latest, &results[i])
			continue
		}

		for j := 0; j < len(latest); j++ {
			if results[i].Created.After(latest[j].Created) {
				latest[j] = &results[i]
				break
			}
		}
	}

	// Reverse sort by creation time
	sort.Slice(latest, func(i, j int) bool {
		return latest[j].Created.Before(latest[i].Created)
	})

	return latest, nil
}

func (snippets *Snippets) Insert(title, content, expires string) (int, error) {
	expiresInt, err := strconv.Atoi(expires)
	if err != nil {
		return 0, err
	}

	snippet := Snippet{
		Title:   title,
		Content: content,
		Created: time.Now(),
		Expires: time.Now().AddDate(0, 0, expiresInt),
	}

	id, err := snippets.Create(&snippet)
	if err != nil {
		return 0, err
	}

	return id, nil
}
