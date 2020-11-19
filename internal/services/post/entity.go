package post

import "time"

type Entity struct {
	UUID            string
	Title           string
	Slug            string
	ContentMarkdown string
	ContentHTML     string
	PublishedAt     *time.Time
}
