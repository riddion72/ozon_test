package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.66

import (
	"github.com/riddion72/ozon_test/internal/graph/generated"
)

// // User is the resolver for the user field.
// func (r *commentResolver) User(ctx context.Context, obj *domain.Comment) (string, error) {
// 	panic(fmt.Errorf("not implemented: User - user"))
// }

// // PostID is the resolver for the postId field.
// func (r *commentResolver) PostID(ctx context.Context, obj *domain.Comment) (string, error) {
// 	panic(fmt.Errorf("not implemented: PostID - postId"))
// }

// // ParentID is the resolver for the parentId field.
// func (r *commentResolver) ParentID(ctx context.Context, obj *domain.Comment) (*string, error) {
// 	panic(fmt.Errorf("not implemented: ParentID - parentId"))
// }

// // CreatePost is the resolver for the createPost field.
// func (r *mutationResolver) CreatePost(ctx context.Context, input model.NewPost) (*domain.Post, error) {
// 	panic(fmt.Errorf("not implemented: CreatePost - createPost"))
// }

// // CreateComment is the resolver for the createComment field.
// func (r *mutationResolver) CreateComment(ctx context.Context, input model.NewComment) (*domain.Comment, error) {
// 	panic(fmt.Errorf("not implemented: CreateComment - createComment"))
// }

// User is the resolver for the user field.
// func (r *postResolver) User(ctx context.Context, obj *domain.Post) (string, error) {
// 	panic(fmt.Errorf("not implemented: User - user"))
// }

// // Posts is the resolver for the posts field.
// func (r *queryResolver) Posts(ctx context.Context, limit int, offset int) ([]domain.Post, error) {
// 	panic(fmt.Errorf("not implemented: Posts - posts"))
// }

// // Post is the resolver for the post field.
// func (r *queryResolver) Post(ctx context.Context, id string) (*domain.Post, error) {
// 	panic(fmt.Errorf("not implemented: Post - post"))
// }

// // Comments is the resolver for the comments field.
// func (r *queryResolver) Comments(ctx context.Context, postID string, limit int, offset int) ([]domain.Comment, error) {
// 	panic(fmt.Errorf("not implemented: Comments - comments"))
// }

// // CommentAdded is the resolver for the commentAdded field.
// func (r *subscriptionResolver) CommentAdded(ctx context.Context, postID string) (<-chan *domain.Comment, error) {
// 	panic(fmt.Errorf("not implemented: CommentAdded - commentAdded"))
// }

// Comment returns generated.CommentResolver implementation.
func (r *Resolver) Comment() generated.CommentResolver { return &commentResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Post returns generated.PostResolver implementation.
func (r *Resolver) Post() generated.PostResolver { return &postResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type commentResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type postResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
