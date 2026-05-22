package services

import (
	"time"

	"wechiye/backend/database"
	"wechiye/backend/models"
)

type BudgetService struct {
	db *database.DB
}

func NewBudgetService(db *database.DB) *BudgetService {
	return &BudgetService{db: db}
}

func (s *BudgetService) Set(b *models.Budget) (*models.Budget, error) {
	now := time.Now()
	res, err := s.db.Exec(`INSERT INTO budgets (category, month_year, amount_limit, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?) 
		ON CONFLICT(category, month_year) DO UPDATE SET amount_limit=excluded.amount_limit, updated_at=excluded.updated_at`,
		b.Category, b.MonthYear, b.AmountLimit, now, now)
	if err != nil {
		return nil, err
	}
	if b.ID == 0 {
		id, _ := res.LastInsertId()
		b.ID = id
	}
	b.CreatedAt = now
	b.UpdatedAt = now
	return b, nil
}

func (s *BudgetService) List(monthYear string) ([]*models.Budget, error) {
	rows, err := s.db.Query(`SELECT id, category, month_year, amount_limit, created_at, updated_at 
		FROM budgets WHERE month_year = ?`, monthYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var budgets []*models.Budget
	for rows.Next() {
		var b models.Budget
		err := rows.Scan(&b.ID, &b.Category, &b.MonthYear, &b.AmountLimit, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, &b)
	}
	return budgets, nil
}

func (s *BudgetService) Delete(id int64) error {
	_, err := s.db.Exec(`DELETE FROM budgets WHERE id=?`, id)
	return err
}