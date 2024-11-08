package database

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Amele9/call-manager/internal/config"
	"github.com/Amele9/call-manager/internal/models"
)

type Database interface {
	// CreateCall creates a call
	CreateCall(call *models.CallInfo) (int, error)

	// GetCalls returns all calls
	GetCalls() ([]*models.CallInfo, error)

	// GetCallInfo returns call information
	GetCallInfo(ID int) (*models.CallInfo, error)

	// UpdateCallStatus updates a call status
	UpdateCallStatus(ID int) error

	// DeleteCall deletes a call
	DeleteCall(ID int) error

	// Shutdown closes the database connection
	Shutdown(ctx context.Context) error
}

var notFoundError = errors.New("record not found")

// New returns an instance of the PostgreSQLDatabase
func New(config *config.Configuration) (*PostgreSQLDatabase, error) {
	// We need to wait for the database to start
	time.Sleep(time.Second)

	connection, err := pgxpool.New(context.Background(), config.ConnectionString)
	if err != nil {
		return nil, err
	}

	err = connection.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return &PostgreSQLDatabase{connection: connection}, nil
}

// PostgreSQLDatabase is a structure for working with PostgreSQL database
type PostgreSQLDatabase struct {
	connection *pgxpool.Pool
}

// CreateCall creates a call in the PostgreSQL database
func (d *PostgreSQLDatabase) CreateCall(call *models.CallInfo) (int, error) {
	query := `
		INSERT INTO calls (client_name, phone_number, description)
		VALUES ($1, $2, $3)
		RETURNING id, status, created_at
	`

	err := d.connection.
		QueryRow(context.Background(), query, call.ClientName, call.PhoneNumber, call.Description).
		Scan(&call.ID, &call.Status, &call.CreatedAt)
	if err != nil {
		return 0, err
	}

	return call.ID, nil
}

// GetCalls returns all calls from the PostgreSQL database
func (d *PostgreSQLDatabase) GetCalls() ([]*models.CallInfo, error) {
	var count int
	err := d.connection.QueryRow(context.Background(), "SELECT COUNT(*) FROM calls").Scan(&count)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id,
		       client_name,
		       phone_number,
		       description,
		       status,
		       created_at
		FROM calls
	`

	rows, err := d.connection.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	calls := make([]*models.CallInfo, 0, count)

	for rows.Next() {
		var call models.CallInfo

		err = rows.Scan(
			&call.ID,
			&call.ClientName,
			&call.PhoneNumber,
			&call.Description,
			&call.Status,
			&call.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		calls = append(calls, &call)
	}

	return calls, nil
}

// GetCallInfo returns call information from the PostgreSQL database
func (d *PostgreSQLDatabase) GetCallInfo(ID int) (*models.CallInfo, error) {
	query := `
		SELECT id,
		       client_name,
		       phone_number,
		       description,
		       status,
		       created_at
		FROM calls
		WHERE id = $1
	`

	var call models.CallInfo

	err := d.connection.QueryRow(context.Background(), query, ID).Scan(
		&call.ID,
		&call.ClientName,
		&call.PhoneNumber,
		&call.Description,
		&call.Status,
		&call.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, notFoundError
		}

		return nil, err
	}

	return &call, nil
}

// UpdateCallStatus updates a call status in the PostgreSQL database
func (d *PostgreSQLDatabase) UpdateCallStatus(ID int) error {
	query := `UPDATE calls SET status = 'closed' WHERE id = $1`

	result, err := d.connection.Exec(context.Background(), query, ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return notFoundError
	}

	return nil
}

// DeleteCall deletes a call from the PostgreSQL database
func (d *PostgreSQLDatabase) DeleteCall(ID int) error {
	query := `DELETE FROM calls WHERE id = $1`

	result, err := d.connection.Exec(context.Background(), query, ID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return notFoundError
	}

	return nil
}

// Shutdown closes the PostgreSQL database connection
func (d *PostgreSQLDatabase) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		d.connection.Close()

		done <- struct{}{}

		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
