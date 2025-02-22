package domain

import "time"

type Comment struct {
	ID        int        `json:"id"`
	User      string     `json:"user"`
	PostID    int        `json:"postId"`
	ParentID  *int       `json:"parentId,omitempty"` // Указатель для nullable
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"createdAt"`
	Replies   []*Comment `json:"replies,omitempty"`
}
