package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Database function list (derived from docker/db/init/002_seed.sql).
//
// The seed flow is:
// vendors -> features -> actions -> action_revisions -> commands
//
// Connection / lifecycle:
// - Open(dsn string) (*Store, error): centralizes pool creation and startup validation.
// - Close() error: cleanly releases pool resources on shutdown.
// - Ping() error: verifies DB reachability before handling requests.
// - WithTx(fn func(*TxStore) error) error: runs grouped writes atomically.
//
// Vendor functions:
// - CreateVendor(name string, deviceType string) (int, error): inserts a vendor row.
// - GetVendorByName(name string) (Vendor, error): resolves vendor identity for downstream queries.
// - ListVendors() ([]Vendor, error): returns all vendors for selection/display.
//
// Feature functions:
// - CreateFeature(vendorID int, featureName string) (int, error): adds a feature under a vendor.
// - GetFeature(vendorID int, featureName string) (Feature, error): finds a specific vendor feature.
// - ListFeaturesByVendor(vendorID int) ([]Feature, error): lists feature options per vendor.
//
// Action functions:
// - CreateAction(featureID int, action string) (int, error): creates a CRUD action for a feature.
// - GetAction(featureID int, action string) (Action, error): looks up one action record.
// - ListActionsByFeature(featureID int) ([]Action, error): enumerates allowed actions for a feature.
//
// Revision functions:
// - CreateActionRevision(actionID int, createdBy string, comment string) (int, error): records version history for action commands.
// - SetCurrentRevision(actionID int, revisionID int) error: switches which revision is active.
// - GetCurrentRevisionID(actionID int) (int, error): resolves active revision for reads/writes.
//
// Command functions:
// - AddCommand(revisionID int, position int, command string) error: inserts one ordered command line.
// - ReplaceCommandsForRevision(revisionID int, commands []string) error: rewrites full command set for a revision.
// - ListCommandsByRevision(revisionID int) ([]Command, error): fetches ordered commands for execution/output.
//
// Read path used by command generation:
// - GetCommands(vendorName string, featureName string, action string) ([]string, error): returns final command templates for CLI use.
// - BuildConfigSnapshot() (ConfigSnapshot, error): materializes DB data into app-friendly config shape.
//
// Bootstrap / seed helpers:
// - EnsureInitialRevision(actionID int, createdBy string) (int, error): guarantees each action has a starting revision.
// - SeedDefaults() error: loads baseline vendors/features/actions/commands.

type VendorTable struct {
	Id         int
	Name       string
	DeviceType string
}

func OpenDatabaseConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	var err error

	dbpool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return dbpool, nil
}

func getConnectionFromPool(ctx context.Context, dbpool *pgxpool.Pool) (*pgxpool.Conn, error) {
	var err error

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	connection, err := dbpool.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	return connection, err
}

func startTransaction(ctx context.Context, connection *pgxpool.Conn) (pgx.Tx, error) {
	var err error

	transaction, err := connection.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func processRows[T any](rows pgx.Rows, scan func(pgx.Rows) (T, error)) ([]T, error) {
	defer rows.Close()

	out := make([]T, 0)
	for rows.Next() {
		item, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()

}

func SelectAllOnTable[T any](ctx context.Context, dbpool *pgxpool.Pool, table string, tableProcessor func(pgx.Rows) (T, error)) ([]T, error) {
	var err error
	var output []T

	query := fmt.Sprintln("SELECT * FROM ", table)

	connection, err := getConnectionFromPool(ctx, dbpool)
	if err != nil {
		return output, err
	}
	defer connection.Release()

	rows, err := connection.Query(ctx, query)
	if err != nil {
		return output, err
	}
	output, err = processRows(rows, tableProcessor)
	if err != nil {
		return output, err
	}

	return output, nil
}

func ProcessVendorTable(rows pgx.Rows) (VendorTable, error) {
	var output VendorTable

	err := rows.Scan(output.Id, output.Name, output.DeviceType)
	if err != nil {
		return output, err
	}
	return output, nil

}

func UpdateVendorName(ctx context.Context, dbpool *pgxpool.Pool, vendorName string, vendorId string) (string, error) {
	var err error
	var output string

	connection, err := getConnectionFromPool(ctx, dbpool)
	if err != nil {
		return output, err
	}
	defer connection.Release()

	tx, err := startTransaction(ctx, connection)
	if err != nil {
		return output, err
	}
	defer tx.Rollback(ctx)

	rows, err := tx.Exec(
		ctx,
		`UPDATE name = $1
		FROM vendors
		WHERE name = $2`,
		vendorName,
		vendorId,
	)
	if err != nil {
		return output, err
	}
	tx.Commit(ctx)

	return rows.String(), nil
}

func SelectVendorbyId(ctx context.Context, dbpool *pgxpool.Pool, vendorId string) ([]VendorTable, error) {
	var err error
	var output []VendorTable

	connection, err := getConnectionFromPool(ctx, dbpool)
	if err != nil {
		return output, err
	}
	defer connection.Release()

	rows, err := connection.Query(
		ctx,
		`Select *
		FROM vendors
		WHERE id = $1`,
		vendorId,
	)
	if err != nil {
		return output, err
	}
	output, err = processRows(rows, ProcessVendorTable)
	if err != nil {
		return output, err
	}

	return output, err
}
