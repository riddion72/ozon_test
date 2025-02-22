package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
)

type CommentRepo struct {
	sync.RWMutex
	comments     map[int]domain.Comment
	postComments map[int][]int
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments:     make(map[int]domain.Comment),
		postComments: make(map[int][]int),
	}
}

func (r *CommentRepo) Create(ctx context.Context, comment domain.Comment) error {
	r.Lock()
	defer r.Unlock()

	comment.CreatedAt = time.Now()
	r.comments[comment.ID] = comment
	r.postComments[comment.PostID] = append(r.postComments[comment.PostID], comment.ID)
	return nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (domain.Comment, bool) {
	r.RLock()
	defer r.RUnlock()

	comment, exists := r.comments[id]
	return comment, exists
}

func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error) {
	r.RLock()
	defer r.RUnlock()

	commentIDs := r.postComments[postID]
	comments := make([]domain.Comment, 0, len(commentIDs))
	for _, id := range commentIDs {
		comments = append(comments, r.comments[id])
	}
	ans := paginate(comments, limit, offset)
	return ans, nil
}
func (r *CommentRepo) GetReplies(ctx context.Context, commentID int, limit, offset int) ([]domain.Comment, error) {
	r.RLock()
	defer r.RUnlock()

	parentID := r.comments[commentID].ParentID
	if parentID == nil {
		return nil, nil
	}

	commentIDs := r.postComments[*parentID]
	comments := make([]domain.Comment, 0, len(commentIDs))
	for _, id := range commentIDs {
		comments = append(comments, r.comments[id])
	}
	ans := paginate(comments, limit, offset)
	return ans, nil
}
