package inmemory_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/riddion72/ozon_test/internal/domain"
	"github.com/riddion72/ozon_test/internal/logger"
	"github.com/riddion72/ozon_test/internal/storage/inmemory"
)

func TestCommentRepo_Create(t *testing.T) {
	ctx := context.Background()
	logger.MustInit("prod")

	t.Run("Create comment", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()
		comment := &domain.Comment{
			User:     "test user",
			PostID:   1,
			ParentID: nil,
			Text:     "comment",
		}

		created, err := repo.Create(ctx, comment)
		require.NoError(t, err)
		assert.Equal(t, 1, created.ID)
		assert.Equal(t, "test user", created.User)
		assert.Equal(t, 1, created.PostID)
		assert.Nil(t, created.ParentID)
		assert.Equal(t, "comment", created.Text)
		assert.WithinDuration(t, time.Now(), created.CreatedAt, 1*time.Second)

	})

	t.Run("Create comment with comment", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		parent := &domain.Comment{PostID: 1, Text: "Parent"}
		parent, err := repo.Create(ctx, parent)
		require.NoError(t, err)

		child := &domain.Comment{
			PostID:   1,
			Text:     "Child",
			ParentID: &parent.ID,
		}

		created, err := repo.Create(ctx, child)
		require.NoError(t, err)

		assert.Equal(t, 2, created.ID)
		assert.Equal(t, *created.ParentID, parent.ID)
		assert.Equal(t, "Child", created.Text)
	})

	t.Run("Create comments", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		for i := 1; i <= 10; i++ {
			c := &domain.Comment{PostID: 1}
			created, err := repo.Create(ctx, c)
			require.NoError(t, err)
			assert.Equal(t, i, created.ID)
		}

	})

	t.Run("Comment added to correct post and repliers", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()
		// Post 1:
		// - Comment 1
		//   - Comment 2
		// - Comment 3
		// Post 2:
		// - Comment 4

		c1 := &domain.Comment{PostID: 1}
		c1, err := repo.Create(ctx, c1)
		require.NoError(t, err)

		c2 := &domain.Comment{PostID: 1, ParentID: &c1.ID}
		c2, err = repo.Create(ctx, c2)
		require.NoError(t, err)

		c3 := &domain.Comment{PostID: 1}
		c3, err = repo.Create(ctx, c3)
		require.NoError(t, err)

		c4 := &domain.Comment{PostID: 2}
		c4, err = repo.Create(ctx, c4)
		require.NoError(t, err)

	})
}

func TestGetByID(t *testing.T) {
	ctx := context.Background()

	t.Run("exist comment", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		created, err := repo.Create(ctx, &domain.Comment{
			PostID: 1,
			Text:   "comment",
		})
		require.NoError(t, err)

		comment, exists := repo.GetByID(ctx, created.ID)

		assert.True(t, exists, "Should exist")
		assert.Equal(t, created.ID, comment.ID)
		assert.Equal(t, "comment", comment.Text)
	})

	t.Run("!exist comment", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		comment, exists := repo.GetByID(ctx, 123)

		assert.False(t, exists, "Should not exist")
		assert.Empty(t, comment, "Comment should be empty")
	})

	t.Run("Verify all fields", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		parent, err := repo.Create(ctx, &domain.Comment{
			User:     "test user",
			PostID:   1,
			ParentID: nil,
			Text:     "comment",
		})
		require.NoError(t, err)

		testComment := &domain.Comment{
			User:     "test user",
			PostID:   2,
			Text:     "Child",
			ParentID: &parent.ID,
		}
		created, err := repo.Create(ctx, testComment)
		require.NoError(t, err)

		comment, exists := repo.GetByID(ctx, created.ID)
		require.True(t, exists, "Comment should exist")

		assert.Equal(t, created.ID, comment.ID)
		assert.Equal(t, 2, comment.PostID)
		assert.Equal(t, "Child", comment.Text)
		require.NotNil(t, comment.ParentID)
		assert.Equal(t, parent.ID, *comment.ParentID)
		assert.WithinDuration(t, time.Now(), comment.CreatedAt, 10*time.Millisecond)
	})

	t.Run("zero ID", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()
		comment, exists := repo.GetByID(ctx, -0)
		assert.False(t, exists)
		assert.Empty(t, comment)
	})

	t.Run("negative ID", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()
		comment, exists := repo.GetByID(ctx, -1)
		assert.False(t, exists)
		assert.Empty(t, comment)
	})
}

