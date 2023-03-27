package models

import "time"

type Issue struct {
	Title     string	`json:"title"`
	Created	  time.Time	`json:"created_at"`
	URL       string	`json:"html_url"`
}