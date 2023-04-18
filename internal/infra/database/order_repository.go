package database

import (
	"context"
	"database/sql"

	"github.com/sesaquecruz/goexpert-clean-architecture-lab/internal/entity"
)

type OrderRepository struct {
	Db *sql.DB
}

func NewOrderRepository(db *sql.DB) *OrderRepository {
	return &OrderRepository{
		Db: db,
	}
}

func (r *OrderRepository) Save(ctx context.Context, order entity.Order) error {
	stmt, err := r.Db.Prepare("INSERT INTO orders (id, price, tax, final_price) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	stmt.Close()

	_, err = stmt.ExecContext(ctx, order.Id, order.Price, order.Tax, order.FinalPrice)
	return err
}

func (r *OrderRepository) FindAll(ctx context.Context) ([]entity.Order, error) {
	rows, err := r.Db.QueryContext(ctx, "SELECT id, price, tax, final_price FROM orders")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []entity.Order
	for rows.Next() {
		var order entity.Order
		if err := rows.Scan(&order.Id, &order.Price, &order.Tax, &order.FinalPrice); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
