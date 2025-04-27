package postgres

import (
	"WebMessanger/internal/model"
	"database/sql"
	"github.com/google/uuid"
)

type likeRepo struct {
	db *sql.DB
}

func NewLikeRepo(db *sql.DB) *likeRepo {
	return &likeRepo{db: db}
}
func (r *likeRepo) AddLike(like *model.Like) error {
	_, err := r.db.Exec(`
		INSERT INTO likes (id, thread_id, user_id, user_name, created_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (thread_id, user_id) DO NOTHING
	`, like.ID, like.ThreadID, like.UserID, like.UserName, like.CreatedAt)

	return err
}

func (r *likeRepo) RemoveLike(threadID uuid.UUID, userID uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM likes WHERE thread_id = $1 AND user_id = $2`, threadID, userID)
	return err
}
func (r *likeRepo) GetLikesByThread(threadID uuid.UUID) ([]*model.Like, error) {
	rows, err := r.db.Query(`SELECT id, thread_id, user_id, user_name, created_at
		FROM likes
		WHERE thread_id = $1
		ORDER BY created_at ASC`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []*model.Like
	for rows.Next() {
		var l model.Like
		if err := rows.Scan(&l.ID, &l.ThreadID, &l.UserID, &l.UserName, &l.CreatedAt); err != nil {
			return nil, err
		}
		likes = append(likes, &l)
	}
	return likes, nil
}
func (r *likeRepo) GetLikesByUser(userID uuid.UUID) ([]*model.Like, error) {
	rows, err := r.db.Query(`SELECT id, thread_id, user_id, user_name, created_at FROM likes WHERE user_id = $1 ORDER BY  created_at ASC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []*model.Like
	for rows.Next() {
		var l model.Like
		if err := rows.Scan(&l.ID, &l.ThreadID, &l.UserID, &l.UserName, &l.CreatedAt); err != nil {
			return nil, err
		}
		likes = append(likes, &l)
	}
	return likes, nil
}
