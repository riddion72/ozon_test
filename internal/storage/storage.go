package storage

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage/inmemory"
	"github.com/riddion72/ozon_test/internal/storage/postgres"
)

type PostStorage interface {
	Create(ctx context.Context, post *domain.Post) (*domain.Post, error)
	GetByID(ctx context.Context, postID int) (domain.Post, bool)
	List(ctx context.Context, limit, offset int) ([]domain.Post, error)
	CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error)
}

type CommentStorage interface {
	Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error)
	GetByID(ctx context.Context, commentID int) (domain.Comment, bool)
	GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error)
	GetReplies(ctx context.Context, commentID int, limit int, offset int) ([]domain.Comment, error)
	CheckCommentUnderPost(ctx context.Context, postID, commentID int) (bool, error)
}

type Storage struct {
	Post    *PostStorage
	Comment *CommentStorage
}

func NewStorage(cfg config.DB) *Storage {
	const f = "storage.NewStorage"
	var postRepo PostStorage
	var commentRepo CommentStorage
	var db *sqlx.DB

	// Выбор хранилища
	if cfg.Host != "" {
		// Подключение к PostgreSQL с повторами
		var err error
		db, err = postgres.ConnectWithRetries(cfg)
		if err != nil {
			logger.Error("Failed to connect to PostgreSQL: ", slog.String("func", f), slog.String("error", err.Error()))
			// Если к PostgreSQL подключиться не получилось используем In-memory реализацию
			postRepo = inmemory.NewPostRepo()
			commentRepo = inmemory.NewCommentRepo()
			logger.Info("Using in-memory storage", slog.String("func", f))
		} else {

			postRepo = postgres.NewPostRepository(db)
			commentRepo = postgres.NewCommentRepository(db)
		}
	} else {
		// In-memory реализация
		postRepo = inmemory.NewPostRepo()
		commentRepo = inmemory.NewCommentRepo()
		logger.Info("Using in-memory storage", slog.String("func", f))
	}

	return &Storage{Post: &postRepo, Comment: &commentRepo}
}
