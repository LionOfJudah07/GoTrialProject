package services

import (
	"database/sql"
	"time"

	"wechiye/backend/database"
	"wechiye/backend/models"
)

type UserService struct {
	db *database.DB
}

func NewUserService(db *database.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Get() (*models.User, error) {
	var u models.User
	row := s.db.QueryRow(`SELECT id, full_name, email, username, avatar, gender, education_level, occupation, 
		has_kids, kids_allowance_amount, kids_allowance_interval, created_at, updated_at 
		FROM users LIMIT 1`)
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Username, &u.Avatar, &u.Gender,
		&u.EducationLevel, &u.Occupation, &u.HasKids, &u.KidsAllowanceAmount,
		&u.KidsAllowanceInterval, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return s.createDefault()
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *UserService) createDefault() (*models.User, error) {
	now := time.Now()
	res, err := s.db.Exec(`INSERT INTO users (full_name, created_at, updated_at) VALUES ('', ?, ?)`, now, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.User{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s *UserService) Update(u *models.User) error {
	u.UpdatedAt = time.Now()
	_, err := s.db.Exec(`UPDATE users SET full_name=?, email=?, username=?, avatar=?, gender=?, 
		education_level=?, occupation=?, has_kids=?, kids_allowance_amount=?, kids_allowance_interval=?, 
		updated_at=? WHERE id=?`,
		u.FullName, u.Email, u.Username, u.Avatar, u.Gender, u.EducationLevel,
		u.Occupation, u.HasKids, u.KidsAllowanceAmount, u.KidsAllowanceInterval,
		u.UpdatedAt, u.ID)
	return err
}

func (s *UserService) UpdateAvatar(avatar string) error {
	_, err := s.db.Exec(`UPDATE users SET avatar=?, updated_at=? WHERE id=(SELECT id FROM users LIMIT 1)`,
		avatar, time.Now())
	return err
}