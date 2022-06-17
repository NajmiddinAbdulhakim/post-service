package service

import (
	"context"

	pb "github.com/NajmiddinAbdulhakim/post-service/genproto"
	l "github.com/NajmiddinAbdulhakim/post-service/pkg/logger"
	cl "github.com/NajmiddinAbdulhakim/post-service/service/grpc_client"
	"github.com/NajmiddinAbdulhakim/post-service/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//UserService ...
type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client  cl.GrpcClientI
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
		return nil, status.Error(codes.Internal, `failed while inserting post`)
	}
	return post, nil
}

func (s *PostService) UpdatePost(ctx context.Context, req *pb.Post) (*pb.BoolResponse, error) {
	post, err := s.storage.Post().UpdatePost(req)
	if err != nil {
		s.logger.Error(`failed while update post`, l.Error(err))
		return nil, status.Error(codes.Internal, `failed while update post`)
	}
	return &pb.BoolResponse{Success: post}, nil
}

func (s *PostService) DeletePost(ctx context.Context, req *pb.PostByIdReq) (*pb.BoolResponse, error) {
	res, err := s.storage.Post().DeletePost(req.PostId)
	if err != nil {
		s.logger.Error(`Filed while delete post`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while delete post`)
	}
	return &pb.BoolResponse{Success: res}, nil
}

func (s *PostService) GetUserPosts(ctx context.Context, req *pb.GetUserPostsReq) (*pb.GetUserPostsRes, error) {
	posts, err := s.storage.Post().GetUserPosts(req.UserId)
	if err != nil {
		s.logger.Error(`failed while getting post user by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `failed while getting post user by id`)
	}
	return &pb.GetUserPostsRes{Posts: posts}, nil
}

func (s *PostService) GetPostById(ctx context.Context, req *pb.PostByIdReq) (*pb.Post, error) {
	post, err := s.storage.Post().GetPostById(req.PostId)
	if err != nil {
		s.logger.Error(`Filed while getting post by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting post by id`)
	}
	return post, nil
}

func (s *PostService) GetAllPosts(ctx context.Context, req *pb.Empty) (*pb.GetAllPostsRes, error) {
	posts, err := s.storage.Post().GetAllPosts()
	if err != nil {
		s.logger.Error(`Filed while getting all posts`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting all posts`)
	}
	return &pb.GetAllPostsRes{Posts: posts}, nil
}

func (s *PostService) GetPostWithUser(ctx context.Context, req *pb.PostByIdReq) (*pb.GetPostWithUserRes, error) {
	post, err := s.storage.Post().GetPostById(req.PostId)
	if err != nil {
		s.logger.Error(`Filed while getting posts user by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting posts user by id`)
	}
	user, err := s.client.UserService().GetUserById(ctx, &pb.UserByIdReq{Id: post.UserId})
	if err != nil {
		s.logger.Error(`Filed while getting user by id`, l.Error(err))
		return nil, status.Error(codes.Internal, `Filed while getting user by id`)
	}

	return &pb.GetPostWithUserRes{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Post:      post,
	}, nil
}
