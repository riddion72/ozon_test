package inmemory

import (
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
)

type CommentRepo struct {
	sync.RWMutex
	comments     map[string]domain.Comment
	postComments map[string][]string
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments:     make(map[string]domain.Comment),
		postComments: make(map[string][]string),
	}
}

func (r *CommentRepo) Create(comment domain.Comment) error {
	r.Lock()
	defer r.Unlock()

	comment.CreatedAt = time.Now()
	r.comments[comment.ID] = comment
	r.postComments[comment.PostID] = append(r.postComments[comment.PostID], comment.ID)
	return nil
}

func (r *CommentRepo) GetByPostID(postID string, limit, offset int) []domain.Comment {
	r.RLock()
	defer r.RUnlock()

	commentIDs := r.postComments[postID]
	comments := make([]domain.Comment, 0, len(commentIDs))
	for _, id := range commentIDs {
		comments = append(comments, r.comments[id])
	}
	return paginate(comments, limit, offset)
}

func (r *CommentRepo) GetByID(id string) (domain.Comment, bool) {
	r.RLock()
	defer r.RUnlock()

	comment, exists := r.comments[id]
	return comment, exists
}
