package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/riddion72/ozon_test/internal/config"
	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage/inmemory"
	"github.com/riddion72/ozon_test/internal/storage/postgres"
)

type PostStorage interface {
	Create(ctx context.Context, post domain.Post) error
	GetByID(ctx context.Context, id int) (domain.Post, bool)
	List(ctx context.Context, limit, offset int) []domain.Post
	CommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) (*domain.Post, error)
}

type CommentStorage interface {
	Create(ctx context.Context, comment domain.Comment) error
	GetByID(ctx context.Context, id int) (domain.Comment, bool)
	GetByPostID(ctx context.Context, postID int, limit, offset int) ([]domain.Comment, error)
	GetReplies(ctx context.Context, commentID int, limit int, offset int) ([]domain.Comment, error)
}

type Storage struct {
	Post    *PostStorage
	Comment *CommentStorage
}

func NewStorage(post *PostStorage, comment *CommentStorage) *Storage {
	return &Storage{Post: post, Comment: comment}
}

func CreateStorages(cfg config.DB) (PostStorage, CommentStorage) {
	var postRepo PostStorage
	var commentRepo CommentStorage
	var db *sqlx.DB

	// Выбор хранилища
	if cfg.Host != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// Подключение к PostgreSQL с повторами
		var err error
		db, err = postgres.ConnectWithRetries(ctx, cfg)
		if err != nil {
			logger.Error("Failed to connect to PostgreSQL: ", slog.String("error", err.Error()))
			// Если к PostgreSQL подключиться не получилось используем In-memory реализацию
			postRepo = inmemory.NewPostRepo()
			commentRepo = inmemory.NewCommentRepo()
			logger.Info("Using in-memory storage")
		} else {
			defer db.Close()

			postRepo = postgres.NewPostRepository(db)
			commentRepo = postgres.NewCommentRepository(db)
		}
	} else {
		// In-memory реализация
		postRepo = inmemory.NewPostRepo()
		commentRepo = inmemory.NewCommentRepo()
		logger.Info("Using in-memory storage")
	}

	return postRepo, commentRepo
}
