package service

import (
	"context"

	"github.com/jmoiron/sqlx"
	pb "github.com/template-service/genproto"
	l "github.com/template-service/pkg/logger"
	"github.com/template-service/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	cl "github.com/template-service/service/grpc_client"

)

//UserService ...
type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client cl.GrpcClientI
}

//NewUserService ...
func NewPostService(db *sqlx.DB, log l.Logger, client cl.GrpcClientI) *PostService {
	return &PostService{
		storage: storage.NewStoragePg(db),
		logger:  log,
		client:  client,
	}
}


func (s *PostService) CreatePost(ctx context.Context, req *pb.Post) (*pb.Post, error) {
	post, err := s.storage.Post().CreatePost(req)
	if err != nil {
		s.logger.Error(`failed while inserting post`, l.Error(err))
		return nil, status.Error(codes.Internal,`failed while inserting post` )	
	}
	return post, nil
}

func (s *PostService) GetUserPosts(ctx context.Context, req *pb.GetUserPostsReq) (*pb.GetUserPostsRes, error) {
	posts, err := s.storage.Post().GetUserPosts(req.UserId)
	if err != nil {
		s.logger.Error(`failed while getting post user by id`, l.Error(err))
		return nil, status.Error(codes.Internal,`failed while getting post user by id` )	
	}
	return &pb.GetUserPostsRes{Posts: posts}, nil
}

func (s *PostService) GetPostById(ctx context.Context, req *pb.GetPostReq) (*pb.Post, error) {
	post, err := s.storage.Post().GetPostById(req.PostId)
	if err != nil {
		s.logger.Error(`Filed while getting post by id`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while getting post by id` )
	}
	return post, nil
}

func (s *PostService) GetPostWithUser(ctx context.Context, req *pb.GetPostReq) (*pb.GetUserWithPostRes, error) {
	post, err := s.storage.Post().GetPostById(req.PostId)
	if err != nil {
		s.logger.Error(`Filed while getting posts user by id`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while getting posts user by id` )
	}
	user, err := s.client.UserService().GetUserById(ctx, &pb.UserByIdReq{Id: post.UserId })
	if err != nil {
		s.logger.Error(`Filed while getting user by id`, l.Error(err))
		return nil, status.Error(codes.Internal,`Filed while getting user by id` )
	}

	return &pb.GetUserWithPostRes{
		FirstName: user.FirstName,
		LastName: user.LastName,
		Post: post,
		}, nil
}

	// func (s *PostService) GetUserWithPost(ctx context.Context, req *pb.GetUserPostsReq) (*pb.GetUserWithPostRes, error) {
	// 	user, err := s.client.UserService().GetUserById(ctx, &pb.UserByIdReq{Id: req.UserId})
	// 	if err != nil {
	// 		s.logger.Error(`Filed while getting user by id`, l.Error(err))
	// 		return nil, status.Error(codes.Internal,`Filed while getting user by id` )
	// 	}
		
	// 	posts, err := s.storage.Post().GetUserPosts(user.Id)
	// 	if err != nil {
	// 		s.logger.Error(`Filed while getting posts user by id`, l.Error(err))
	// 		return nil, status.Error(codes.Internal,`Filed while getting posts user by id` )
	// 	}
	// 	user.Posts = posts
	
	// 	return &pb.GetUserWithPostRes{
	// 		FirstName: user.FirstName,
	// 		LastName: user.LastName,
	// 		Post: posts,
	// 		}, nil
	// }