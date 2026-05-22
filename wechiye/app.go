package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"wechiye/backend/config"
	"wechiye/backend/crypto"
	"wechiye/backend/database"
	"wechiye/backend/logger"
	"wechiye/backend/models"
	"wechiye/backend/services"
	"wechiye/backend/sync"
)

type App struct {
	ctx         context.Context
	db          *database.DB
	userService *services.UserService
	accService  *services.AccountService
	txnService  *services.TransactionService
	budService  *services.BudgetService
	coupleSvc   *services.CoupleService
	syncSvc     *sync.SyncService
	backupSvc   *services.BackupService
	config      *config.Config
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logger.Init()

	configDir, err := config.GetConfigDir()
	if err != nil {
		logger.Fatal("Failed to get config directory", "error", err)
	}
	if err := os.MkdirAll(configDir, 0700); err != nil {
		logger.Fatal("Failed to create config directory", "error", err)
	}

	a.config = config.NewConfig(configDir)

	key, err := a.config.LoadMasterKey()
	if err == nil && key != nil {
		dbPath := filepath.Join(configDir, "data.db")
		db, err := database.Open(dbPath, key)
		if err != nil {
			_ = a.config.DeleteMasterKey()
			key = nil
		} else {
			a.db = db
		}
	}

	if a.db == nil {
		runtime.EventsEmit(ctx, "needs-password")
	} else {
		a.initServices()
		runtime.EventsEmit(ctx, "db-ready")
	}
}

func (a *App) initServices() {
	a.userService = services.NewUserService(a.db)
	a.accService = services.NewAccountService(a.db)
	a.txnService = services.NewTransactionService(a.db)
	a.budService = services.NewBudgetService(a.db)
	a.coupleSvc = services.NewCoupleService(a.db)
	a.syncSvc = sync.NewSyncService(a.db, a.coupleSvc)
	a.backupSvc = services.NewBackupService(a.config.DataDir)
}

func (a *App) SetupMasterPassword(password string, remember bool) error {
	key, err := database.DeriveKey(password)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	dbPath := filepath.Join(a.config.DataDir, "data.db")
	_, err = os.Stat(dbPath)
	dbExists := !os.IsNotExist(err)

	db, err := database.Open(dbPath, key)
	if err != nil {
		return fmt.Errorf("invalid password or database corrupt")
	}

	a.db = db
	if !dbExists {
		if err := database.Migrate(db); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	if remember {
		if err := a.config.SaveMasterKey(key); err != nil {
			logger.Warn("Failed to save key to keychain", "error", err)
		}
	} else {
		_ = a.config.DeleteMasterKey()
	}

	a.initServices()
	runtime.EventsEmit(a.ctx, "db-ready")
	return nil
}

func (a *App) Shutdown() {
	if a.db != nil {
		a.db.Close()
	}
}

// User Profile
func (a *App) GetUser() (*models.User, error) {
	return a.userService.Get()
}

func (a *App) UpdateUser(user *models.User) error {
	return a.userService.Update(user)
}

func (a *App) UpdateAvatar(avatar string) error {
	return a.userService.UpdateAvatar(avatar)
}

// Accounts
func (a *App) AddAccount(name string, initialBalance float64) (*models.Account, error) {
	return a.accService.Create(name, initialBalance)
}

func (a *App) UpdateAccount(id int64, name string, initialBalance float64) error {
	return a.accService.Update(id, name, initialBalance)
}

func (a *App) DeleteAccount(id int64) error {
	return a.accService.Delete(id)
}

func (a *App) GetAccounts() ([]*models.Account, error) {
	return a.accService.List()
}

// Transactions
func (a *App) AddTransaction(txn *models.Transaction) (*models.Transaction, error) {
	return a.txnService.Create(txn)
}

func (a *App) UpdateTransaction(txn *models.Transaction) error {
	return a.txnService.Update(txn)
}

func (a *App) DeleteTransaction(id int64) error {
	return a.txnService.Delete(id)
}

func (a *App) GetTransactions(filter map[string]interface{}) ([]*models.Transaction, error) {
	return a.txnService.List(filter)
}

// Budgets
func (a *App) SetBudget(budget *models.Budget) (*models.Budget, error) {
	return a.budService.Set(budget)
}

func (a *App) GetBudgets(monthYear string) ([]*models.Budget, error) {
	return a.budService.List(monthYear)
}

func (a *App) DeleteBudget(id int64) error {
	return a.budService.Delete(id)
}

// Couple Linking
func (a *App) GeneratePairingQR() (string, error) {
	pubKey, err := a.coupleSvc.GetLocalPublicKey()
	if err != nil {
		return "", err
	}
	nonce := a.coupleSvc.GenerateNonce()
	data := fmt.Sprintf("wechiye:link:%s:%s:1.0", pubKey, nonce)
	qrPNG, err := crypto.GenerateQRCode(data)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(qrPNG), nil
}

func (a *App) ProcessScannedQR(qrData string) error {
	return a.coupleSvc.ProcessPairingQR(qrData)
}

func (a *App) GetCoupleStatus() (map[string]interface{}, error) {
	return a.coupleSvc.GetStatus()
}

// Sync
func (a *App) ExportEncryptedSharedData() (string, error) {
	data, err := a.syncSvc.ExportSharedData()
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) ImportEncryptedSharedData(b64Data string) error {
	data, err := base64.StdEncoding.DecodeString(b64Data)
	if err != nil {
		return err
	}
	return a.syncSvc.ImportSharedData(data)
}

// Backup
func (a *App) ExportBackup() (string, error) {
	return a.backupSvc.Export()
}

func (a *App) RestoreBackup(backupPath string) error {
	return a.backupSvc.Restore(backupPath)
}

// Utility
func (a *App) GetAppVersion() string {
	return "1.0.0"
}