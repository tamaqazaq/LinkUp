package postgres

import (
	"WebMessanger/internal/model"
	"database/sql"
	"github.com/google/uuid"
)

type commentRepo struct {
	db *sql.DB
}

func NewCommentRepo(db *sql.DB) *commentRepo {
	return &commentRepo{db: db}
}
func (r *commentRepo) Create(cmt *model.Comment) (*model.Comment, error) {
	_, err := r.db.Exec(`INSERT INTO comments (id, thread_id, user_id, user_name, content, created_at) VALUES ($1, $2, $3, $4, $5, $6)`, cmt.ID, cmt.ThreadID, cmt.UserID, cmt.UserName, cmt.Content, cmt.CreatedAt)
	if err != nil {
		return nil, err
	}
	return cmt, nil
}
func (r *commentRepo) GetByThread(threadID uuid.UUID) ([]*model.Comment, error) {
	rows, err := r.db.Query(`
		SELECT id, thread_id, user_id, user_name, content, created_at
		FROM comments
		WHERE thread_id = $1
		ORDER BY created_at ASC
	`, threadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.Comment
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.ID, &c.ThreadID, &c.UserID, &c.UserName, &c.Content, &c.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}

	return comments, nil
}
func (r *commentRepo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`DELETE FROM comments WHERE id = $1`, id)
	return err
}
