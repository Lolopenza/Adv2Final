package postgres

import (
	"context"
	"database/sql"
	"log"

	"payment-service/internal/domain"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) domain.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, customer_email, plan_name, price, currency, status, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query,
		subscription.ID,
		subscription.CustomerEmail,
		subscription.PlanName,
		subscription.Price,
		subscription.Currency,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)
	return err
}

func (r *subscriptionRepository) GetByID(id string) (*domain.Subscription, error) {
	query := `
		SELECT id, customer_email, plan_name, price, currency, status, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`
	subscription := &domain.Subscription{}
	err := r.db.QueryRow(query, id).Scan(
		&subscription.ID,
		&subscription.CustomerEmail,
		&subscription.PlanName,
		&subscription.Price,
		&subscription.Currency,
		&subscription.Status,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return subscription, nil
}

func (r *subscriptionRepository) Update(subscription *domain.Subscription) error {
	query := `
		UPDATE subscriptions
		SET customer_email = $1, plan_name = $2, price = $3, currency = $4, status = $5, 
			start_date = $6, end_date = $7, updated_at = $8
		WHERE id = $9
	`
	_, err := r.db.Exec(query,
		subscription.CustomerEmail,
		subscription.PlanName,
		subscription.Price,
		subscription.Currency,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.UpdatedAt,
		subscription.ID,
	)
	return err
}

func (r *subscriptionRepository) Delete(id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *subscriptionRepository) List(customerEmail string, page, limit int32) ([]*domain.Subscription, int32, error) {
	log.Printf("Repository List method called with: customerEmail=%q, page=%d, limit=%d", customerEmail, page, limit)

	// Calculate offset
	offset := (page - 1) * limit

	// Build the query based on whether customerEmail is provided
	var countQuery, listQuery string
	var args []interface{}
	var countArgs []interface{}

	// Base queries
	countBase := "SELECT COUNT(*) FROM subscriptions"
	listBase := `
		SELECT id, customer_email, plan_name, price, currency, status, start_date, end_date, created_at, updated_at
		FROM subscriptions
	`

	// Check if we need to filter by customer email
	if customerEmail != "" {
		countQuery = countBase + " WHERE customer_email = $1"
		listQuery = listBase + " WHERE customer_email = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3"
		args = []interface{}{customerEmail, limit, offset}
		countArgs = []interface{}{customerEmail}
	} else {
		countQuery = countBase
		listQuery = listBase + " ORDER BY created_at DESC LIMIT $1 OFFSET $2"
		args = []interface{}{limit, offset}
		countArgs = []interface{}{} // Empty slice for no args
	}

	// First get the total count
	var total int32
	var err error
	if customerEmail != "" {
		err = r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	} else {
		err = r.db.QueryRow(countQuery).Scan(&total)
	}

	if err != nil {
		return nil, 0, err
	}

	// Now get the actual subscriptions
	rows, err := r.db.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var subscriptions []*domain.Subscription
	for rows.Next() {
		subscription := &domain.Subscription{}
		err := rows.Scan(
			&subscription.ID,
			&subscription.CustomerEmail,
			&subscription.PlanName,
			&subscription.Price,
			&subscription.Currency,
			&subscription.Status,
			&subscription.StartDate,
			&subscription.EndDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

// Implements a transaction that creates both a subscription and a payment
func (r *subscriptionRepository) CreateWithPayment(ctx context.Context, subscription *domain.Subscription, payment *domain.Payment) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	// Insert subscription
	subQuery := `
		INSERT INTO subscriptions (id, customer_email, plan_name, price, currency, status, start_date, end_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err = tx.ExecContext(ctx, subQuery,
		subscription.ID,
		subscription.CustomerEmail,
		subscription.PlanName,
		subscription.Price,
		subscription.Currency,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.CreatedAt,
		subscription.UpdatedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// Insert payment
	payQuery := `
		INSERT INTO payments (id, amount, currency, status, customer_email, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, payQuery,
		payment.ID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.CustomerEmail,
		payment.Description,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
