package item

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
	"test_db_server/internal/item"
	"test_db_server/pkg/client/postgresql"
	"test_db_server/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, item *item.Item, orderId string) (bool, error) {
	q := `INSERT INTO items(chrt_id, track_number, price, rid, name,
                  sale, size, total_price, nm_id, brand, status, fk_items_to_orders) 
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		  RETURNING chrt_id`
	r.logger.Info(fmt.Sprintf("SQL query %s", strings.ReplaceAll(q, "\t", "")))
	var tempId int
	if err := r.client.QueryRow(ctx, q, item.ChartId, item.TrackNumber, item.Price, item.RId, item.Name, item.Sale, item.Size,
		item.TotalPrice, item.NmId, item.Brand, item.Status, orderId).Scan(&tempId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			r.logger.Error(fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detaii: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where)))
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *repository) FindAll(ctx context.Context) ([]item.Item, error) {
	q := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
		  FROM items;`

	r.logger.Info(fmt.Sprintf("SQL query %s", strings.ReplaceAll(q, "\t", "")))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	items := make([]item.Item, 0)
	for rows.Next() {
		var it item.Item
		err = rows.Scan(&it.ChartId, &it.TrackNumber, &it.Price, &it.RId, &it.Name, &it.Sale, &it.Size,
			&it.TotalPrice, &it.NmId, &it.Brand, &it.Status)
		if err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	r.logger.Info("exit from findall without error")
	return items, nil
}

func (r *repository) FindOne(ctx context.Context, orderId string) (item.Item, error) {
	q := `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, fk_items_to_orders
		  FROM items
		  WHERE fk_items_to_orders = $1`

	r.logger.Info(fmt.Sprintf("SQL query %s", strings.ReplaceAll(q, "\t", "")))

	var it item.Item
	err := r.client.QueryRow(ctx, q, orderId).Scan(&it.ChartId, &it.TrackNumber, &it.Price, &it.RId, &it.Name, &it.Sale, &it.Size,
		&it.TotalPrice, &it.NmId, &it.Brand, &it.Status)

	if err != nil {
		return item.Item{}, err
	}

	return it, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) item.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
