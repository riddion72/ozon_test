package domain

import "time"

type Comment struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	PostID    string    `json:"postId"`
	ParentID  *string   `json:"parentId"` // Указатель для nullable
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"createdAt"`
}
