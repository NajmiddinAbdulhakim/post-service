package repo

import (
    pb "github.com/template-service/genproto"
)

//UserStorageI ...
type PostStorageI interface {
    CreatePost(*pb.Post) (*pb.Post, error)
    GetPostById(postID string) (*pb.Post, error)
    UpdatePost(post *pb.Post) (bool, error)
    DeletePost(PostID string) (bool, error)
    GetAllPosts() ([]*pb.Post, error)
    GetUserPosts(userID string) ([]*pb.Post, error)
}
