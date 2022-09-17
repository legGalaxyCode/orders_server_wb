CREATE TABLE IF NOT EXISTS orders(
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TEXT,
    oof_shard TEXT
);

CREATE TABLE IF NOT EXISTS deliveries(
    delivery_uid INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT,
    fk_deliveries_to_orders TEXT REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS payments(
    transaction TEXT PRIMARY KEY,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INT,
    payment_dt INT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    fk_payments_to_orders TEXT REFERENCES orders(order_uid)
);

CREATE TABLE IF NOT EXISTS items(
    chrt_id INT PRIMARY KEY,
    track_number TEXT,
    price INT,
    rid TEXT,
    name TEXT,
    sale INT,
    size TEXT,
    total_price INT,
    nm_id INT,
    brand TEXT,
    status INT,
    fk_items_to_orders TEXT REFERENCES orders(order_uid)
);

CREATE INDEX IF NOT EXISTS index_order_uid ON orders(order_uid);

-- INSERT INTO orders(order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
-- VALUES ('b563feb7b2b84b6test', 'WBILMTESTTRACK', 'WBIL', )