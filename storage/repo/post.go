package repo

import (
    pb "github.com/template-service/genproto"
)

//UserStorageI ...
type PostStorageI interface {
    CreatePost(*pb.Post) (*pb.Post, error)
    GetUserPosts(userID string) ([]*pb.Post, error)
}
