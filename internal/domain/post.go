package domain

import "time"

type Post struct {
	ID              int        `json:"id"`
	Title           string     `json:"title"`
	User            string     `json:"user"`
	Content         string     `json:"content"`
	CommentsAllowed bool       `json:"commentsAllowed"`
	CreatedAt       time.Time  `json:"createdAt"`
	Comments        []*Comment `json:"comments"`
}
