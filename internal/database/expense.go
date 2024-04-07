package database

import (
	"context"
	"time"

	"github.com/Svoevolin/workshop_1_bot/internal/domain"
)

type ExpenseDB struct {
	store map[int64][]domain.Expense
}

func NewExpenseDB() (*ExpenseDB, error) {
	return &ExpenseDB{
		store: make(map[int64][]domain.Expense),
	}, nil
}

func (db *ExpenseDB) AddExpense(ctx context.Context, userID int64, kopecks int64, title string, date time.Time) error {
	db.store[userID] = append(db.store[userID], domain.Expense{
		Title:  title,
		Date:   date,
		Amount: kopecks,
	})
	return nil
}

func (db *ExpenseDB) GetExpenses(ctx context.Context, userID int64) ([]domain.Expense, error) {
	return db.store[userID], nil
}