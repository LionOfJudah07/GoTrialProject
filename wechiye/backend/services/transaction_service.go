package services

import (
	"time"

	"wechiye/backend/database"
	"wechiye/backend/models"
)

type TransactionService struct {
	db *database.DB
}

func NewTransactionService(db *database.DB) *TransactionService {
	return &TransactionService{db: db}
}

func (s *TransactionService) Create(txn *models.Transaction) (*models.Transaction, error) {
	now := time.Now()
	txn.CreatedAt = now
	txn.UpdatedAt = now

	res, err := s.db.Exec(`INSERT INTO transactions 
		(amount, category, type, date, note, account_id, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		txn.Amount, txn.Category, txn.Type, txn.Date, txn.Note, txn.AccountID, now, now)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	txn.ID = id

	multiplier := 1.0
	if txn.Type == "expense" {
		multiplier = -1.0
	}
	_, err = s.db.Exec(`UPDATE accounts SET current_balance = current_balance + ?, updated_at = ? WHERE id = ?`,
		txn.Amount*multiplier, now, txn.AccountID)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

func (s *TransactionService) Update(txn *models.Transaction) error {
	var old models.Transaction
	row := s.db.QueryRow(`SELECT amount, type, account_id FROM transactions WHERE id=?`, txn.ID)
	if err := row.Scan(&old.Amount, &old.Type, &old.AccountID); err != nil {
		return err
	}
	oldMult := 1.0
	if old.Type == "expense" {
		oldMult = -1.0
	}
	_, err := s.db.Exec(`UPDATE accounts SET current_balance = current_balance - ?, updated_at = ? WHERE id = ?`,
		old.Amount*oldMult, time.Now(), old.AccountID)
	if err != nil {
		return err
	}

	txn.UpdatedAt = time.Now()
	_, err = s.db.Exec(`UPDATE transactions SET amount=?, category=?, type=?, date=?, note=?, account_id=?, updated_at=? WHERE id=?`,
		txn.Amount, txn.Category, txn.Type, txn.Date, txn.Note, txn.AccountID, txn.UpdatedAt, txn.ID)
	if err != nil {
		return err
	}

	newMult := 1.0
	if txn.Type == "expense" {
		newMult = -1.0
	}
	_, err = s.db.Exec(`UPDATE accounts SET current_balance = current_balance + ?, updated_at = ? WHERE id = ?`,
		txn.Amount*newMult, time.Now(), txn.AccountID)
	return err
}

func (s *TransactionService) Delete(id int64) error {
	var txn models.Transaction
	row := s.db.QueryRow(`SELECT amount, type, account_id FROM transactions WHERE id=?`, id)
	if err := row.Scan(&txn.Amount, &txn.Type, &txn.AccountID); err != nil {
		return err
	}
	mult := 1.0
	if txn.Type == "expense" {
		mult = -1.0
	}
	_, err := s.db.Exec(`UPDATE accounts SET current_balance = current_balance - ?, updated_at = ? WHERE id = ?`,
		txn.Amount*mult, time.Now(), txn.AccountID)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`DELETE FROM transactions WHERE id=?`, id)
	return err
}

func (s *TransactionService) List(filter map[string]interface{}) ([]*models.Transaction, error) {
	query := `SELECT id, amount, category, type, date, note, account_id, created_at, updated_at FROM transactions WHERE 1=1`
	args := []interface{}{}
	if v, ok := filter["start_date"]; ok && v != "" {
		query += " AND date >= ?"
		args = append(args, v)
	}
	if v, ok := filter["end_date"]; ok && v != "" {
		query += " AND date <= ?"
		args = append(args, v)
	}
	if v, ok := filter["category"]; ok && v != "" {
		query += " AND category = ?"
		args = append(args, v)
	}
	if v, ok := filter["type"]; ok && v != "" {
		query += " AND type = ?"
		args = append(args, v)
	}
	if v, ok := filter["search"]; ok && v != "" {
		query += " AND (note LIKE ? OR category LIKE ?)"
		like := "%" + v.(string) + "%"
		args = append(args, like, like)
	}
	query += " ORDER BY date DESC, created_at DESC"

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var txns []*models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.Amount, &t.Category, &t.Type, &t.Date, &t.Note, &t.AccountID, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		txns = append(txns, &t)
	}
	return txns, nil
}