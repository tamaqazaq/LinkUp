package postgres

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"database/sql"
	"github.com/google/uuid"
)

type threadRepo struct {
	db *sql.DB
}

func NewThreadRepo(db *sql.DB) usecase.ThreadRepository {
	return &threadRepo{db: db}
}

func (r *threadRepo) Create(thr *model.Thread) (*model.Thread, error) {
	_, err := r.db.Exec(`INSERT INTO threads (id, user_id, content, media_url, created_at) VALUES ($1, $2, $3, $4, $5)`,
		thr.ID, thr.UserID, thr.Content, thr.MediaURL, thr.CreatedAt)
	if err != nil {
		return nil, err
	}
	return thr, nil
}

func (r *threadRepo) Update(thr *model.Thread) (*model.Thread, error) {
	_, err := r.db.Exec(`UPDATE threads SET content = $1, media_url = $2 WHERE id = $3`, thr.Content, thr.MediaURL, thr.ID)
	if err != nil {
		return nil, err
	}
	return thr, nil
}

func (r *threadRepo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM threads WHERE id = $1`, id)
	return err
}

func (r *threadRepo) GetAllThreads() ([]*model.Thread, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.user_id, u.name,u.avatar_url, t.content, t.media_url, t.created_at
		FROM threads t
		JOIN users u ON t.user_id = u.id
		ORDER BY t.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*model.Thread
	for rows.Next() {
		var t model.Thread
		err := rows.Scan(&t.ID, &t.UserID, &t.UserName, &t.AvatarURL, &t.Content, &t.MediaURL, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		threads = append(threads, &t)
	}
	return threads, nil
}

func (r *threadRepo) GetThreadById(id uuid.UUID) (*model.Thread, error) {
	var t model.Thread
	err := r.db.QueryRow(`
		SELECT t.id, t.user_id, u.name, u.avatar_url, t.content, t.media_url, t.created_at
		FROM threads t
		JOIN users u ON t.user_id = u.id
		WHERE t.id = $1
	`, id).
		Scan(&t.ID, &t.UserID, &t.UserName, &t.AvatarURL, &t.Content, &t.MediaURL, &t.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *threadRepo) GetByUser(userID uuid.UUID) ([]*model.Thread, error) {
	rows, err := r.db.Query(`
		SELECT t.id, t.user_id, u.name, u.avatar_url, t.content, t.media_url, t.created_at
		FROM threads t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = $1
		ORDER BY t.created_at DESC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var threads []*model.Thread
	for rows.Next() {
		var t model.Thread
		err := rows.Scan(&t.ID, &t.UserID, &t.UserName, &t.AvatarURL, &t.Content, &t.MediaURL, &t.CreatedAt)
		if err != nil {
			return nil, err
		}
		threads = append(threads, &t)
	}
	return threads, nil
}
