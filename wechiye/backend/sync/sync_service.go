package sync

import (
	"encoding/json"
	"time"

	"wechiye/backend/crypto"
	"wechiye/backend/database"
	"wechiye/backend/models"
	"wechiye/backend/services"
)

type SyncService struct {
	db        *database.DB
	coupleSvc *services.CoupleService
}

func NewSyncService(db *database.DB, coupleSvc *services.CoupleService) *SyncService {
	return &SyncService{db: db, coupleSvc: coupleSvc}
}

type SharedData struct {
	Transactions []*models.Transaction `json:"transactions"`
	Budgets      []*models.Budget      `json:"budgets"`
	Timestamp    time.Time             `json:"timestamp"`
}

func (s *SyncService) ExportSharedData() ([]byte, error) {
	txns, err := s.getTransactions()
	if err != nil {
		return nil, err
	}
	budgets, err := s.getBudgets()
	if err != nil {
		return nil, err
	}
	shared := SharedData{
		Transactions: txns,
		Budgets:      budgets,
		Timestamp:    time.Now().UTC(),
	}
	plain, err := json.Marshal(shared)
	if err != nil {
		return nil, err
	}
	secret, err := s.coupleSvc.GetSharedSecret()
	if err != nil {
		return nil, err
	}
	encrypted, err := crypto.EncryptAESGCM(plain, secret)
	if err != nil {
		return nil, err
	}
	return encrypted, nil
}

func (s *SyncService) ImportSharedData(encrypted []byte) error {
	secret, err := s.coupleSvc.GetSharedSecret()
	if err != nil {
		return err
	}
	plain, err := crypto.DecryptAESGCM(encrypted, secret)
	if err != nil {
		return err
	}
	var shared SharedData
	if err := json.Unmarshal(plain, &shared); err != nil {
		return err
	}
	for _, txn := range shared.Transactions {
		s.mergeTransaction(txn)
	}
	for _, bud := range shared.Budgets {
		s.mergeBudget(bud)
	}
	return nil
}

func (s *SyncService) mergeTransaction(txn *models.Transaction) {
	var existing models.Transaction
	row := s.db.QueryRow(`SELECT id, updated_at FROM transactions WHERE id = ? OR (date = ? AND amount = ? AND category = ? AND type = ? AND account_id = ?)`,
		txn.ID, txn.Date, txn.Amount, txn.Category, txn.Type, txn.AccountID)
	err := row.Scan(&existing.ID, &existing.UpdatedAt)
	if err == nil {
		if txn.UpdatedAt.After(existing.UpdatedAt) {
			_, _ = s.db.Exec(`UPDATE transactions SET amount=?, category=?, type=?, date=?, note=?, account_id=?, updated_at=? WHERE id=?`,
				txn.Amount, txn.Category, txn.Type, txn.Date, txn.Note, txn.AccountID, txn.UpdatedAt, existing.ID)
		}
	} else {
		_, _ = s.db.Exec(`INSERT INTO transactions (id, amount, category, type, date, note, account_id, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			txn.ID, txn.Amount, txn.Category, txn.Type, txn.Date, txn.Note, txn.AccountID, txn.CreatedAt, txn.UpdatedAt)
	}
}

func (s *SyncService) mergeBudget(b *models.Budget) {
	var existing models.Budget
	row := s.db.QueryRow(`SELECT id, updated_at FROM budgets WHERE category = ? AND month_year = ?`, b.Category, b.MonthYear)
	err := row.Scan(&existing.ID, &existing.UpdatedAt)
	if err == nil {
		if b.UpdatedAt.After(existing.UpdatedAt) {
			_, _ = s.db.Exec(`UPDATE budgets SET amount_limit=?, updated_at=? WHERE id=?`, b.AmountLimit, b.UpdatedAt, existing.ID)
		}
	} else {
		_, _ = s.db.Exec(`INSERT INTO budgets (category, month_year, amount_limit, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
			b.Category, b.MonthYear, b.AmountLimit, b.CreatedAt, b.UpdatedAt)
	}
}

func (s *SyncService) getTransactions() ([]*models.Transaction, error) {
	rows, err := s.db.Query(`SELECT id, amount, category, type, date, note, account_id, created_at, updated_at FROM transactions`)
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

func (s *SyncService) getBudgets() ([]*models.Budget, error) {
	rows, err := s.db.Query(`SELECT id, category, month_year, amount_limit, created_at, updated_at FROM budgets`)
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