func TestGetByPostID(t *testing.T) {
	ctx := context.Background()

	t.Run("comments without parent", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		_, err := repo.Create(ctx, &domain.Comment{PostID: 1, Text: "Root 1"})
		require.NoError(t, err)
		_, err = repo.Create(ctx, &domain.Comment{PostID: 1, Text: "Reply", ParentID: ptr(1)})
		require.NoError(t, err)
		_, err = repo.Create(ctx, &domain.Comment{PostID: 1, Text: "Root 2"})
		require.NoError(t, err)

		comments, err := repo.GetByPostID(ctx, 1, 10, 0)
		require.NoError(t, err)

		assert.Len(t, comments, 2, "Should return only root comments")
		assert.Nil(t, comments[0].ParentID)
		assert.Nil(t, comments[1].ParentID)
		assert.ElementsMatch(t, []string{"Root 1", "Root 2"}, []string{
			comments[0].Text,
			comments[1].Text,
		})
	})

	t.Run("Pagination", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		for i := 1; i <= 5; i++ {
			_, err := repo.Create(ctx, &domain.Comment{
				PostID: 1,
				Text:   "Comment " + string(rune('A'+i-1)),
			})
			require.NoError(t, err)
		}

		tests := []struct {
			name     string
			limit    int
			offset   int
			expected []string
		}{
			{"First page", 2, 0, []string{"Comment A", "Comment B"}},
			{"Second page", 2, 2, []string{"Comment C", "Comment D"}},
			{"Partial page", 1, 4, []string{"Comment E"}},
			{"Offset beyond data", 5, 10, []string(nil)},
			{"Zero limit", 0, 0, []string(nil)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				comments, err := repo.GetByPostID(ctx, 1, tt.limit, tt.offset)
				require.NoError(t, err)

				var Texts []string
				for _, c := range comments {
					Texts = append(Texts, c.Text)
				}
				assert.Equal(t, tt.expected, Texts)
			})
		}
	})

	t.Run("Empty result", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		comments, err := repo.GetByPostID(ctx, 123, 10, 0)
		require.NoError(t, err)
		assert.Empty(t, comments)
	})

	t.Run("Ignore ParentID", func(t *testing.T) {
		repo := inmemory.NewCommentRepo()

		// - Root comment (ID 1)
		//   - Reply 1 (ID 2)
		//   - Reply 2 (ID 3)
		// - Root (ID 4)
		root1, _ := repo.Create(ctx, &domain.Comment{PostID: 1})
		repo.Create(ctx, &domain.Comment{PostID: 1, ParentID: &root1.ID})
		repo.Create(ctx, &domain.Comment{PostID: 1, ParentID: &root1.ID})
		repo.Create(ctx, &domain.Comment{PostID: 1})

		comments, err := repo.GetByPostID(ctx, 1, 10, 0)
		require.NoError(t, err)

		assert.Len(t, comments, 2, "only root comments")
		for _, c := range comments {
			assert.Nil(t, c.ParentID)
		}
	})
}

func TestGetReplies(t *testing.T) {
	repo := inmemory.NewCommentRepo()
	ctx := context.Background()

	parentComment := &domain.Comment{
		PostID: 1,
		Text:   "Parent comment",
	}
	repo.Create(ctx, parentComment)

	replyComment := &domain.Comment{
		PostID:   1,
		Text:     "Reply comment",
		ParentID: &parentComment.ID,
	}
	repo.Create(ctx, replyComment)

	replies, err := repo.GetReplies(ctx, parentComment.ID, 10, 0)
	if err != nil {
		t.Fatalf("GetReplies() error = %v", err)
	}

	if len(replies) != 1 {
		t.Errorf("GetReplies() got %v replies, want %v", len(replies), 1)
	}

	if replies[0].Text != "Reply comment" {
		t.Errorf("GetReplies() got Text = %v, want %v", replies[0].Text, "Reply comment")
	}
}

func TestCheckCommentUnderPost(t *testing.T) {
	repo := inmemory.NewCommentRepo()
	ctx := context.Background()

	comment := &domain.Comment{
		PostID: 1,
		Text:   "Test comment",
	}
	repo.Create(ctx, comment)

	exists, err := repo.CheckCommentUnderPost(ctx, 1, comment.ID)
	if err != nil {
		t.Fatalf("CheckCommentUnderPost() error = %v", err)
	}

	if !exists {
		t.Error("CheckCommentUnderPost() comment not found under post")
	}
}

func ptr(i int) *int { return &i }
