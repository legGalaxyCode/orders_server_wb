package payment

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
	"test_db_server/internal/payment"
	"test_db_server/pkg/client/postgresql"
	"test_db_server/pkg/logging"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func (r repository) Create(ctx context.Context, pay *payment.Payment, id string) (bool, error) {
	q := `INSERT INTO payments(transaction, request_id, currency, provider, amount,
                     payment_dt, bank, delivery_cost, goods_total, custom_fee, fk_payments_to_orders)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		  RETURNING transaction`
	var trans string
	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, pay.Transaction, pay.RequestId, pay.Currency, pay.Provider, pay.Amount,
		pay.PaymentDt, pay.Bank, pay.DeliveryCost, pay.GoodsTotal, pay.CustomFee, id).Scan(&trans); err != nil {
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

func (r repository) FindAll(ctx context.Context) ([]payment.Payment, error) {
	q := `SELECT transaction, request_id, currency, provider, amount, payment_dt,
       bank, delivery_cost, goods_total, custom_fee
		  FROM payments;`

	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	items := make([]payment.Payment, 0)
	for rows.Next() {
		var pay payment.Payment
		err = rows.Scan(&pay.Transaction, &pay.RequestId, &pay.Currency, &pay.Provider, &pay.Amount,
			&pay.PaymentDt, &pay.Bank, &pay.DeliveryCost, &pay.GoodsTotal, &pay.CustomFee)
		if err != nil {
			return nil, err
		}
		items = append(items, pay)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	r.logger.Info("exit from findall without error")
	return items, nil
}

func (r repository) FindOne(ctx context.Context, id string) (payment.Payment, error) {
	q := `SELECT transaction, request_id, currency, provider, amount, payment_dt,
       bank, delivery_cost, goods_total, custom_fee
		  FROM payments
		  WHERE fk_payments_to_orders = $1`

	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))

	var pay payment.Payment
	err := r.client.QueryRow(ctx, q, id).Scan(&pay.Transaction, &pay.RequestId, &pay.Currency, &pay.Provider, &pay.Amount,
		&pay.PaymentDt, &pay.Bank, &pay.DeliveryCost, &pay.GoodsTotal, &pay.CustomFee)

	if err != nil {
		return payment.Payment{}, err
	}

	return pay, nil
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) payment.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
