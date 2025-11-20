package orders

import (
	"context"
	"database/sql"
)

type MySQLRepository struct {
	db *sql.DB

	stmtGetAll  *sql.Stmt
	stmtInsert  *sql.Stmt
	stmtGetByID *sql.Stmt
}

func NewMySQLRepository(db *sql.DB) (*MySQLRepository, error) {
	r := &MySQLRepository{db: db}

	var err error

	r.stmtGetAll, err = db.Prepare(`
		SELECT id, amount, status, created_at, updated_at
		FROM orders
		ORDER BY id DESC
		LIMIT 100
	`)
	if err != nil {
		return nil, err
	}

	r.stmtInsert, err = db.Prepare(`
		INSERT INTO orders (amount, status)
		VALUES (?, ?)
	`)
	if err != nil {
		return nil, err
	}

	r.stmtGetByID, err = db.Prepare(`
		SELECT id, amount, status, created_at, updated_at
		FROM orders
		WHERE id = ?
	`)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *MySQLRepository) GetAll(ctx context.Context) ([]Order, error) {
	rows, err := r.stmtGetAll.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Order

	for rows.Next() {
		var o Order
		if err := rows.Scan(&o.ID, &o.Amount, &o.Status, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, o)
	}

	return result, rows.Err()
}

func (r *MySQLRepository) Insert(ctx context.Context, amount int, status string) (*Order, error) {
	res, err := r.stmtInsert.ExecContext(ctx, amount, status)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	var o Order
	err = r.stmtGetByID.QueryRowContext(ctx, id).
		Scan(&o.ID, &o.Amount, &o.Status, &o.CreatedAt, &o.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *MySQLRepository) Close() {
	if r.stmtGetAll != nil {
		r.stmtGetAll.Close()
	}
	if r.stmtInsert != nil {
		r.stmtInsert.Close()
	}
	if r.stmtGetByID != nil {
		r.stmtGetByID.Close()
	}
}
