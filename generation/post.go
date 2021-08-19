package main

import "time"

type Post struct {
	PublishDate time.Time
	ShortTitle  string
	LongTitle   string
	Path        string
}
