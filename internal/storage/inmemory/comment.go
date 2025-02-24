package inmemory

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
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
	const f = "inmemory.CommentRepo.Create"
	r.Lock()
	defer r.Unlock()

	comment.CreatedAt = time.Now()
	comment.ID = len(r.comments) + 1
	r.comments[comment.ID] = *comment
	r.postComments[comment.PostID] = append(r.postComments[comment.PostID], comment.ID)
	if comment.ParentID != nil {
		r.repliers[*comment.ParentID] = append(r.repliers[*comment.ParentID], comment.ID)
	}
	logger.Info("Comment created", slog.String("func", f), slog.Int("commentID", comment.ID))
	return comment, nil
}

func (r *CommentRepo) GetByID(ctx context.Context, id int) (domain.Comment, bool) {
	const f = "inmemory.CommentRepo.GetByID"
	r.RLock()
	defer r.RUnlock()

	comment, exists := r.comments[id]
	if !exists {
		logger.Warn("Comment not found", slog.String("func", f), slog.Int("commentID", id))
	}
	return comment, exists
}

func (r *CommentRepo) GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error) {
	const f = "inmemory.CommentRepo.GetByPostID"
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
	logger.Info("Fetched comments by post ID", slog.String("func", f), slog.Int("postID", postID), slog.Int("count", len(ans)))
	return ans, nil
}

func (r *CommentRepo) GetReplies(ctx context.Context, commentID int, limit, offset int) ([]domain.Comment, error) {
	const f = "inmemory.CommentRepo.GetReplies"
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
	logger.Info("Fetched replies for comment", slog.String("func", f), slog.Int("commentID", commentID), slog.Int("count", len(ans)))
	return ans, nil
}

func (r *CommentRepo) CheckCommentUnderPost(ctx context.Context, postID, commentID int) (bool, error) {
	const f = "inmemory.CommentRepo.CheckCommentUnderPost"
	r.RLock()
	defer r.RUnlock()
	comment, exist := r.comments[commentID]
	if !exist {
		logger.Warn("No replies found for comment", slog.String("func", f), slog.Int("commentID", commentID))
		return false, nil
	}

	if comment.PostID == postID {
		return true, nil
	}

	return false, nil
}
