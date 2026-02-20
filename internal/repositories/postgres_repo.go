package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/debt-optimization-engine/internal/models"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	CreateGroup(ctx context.Context, group *models.Group) error
	AddMemberToGroup(ctx context.Context, groupID, userID string) error
	CreateExpense(ctx context.Context, expense *models.Expense) error
	GetGroupMembers(ctx context.Context, groupID string) ([]models.User, error)
	GetExpensesByGroup(ctx context.Context, groupID string, from, to *time.Time) ([]models.Expense, error)
}

type PostgresRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(pool *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{pool: pool}
}

func (r *PostgresRepo) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email) VALUES ($1, $2) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, user.Username, user.Email).Scan(&user.ID, &user.CreatedAt)
}

func (r *PostgresRepo) CreateGroup(ctx context.Context, group *models.Group) error {
	query := `INSERT INTO groups (name) VALUES ($1) RETURNING id, created_at`
	return r.pool.QueryRow(ctx, query, group.Name).Scan(&group.ID, &group.CreatedAt)
}

func (r *PostgresRepo) AddMemberToGroup(ctx context.Context, groupID, userID string) error {
	query := `INSERT INTO group_members (group_id, user_id) VALUES ($1, $2)`
	_, err := r.pool.Exec(ctx, query, groupID, userID)
	return err
}

func (r *PostgresRepo) CreateExpense(ctx context.Context, expense *models.Expense) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO expenses (group_id, payer_id, amount, description, split_type) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err = tx.QueryRow(ctx, query, expense.GroupID, expense.PayerID, expense.Amount, expense.Description, expense.SplitType).
		Scan(&expense.ID, &expense.CreatedAt)
	if err != nil {
		return err
	}

	for _, split := range expense.Splits {
		splitQuery := `INSERT INTO expense_splits (expense_id, user_id, amount) VALUES ($1, $2, $3)`
		_, err = tx.Exec(ctx, splitQuery, expense.ID, split.UserID, split.Amount)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgresRepo) GetGroupMembers(ctx context.Context, groupID string) ([]models.User, error) {
	query := `SELECT u.id, u.username, u.email, u.created_at FROM users u
	          JOIN group_members gm ON u.id = gm.user_id WHERE gm.group_id = $1`
	rows, err := r.pool.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *PostgresRepo) GetExpensesByGroup(ctx context.Context, groupID string, from, to *time.Time) ([]models.Expense, error) {
	query := `SELECT id, payer_id, amount, description, split_type, created_at FROM expenses WHERE group_id = $1`
	args := []interface{}{groupID}

	if from != nil {
		query += ` AND created_at >= $2`
		args = append(args, *from)
	}
	if to != nil {
		if from != nil {
			query += ` AND created_at <= $3`
		} else {
			query += ` AND created_at <= $2`
		}
		args = append(args, *to)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.PayerID, &e.Amount, &e.Description, &e.SplitType, &e.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Get splits for this expense
		splitQuery := `SELECT user_id, amount FROM expense_splits WHERE expense_id = $1`
		sRows, err := r.pool.Query(ctx, splitQuery, e.ID)
		if err != nil {
			return nil, err
		}
		var splits []models.ExpenseSplit
		for sRows.Next() {
			var s models.ExpenseSplit
			if err := sRows.Scan(&s.UserID, &s.Amount); err != nil {
				sRows.Close()
				return nil, err
			}
			splits = append(splits, s)
		}
		sRows.Close()
		e.Splits = splits
		expenses = append(expenses, e)
	}
	return expenses, nil
}
