package delivery

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
	"test_db_server/internal/delivery"
	"test_db_server/pkg/client/postgresql"
	"test_db_server/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func (r *repository) Create(ctx context.Context, del *delivery.Delivery, orderId string) (bool, error) {
	q := `INSERT INTO deliveries(name, phone, zip, city, address, region, email, fk_deliveries_to_orders)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		  returning delivery_uid`
	var id int
	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, del.Name, del.Phone, del.Zip, del.City, del.Address,
		del.Region, del.Email, orderId).Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detaii: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			r.logger.Error(newErr)
			return false, newErr
		}
		return false, err
	}
	return true, nil
}

func (r *repository) FindAll(ctx context.Context) ([]delivery.Delivery, error) {
	q := `SELECT name, phone, zip, city, address, region, email
		  FROM deliveries;`

	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	items := make([]delivery.Delivery, 0)
	for rows.Next() {
		var del delivery.Delivery
		err = rows.Scan(&del.Name, &del.Phone, &del.Zip, &del.City, &del.Address, &del.Region, &del.Email)
		if err != nil {
			return nil, err
		}
		items = append(items, del)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	r.logger.Info("exit from findall without error")
	return items, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (delivery.Delivery, error) {
	q := `SELECT name, phone, zip, city, address, region, email
		  FROM deliveries
		  WHERE fk_deliveries_to_orders = $1`

	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))

	var del delivery.Delivery
	err := r.client.QueryRow(ctx, q, id).Scan(&del.Name, &del.Phone, &del.Zip, &del.City, &del.Address, &del.Region, &del.Email)

	if err != nil {
		return delivery.Delivery{}, err
	}

	return del, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) delivery.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
