CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    order_id VARCHAR(255) NOT NULL UNIQUE,
    user_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),

    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0),
    total_amount DECIMAL(12, 2) NOT NULL CHECK (total_amount >= 0),
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);

CREATE INDEX IF NOT EXISTS idx_orders_product_id ON orders(product_id);

CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

CREATE INDEX IF NOT EXISTS idx_orders_created_status ON orders(created_at DESC, status);

COMMENT ON TABLE orders IS 'Таблица заказов из Kafka топика';
COMMENT ON COLUMN orders.order_id IS 'Внешний ID заказа (из Kafka)';
COMMENT ON COLUMN orders.total_amount IS 'Общая сумма: quantity * price';