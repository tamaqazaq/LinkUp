package postgres

import (
	"WebMessanger/internal/model"
	"WebMessanger/internal/usecase"
	"database/sql"
	"github.com/google/uuid"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) usecase.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *model.User) error {
	_, err := r.db.Exec(`
  INSERT INTO users (
    id, name, email, hashed_password,
    bio, location, is_online, social_links,
    avatar_url, created_at
  )
  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`,
		user.ID,
		user.Name,
		user.Email,
		user.HashedPassword,
		user.Bio,
		user.Location,
		false,
		user.SocialLinks,
		user.AvatarURL,
		user.CreatedAt,
	)

	return err
}

func (r *userRepo) GetByName(name string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow("SELECT id, name, email, hashed_password, is_verified FROM users WHERE name = $1", name).
		Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.IsVerified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(`
        SELECT
            id,
            name,
            email,
            bio,
            location,
            social_links,
            avatar_url,
            created_at, is_verified
        FROM users
        WHERE id = $1
    `, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Bio,
		&user.Location,
		&user.SocialLinks,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.IsVerified,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) UpdateProfile(user *model.User) error {
	_, err := r.db.Exec(`
  UPDATE users SET
    name = $1,
    bio = $2,
    location = $3,
    avatar_url = $4,
    social_links = $5,
    updated_at = NOW()
  WHERE id = $6
`,
		user.Name, user.Bio, user.Location, user.AvatarURL, user.SocialLinks, user.ID,
	)

	return err
}
func (r *userRepo) Search(query string) ([]*model.User, error) {
	rows, err := r.db.Query("SELECT id, name, avatar_url FROM users WHERE name ILIKE $1", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []*model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.AvatarURL) // ✅ два поля
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
func (r *userRepo) VerifyUserEmail(userID uuid.UUID) error {
	query := `UPDATE users SET is_verified = TRUE WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}
func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.QueryRow(`
		SELECT id, name, email, hashed_password, is_verified
		FROM users
		WHERE email = $1
	`, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.HashedPassword,
		&user.IsVerified,
	)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
