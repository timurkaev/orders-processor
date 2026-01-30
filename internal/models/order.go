package models

import "time"

type Order struct {
	ID          int64   `json:"-" db:"id"`
	OrderID     string  `json:"order_id" db:"order_id"`
	UserID      int64   `json:"user_id" db:"user_id"`
	ProductID   int64   `json:"product_id" db:"product_id"`
	Quantity    int     `json:"quantity" db:"quantity"`
	Price       float64 `json:"price" db:"price"`
	TotalAmount float64 `json:"total_amount" db:"total_amount"`
	Status      string  `json:"status" db:"status"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(msg string) error {
	return &ValidationError{Message: msg}
}

func (o *Order) Validate() error {
	if o.OrderID == "" {
		return NewValidationError("order_id is required")
	}

	if o.UserID <= 0 {
		return NewValidationError("user_id must be positive")
	}

	if o.ProductID <= 0 {
		return NewValidationError("product_id must be positive")
	}

	if o.Quantity <= 0 {
		return NewValidationError("quantity must be positive")
	}

	if o.Price < 0 {
		return NewValidationError("price cannot be negative")
	}

	if o.TotalAmount < 0 {
		return NewValidationError("total_amount cannot be negative")
	}

	return nil
}

func (o *Order) CalculateTotal() {
	o.TotalAmount = float64(o.Quantity) * o.Price
}
