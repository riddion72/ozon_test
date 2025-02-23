package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/storage/inmemory/tools"
)

type CommentRepo struct {
	sync.RWMutex
	comments     map[int]domain.Comment
	postComments map[int][]int
	repliers     map[int][]int
}

func NewCommentRepo() *CommentRepo {
	return &CommentRepo{
		comments:     make(map[int]domain.Comment),
		postComments: make(map[int][]int),
		repliers:     make(map[int][]int),
	}
}

func (r *CommentRepo) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	r.Lock()
	defer r.Unlock()

	comment.CreatedAt = time.Now()
	comment.ID = len(r.comments) + 1
	r.comments[comment.ID] = *comment
	r.postComments[comment.PostID] = append(r.postComments[comment.PostID], comment.ID)
	if comment.ParentID != nil {
		r.repliers[*comment.ParentID] = append(r.repliers[*comment.ParentID], comment.ID)
	}
	return comment, nil
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
	commentsSlice := make([]domain.Comment, 0, len(commentIDs))
	end := limit + offset
	for i := 0; i < len(commentIDs) && end > 0; i++ {
		coment, exist := r.comments[commentIDs[i]]
		if exist && coment.ParentID == nil {
			end--
			commentsSlice = append(commentsSlice, coment)
		}
	}
	ans := tools.Paginate(commentsSlice, limit, offset)
	return ans, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, commentID int, limit, offset int) ([]domain.Comment, error) {
	r.RLock()
	defer r.RUnlock()

	repliersIDs := r.repliers[commentID]
	if repliersIDs == nil {
		return nil, nil
	}

	comments := make([]domain.Comment, 0, len(repliersIDs))
	for _, id := range repliersIDs {
		comments = append(comments, r.comments[id])
	}
	ans := tools.Paginate(comments, limit, offset)
	return ans, nil
}

func (r *CommentRepo) CheckCommentUnderPost(ctx context.Context, postID, commentID int) (bool, error) {
	r.RLock()
	defer r.RUnlock()
	comment, exist := r.comments[commentID]
	if !exist {
		return false, nil
	}

	if comment.PostID == postID {
		return true, nil
	}

	return false, nil
}
