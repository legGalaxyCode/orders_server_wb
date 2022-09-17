package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"strings"
	delivery "test_db_server/internal/delivery/db"
	item22 "test_db_server/internal/item"
	item2 "test_db_server/internal/item/db"
	"test_db_server/internal/order"
	payment "test_db_server/internal/payment/db"
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

func (r *repository) Create(ctx context.Context, order *order.Order) (bool, error) {
	q := `INSERT INTO orders(
                   order_uid, track_number, entry, locale,
                   internal_signature, customer_id, delivery_service, shardkey,
                   sm_id, date_created, oof_shard
                   ) 
	      VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	      RETURNING order_uid`
	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))
	var chartId string
	if err := r.client.QueryRow(ctx, q, order.OrderUid, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSig, order.CustomerId, order.DeliveryService, order.ShardKey,
		order.SmId, order.DateCreated, order.OofShard).Scan(&chartId); err != nil && err != sql.ErrNoRows {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detaii: %s, Where: %s", pgErr.Message, pgErr.Detail, pgErr.Where))
			r.logger.Error(newErr)
			return false, newErr
		}
		return false, err
	}

	// delivery
	deliveryRep := delivery.NewRepository(r.client, r.logger)
	_, err := deliveryRep.Create(ctx, &order.Delivery, order.OrderUid)
	if err != nil {
		return false, err
	}
	// payment
	paymentRep := payment.NewRepository(r.client, r.logger)
	_, err = paymentRep.Create(ctx, &order.Payment, order.OrderUid)
	if err != nil {
		return false, err
	}
	// items
	itemRep := item2.NewRepository(r.client, r.logger)
	for _, it := range order.Items {
		r.logger.Infof("current item: %v", it)
		_, err = itemRep.Create(ctx, &it, order.OrderUid)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (r *repository) CreateMany(ctx context.Context, orders ...*order.Order) (bool, error) {
	for _, ord := range orders {
		create, err := r.Create(ctx, ord)
		if err != nil {
			return create, err
		}
	}
	return true, nil
}

func (r *repository) FindAll(ctx context.Context) ([]order.Order, error) {
	q := `SELECT orders.order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	  deliveries.name, phone, zip, city, address, region, email,
	  transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee,
	  chrt_id, items.track_number, price, rid, items.name, sale, size, total_price, nm_id, brand, status
		  FROM orders, deliveries, payments, items
		  WHERE fk_deliveries_to_orders = orders.order_uid AND fk_payments_to_orders = orders.order_uid AND fk_items_to_orders = orders.order_uid`

	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	orders := make([]order.Order, 0)
	for rows.Next() {
		var ord order.Order
		ord.Items = make([]item22.Item, 1)
		err = rows.Scan(&ord.OrderUid, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSig,
			&ord.CustomerId, &ord.DeliveryService, &ord.ShardKey, &ord.SmId, &ord.DateCreated, &ord.OofShard,
			&ord.Delivery.Name, &ord.Delivery.Phone, &ord.Delivery.Zip, &ord.Delivery.City, &ord.Delivery.Address,
			&ord.Delivery.Region, &ord.Delivery.Email, &ord.Payment.Transaction, &ord.Payment.RequestId, &ord.Payment.Currency,
			&ord.Payment.Provider, &ord.Payment.Amount, &ord.Payment.PaymentDt, &ord.Payment.Bank, &ord.Payment.DeliveryCost,
			&ord.Payment.GoodsTotal, &ord.Payment.CustomFee, &ord.Items[0].ChartId, &ord.Items[0].TrackNumber, &ord.Items[0].Price,
			&ord.Items[0].RId, &ord.Items[0].Name, &ord.Items[0].Sale, &ord.Items[0].Size, &ord.Items[0].TotalPrice,
			&ord.Items[0].NmId, &ord.Items[0].Brand, &ord.Items[0].Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	r.logger.Info("exit from FindAll() without error")

	/*
		for now we can scan orders slice to find duplicates
		and then unite it in one by items slice
	*/

	r.logger.Infof("orders len: %d", len(orders))
	uniqueOrders := make(map[string]order.Order, 0)
	for _, ord := range orders {
		if value, f := uniqueOrders[ord.OrderUid]; f {
			r.logger.Infof("find order with uid: %s; items len: %d and cap: %d", ord.OrderUid, len(value.Items), cap(value.Items))
			value.Items = append(value.Items, ord.Items[0]) // should be only one
			uniqueOrders[ord.OrderUid] = value
			r.logger.Infof("after items len: %d and cap: %d", len(value.Items), cap(value.Items))
		} else {
			uniqueOrders[ord.OrderUid] = ord
		}
	}
	result := make([]order.Order, 0)
	for _, ord := range uniqueOrders {
		r.logger.Infof("order with uid %s has %d items and %d capacity", ord.OrderUid, len(ord.Items), cap(ord.Items))
		result = append(result, ord)
	}

	return result, nil
}

func (r *repository) FindOne(ctx context.Context, orderId string) (order.Order, error) {
	q := `SELECT orders.order_uid, orders.track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard,
	  deliveries.name, phone, zip, city, address, region, email,
	  transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee,
	  chrt_id, items.track_number, price, rid, items.name, sale, size, total_price, nm_id, brand, status
		  FROM orders, deliveries, payments, items
		  WHERE orders.order_uid = $1 AND 
		        fk_deliveries_to_orders = orders.order_uid AND
		        fk_payments_to_orders = orders.order_uid AND
		        fk_items_to_orders = orders.order_uid`

	r.logger.Info(fmt.Sprintf("SQL query %s", formatQuery(q)))

	// here is logic:
	/*
		we get all rows with identical orderId but different items
		then we unite them in one order
		it's not efficiently but works
	*/

	rows, err := r.client.Query(ctx, q, orderId)
	if err != nil {
		return order.Order{}, err
	}
	orders := make([]order.Order, 0)
	for rows.Next() {
		var ord order.Order
		ord.Items = make([]item22.Item, 1)
		err = rows.Scan(&ord.OrderUid, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSig,
			&ord.CustomerId, &ord.DeliveryService, &ord.ShardKey, &ord.SmId, &ord.DateCreated, &ord.OofShard,
			&ord.Delivery.Name, &ord.Delivery.Phone, &ord.Delivery.Zip, &ord.Delivery.City, &ord.Delivery.Address,
			&ord.Delivery.Region, &ord.Delivery.Email, &ord.Payment.Transaction, &ord.Payment.RequestId, &ord.Payment.Currency,
			&ord.Payment.Provider, &ord.Payment.Amount, &ord.Payment.PaymentDt, &ord.Payment.Bank, &ord.Payment.DeliveryCost,
			&ord.Payment.GoodsTotal, &ord.Payment.CustomFee, &ord.Items[0].ChartId, &ord.Items[0].TrackNumber, &ord.Items[0].Price,
			&ord.Items[0].RId, &ord.Items[0].Name, &ord.Items[0].Sale, &ord.Items[0].Size, &ord.Items[0].TotalPrice,
			&ord.Items[0].NmId, &ord.Items[0].Brand, &ord.Items[0].Status)
		if err != nil {
			return order.Order{}, err
		}
		orders = append(orders, ord)
	}
	if err = rows.Err(); err != nil {
		return order.Order{}, err
	}

	finalOrder := order.Order{
		OrderUid:        orders[0].OrderUid,
		TrackNumber:     orders[0].TrackNumber,
		Entry:           orders[0].Entry,
		Delivery:        orders[0].Delivery,
		Payment:         orders[0].Payment,
		Items:           nil,
		Locale:          orders[0].Locale,
		InternalSig:     orders[0].InternalSig,
		CustomerId:      orders[0].CustomerId,
		DeliveryService: orders[0].DeliveryService,
		ShardKey:        orders[0].ShardKey,
		SmId:            orders[0].SmId,
		DateCreated:     orders[0].DateCreated,
		OofShard:        orders[0].OofShard,
	}
	finalOrder.Items = make([]item22.Item, 0)
	for _, it := range orders {
		finalOrder.Items = append(finalOrder.Items, it.Items[0]) // because there's guarantee only one item
	}
	return finalOrder, nil
}

func NewRepository(client postgresql.Client, logger *logging.Logger) order.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
