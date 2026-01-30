package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/timurkaev/order-processor/internal/models"
	"time"
)

type Storage interface {
	SaveOrder(ctx context.Context, order *models.Order) error

	CLose()

	Ping(ctx context.Context) error
}

type PostgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(dsn string) (*PostgresStorage, error) {
	ctx := context.Background()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse DSN: %w", err)
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute
	config.HealthCheckPeriod = time.Minute
	config.ConnConfig.ConnectTimeout = 5 * time.Second

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return &PostgresStorage{
		pool: pool,
	}, nil
}

func (s *PostgresStorage) SaveOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (
			order_id, user_id, product_id, quantity,
			price, total_amount, status, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9	
		)		
		ON CONFLICT (order_id)
		DO UPDATE SET
			user_id = EXCLUDED.user_id,
			product_id = EXCLUDED.product_id,
			quantity = EXCLUDED.quantity,
			price = EXLUDED.price,
			total_amount = EXCLUDED.status,
			status = EXCLUDED.status,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

	defer cancel()

	now := time.Now()
	if order.CreatedAt.IsZero() {
		order.CreatedAt = now
	}
	order.UpdatedAt = now

	err := s.pool.QueryRow(
		ctx,
		query,
		order.OrderID,
		order.UserID,
		order.ProductID,
		order.Quantity,
		order.Price,
		order.TotalAmount,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	).Scan(&order.ID)

	if err != nil {
		return fmt.Errorf("failed to save order: %w", err)
	}

	return nil
}

func (s *PostgresStorage) Ping(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	return s.pool.Ping(ctx)
}

func (s *PostgresStorage) Close() {
	s.pool.Close()
}

func (s *PostgresStorage) Stats() *pgxpool.Stat {
	return s.pool.Stat()
}
