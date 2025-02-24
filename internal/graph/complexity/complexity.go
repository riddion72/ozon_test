package complexity

import (
	"github.com/riddion72/ozon_test/internal/graph/model"
)

type Complexity struct {
	Comment struct {
		CreatedAt func(childComplexity int) int
		ID        func(childComplexity int) int
		ParentID  func(childComplexity int) int
		PostID    func(childComplexity int) int
		Replies   func(childComplexity int, limit *int, offset *int) int
		Text      func(childComplexity int) int
		User      func(childComplexity int) int
	}

	Mutation struct {
		CloseCommentsPost func(childComplexity int, user string, id int, commentsAllowed bool) int
		CreateComment     func(childComplexity int, input model.NewComment) int
		CreatePost        func(childComplexity int, input model.NewPost) int
	}

	Post struct {
		Comments        func(childComplexity int, limit *int, offset *int) int
		CommentsAllowed func(childComplexity int) int
		Content         func(childComplexity int) int
		CreatedAt       func(childComplexity int) int
		ID              func(childComplexity int) int
		Title           func(childComplexity int) int
		User            func(childComplexity int) int
	}

	Query struct {
		Comments func(childComplexity int, postID int, limit *int, offset *int) int
		Post     func(childComplexity int, id int) int
		Posts    func(childComplexity int, limit int, offset int) int
		Replies  func(childComplexity int, commentID int, limit *int, offset *int) int
	}

	Subscription struct {
		CommentAdded func(childComplexity int, postID int) int
	}
}

// TO DO дописать все поля, настроить вес
func NewComplexity() *Complexity {
	c := &Complexity{}

	// logger.Debug("NewComplexity run")
	c.Comment.CreatedAt = func(childComplexity int) int {
		// logger.Debug("calculate complexity", slog.Int("c.Comment.CreatedAt", childComplexity+1))
		return 1 + childComplexity
	}

	c.Comment.ID = func(childComplexity int) int {
		// logger.Debug("calculate complexity", slog.Int("c.Comment.CreatedAt", childComplexity+1))
		return 1 + childComplexity
	}

	c.Comment.Replies = func(childComplexity int, limit *int, offset *int) int {
		// logger.Debug("calculate complexity", slog.Int("c.Comment.ID", childComplexity+2))
		return 2 + childComplexity
	}

	c.Mutation.CreateComment = func(childComplexity int, input model.NewComment) int {
		// logger.Debug("calculate complexity", slog.Int("c.Mutation.CreateComment", childComplexity+1))
		return 1 + childComplexity
	}

	c.Post.Comments = func(childComplexity int, limit *int, offset *int) int {
		// logger.Debug("calculate complexity", slog.Int("c.Post.Comments", childComplexity+2))
		return 2 + childComplexity
	}

	c.Query.Posts = func(childComplexity int, limit int, offset int) int {
		// logger.Debug("calculate complexity", slog.Int("c.Query.Posts", childComplexity+1))
		return 1 + childComplexity
	}

	c.Subscription.CommentAdded = func(childComplexity int, postID int) int {
		// logger.Debug("calculate complexity", slog.Int("c.Subscription.CommentAdded", childComplexity+1))
		return 1 + childComplexity
	}

	return c
}
