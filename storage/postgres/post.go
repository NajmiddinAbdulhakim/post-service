package postgres

import (
	"fmt"
	"time"


	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	pb "github.com/template-service/genproto"
)

type postRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(post *pb.Post) (*pb.Post, error) {
	tx, err := r.db.Begin()
	if err != nil {
		tx.Rollback()
		return post , err 
	}
	
	postID, err := uuid.NewV4()
	time := time.Now()
	var posted pb.Post
	
	if err != nil {
		tx.Rollback()
		return nil, err 
	}
	post.Id = postID.String()
	query := `INSERT INTO posts (id, user_id, title, description, created_at)
	VALUES ($1, $2, $3, $4,$5) RETURNING id, title, description, created_at`
	err = tx.QueryRow(query, post.Id, post.UserId, post.Title, post.Description, time).Scan(	
		&posted.Id,
		&posted.Title, 
		&posted.Description,
		&posted.CreatedAt,
	)
	if err != nil{
		tx.Rollback()
		return nil, fmt.Errorf(`error insert post > %v`,err) 
	}
	
	var me []*pb.Media
	for _, mediaa := range post.Medias {
		mediaID, err := uuid.NewV4()
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf(`error gen uuid by media > %v`,err) 
		}
		// mediaa.Id = mediaID.String()
		mQuery := `INSERT INTO media (id, post_id, link, type) 
		VALUES ($1, $2, $3, $4) RETURNING id, post_id, link, type`
		
		var media pb.Media
		
		err = tx.QueryRow(mQuery, mediaID, posted.Id, mediaa.Link, mediaa.Type).Scan(
			&media.Id,&media.PostId, &media.Link, &media.Type)
		if err != nil{
			tx.Rollback()
			return nil, fmt.Errorf(`error insert media > %v`,err) 
		}	
		me = append(me,&media)
	}
	posted.Medias = me
	tx.Commit()
	return &posted, nil 

}

func (r *postRepo) GetPostById(postID string) (*pb.Post, error) {
	query := `SELECT id, user_id, title, description FROM posts WHERE id = $1`
	rowss, err := r.db.Query(query, postID)
	if err != nil {
		return nil, err
	}

	var post pb.Post
	for rowss.Next() {
		err := rowss.Scan(
			&post.Id,
			&post.UserId,
			&post.Title, 
			&post.Description,
		)
		if err != nil {
			return nil, err
		}
		queryM := `SELECT id, link, type FROM media WHERE post_id = $1`
		rows, err := r.db.Query(queryM, post.Id)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err = rows.Scan(
				&media.Id,
				&media.Link, 
				&media.Type,
			)
			if err != nil {
				return nil, err
			}
			post.Medias = append(post.Medias, &media)
		}
	}
	return &post, nil
}

func (r *postRepo) UpdatePost(post *pb.Post) (bool, error) {
	time := time.Now()
	query := `UPDATE posts SET title = $1, description = $2, updated_at = $3 WHERE id = $4 AND user_id = $5`
	_, err := r.db.Exec(query, post.Title, post.Description,time, post.Id, post.UserId)
	if err != nil {
		return false, err
	}
	for _, media := range post.Medias {
		queryM := `UPDATE media SET link = $1, type = $2 WHERE id = $3 AND post_id = $4`
		_, err := r.db.Exec(queryM, media.Link, media.Type, media.Id, media.PostId)
		if err != nil {
			return false, err
		}
	}
	return true, nil

}

func (r *postRepo) DeletePost(postID string) (bool,error) {
	time := time.Now()

	query := `UPDATE posts SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time, postID)
	if err != nil {
		return false, err
	}
	return true, nil

}

func (r *postRepo) GetAllPosts() ([]*pb.Post, error) {

	var posts []*pb.Post
	
	query := `SELECT id, user_id, title, description FROM posts`
	rowss, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}

	var post pb.Post
	for rowss.Next() {
		err := rowss.Scan(
			&post.Id,
			&post.UserId,
			&post.Title, 
			&post.Description,
		)
		if err != nil {
			return nil, err
		}
		queryM := `SELECT id, link, type FROM media`
		rows, err := r.db.Query(queryM)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err = rows.Scan(
				&media.Id,
				&media.Link, 
				&media.Type,
			)
			if err != nil {
				return nil, err
			}
			post.Medias = append(post.Medias, &media)
		}
		posts = append(posts, &post)
	}
	return posts, nil
	
}

func (r *postRepo) GetUserPosts(userID string) ([]*pb.Post, error) {
	var posts []*pb.Post

	query := `SELECT id, title, description FROM posts WHERE user_id = $1`
	rowss, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	// defer rows.Close()
	for rowss.Next() {
		var post pb.Post
		err := rowss.Scan(
			&post.Id,
			&post.Title, 
			&post.Description,
		)
		if err != nil {
			return nil, err
		}
		queryM := `SELECT id, link, type FROM media WHERE post_id = $1`
		rows, err := r.db.Query(queryM, post.Id)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err = rows.Scan(
				&media.Id,
				&media.Link, 
				&media.Type,
			)
			if err != nil {
				return nil, err
			}
			post.Medias = append(post.Medias, &media)
		}
		posts = append(posts, &post)
	}
	return posts, nil
}




