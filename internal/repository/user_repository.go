package repository

import (
	"database/sql"
	"errors"
	"lab2-terrylneal/internal/models"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// INSERT
func (r *UserRepository) Create(user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	query := `
        INSERT INTO users (username, email, first_name, last_name, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	user.IsActive = true

	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return errors.New("username already exists")
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return errors.New("email already exists")
		}
		return err
	}

	return nil
}

// GetByID
func (r *UserRepository) GetByID(id int) (*models.User, error) {
	user := &models.User{}

	query := `
        SELECT id, username, email, first_name, last_name, is_active, created_at, updated_at
        FROM users 
        WHERE id = $1
    `

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetAll
func (r *UserRepository) GetAll(onlyActive bool) ([]*models.User, error) {
	var query string
	var rows *sql.Rows
	var err error

	if onlyActive {
		query = `
            SELECT id, username, email, first_name, last_name, is_active, created_at, updated_at
            FROM users 
            WHERE is_active = true
            ORDER BY id
        `
		rows, err = r.db.Query(query)
	} else {
		query = `
            SELECT id, username, email, first_name, last_name, is_active, created_at, updated_at
            FROM users 
            ORDER BY id
        `
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UPDATE
func (r *UserRepository) Update(user *models.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	query := `
        UPDATE users 
        SET username = $1, email = $2, first_name = $3, last_name = $4, 
            is_active = $5, updated_at = $6
        WHERE id = $7
        RETURNING id
    `

	user.UpdatedAt = time.Now()

	var id int
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.IsActive,
		user.UpdatedAt,
		user.ID,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return errors.New("user not found")
	}
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return errors.New("username already exists")
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return errors.New("email already exists")
		}
		return err
	}

	return nil
}

// PATCH partial update
func (r *UserRepository) PartialUpdate(id int, updates map[string]interface{}) (*models.User, error) {
	existingUser, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	query := "UPDATE users SET updated_at = $1"
	args := []interface{}{time.Now()}
	argPos := 2

	if username, ok := updates["username"]; ok {
		query += ", username = $" + string(rune('0'+argPos))
		args = append(args, username)
		argPos++
		existingUser.Username = username.(string)
	}

	if email, ok := updates["email"]; ok {
		query += ", email = $" + string(rune('0'+argPos))
		args = append(args, email)
		argPos++
		existingUser.Email = email.(string)
	}

	if firstName, ok := updates["first_name"]; ok {
		query += ", first_name = $" + string(rune('0'+argPos))
		args = append(args, firstName)
		argPos++
		existingUser.FirstName = firstName.(string)
	}

	if lastName, ok := updates["last_name"]; ok {
		query += ", last_name = $" + string(rune('0'+argPos))
		args = append(args, lastName)
		argPos++
		existingUser.LastName = lastName.(string)
	}

	if isActive, ok := updates["is_active"]; ok {
		query += ", is_active = $" + string(rune('0'+argPos))
		args = append(args, isActive)
		argPos++
		existingUser.IsActive = isActive.(bool)
	}

	if err := existingUser.ValidatePartial(); err != nil {
		return nil, err
	}

	query += " WHERE id = $" + string(rune('0'+argPos)) + " RETURNING id"
	args = append(args, id)

	var returnedID int
	err = r.db.QueryRow(query, args...).Scan(&returnedID)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			return nil, errors.New("username already exists")
		}
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return nil, errors.New("email already exists")
		}
		return nil, err
	}

	return r.GetByID(id)
}

// DELETE
func (r *UserRepository) Delete(id int, hardDelete bool) error {
	var query string
	var result sql.Result
	var err error

	if hardDelete {
		// permanently remove from db
		query = `DELETE FROM users WHERE id = $1`
		result, err = r.db.Exec(query, id)
	} else {
		// marking as inactive (soft delete)
		query = `UPDATE users SET is_active = false, updated_at = $1 WHERE id = $2`
		result, err = r.db.Exec(query, time.Now(), id)
	}

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) PermanentDelete(id int) error {
	return r.Delete(id, true)
}
