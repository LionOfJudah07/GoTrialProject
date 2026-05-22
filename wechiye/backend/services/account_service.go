package services

import (
	"time"

	"wechiye/backend/database"
	"wechiye/backend/models"
)

type AccountService struct {
	db *database.DB
}

func NewAccountService(db *database.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) Create(name string, initialBalance float64) (*models.Account, error) {
	now := time.Now()
	res, err := s.db.Exec(`INSERT INTO accounts (name, initial_balance, current_balance, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`, name, initialBalance, initialBalance, now, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return &models.Account{
		ID:             id,
		Name:           name,
		InitialBalance: initialBalance,
		CurrentBalance: initialBalance,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

func (s *AccountService) Update(id int64, name string, initialBalance float64) error {
	_, err := s.db.Exec(`UPDATE accounts SET name=?, updated_at=? WHERE id=?`,
		name, time.Now(), id)
	return err
}

func (s *AccountService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM accounts WHERE id=?`, id)
	return err
}

func (s *AccountService) List() ([]*models.Account, error) {
	rows, err := s.db.Query(`SELECT id, name, initial_balance, current_balance, created_at, updated_at FROM accounts ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var accounts []*models.Account
	for rows.Next() {
		var a models.Account
		err := rows.Scan(&a.ID, &a.Name, &a.InitialBalance, &a.CurrentBalance, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, &a)
	}
	return accounts, nil
}