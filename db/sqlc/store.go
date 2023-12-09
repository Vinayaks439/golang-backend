package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *Store {
	return &Store{
		connPool: connPool,
		Queries:  New(connPool),
	}
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer Transfer `json:"transfer"`
	// SenderAccount and ReceiverAccount are the updated account information
	FromAccount Account `json:"from_account"`
	ToAccount   Account `json:"to_account"`
	FromEntry   Entry   `json:"from_entry"`
	ToEntry     Entry   `json:"to_entry"`
}

// execTx executes a function within a database transaction.
func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	// initiate a transaction
	tx, err := store.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	// Create an instance of the Queries struct with the transaction
	q := New(tx)
	// Call the function that executes the queries
	err = fn(q)
	// if error, rollback the transaction
	if err != nil {
		// rollback the transaction, record the error if any.
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			// if rollback error, return the rollback error
			return rollbackErr
		}
		// else return the error
		return err
	}
	// if no error, commit the transaction
	return tx.Commit(ctx)
}

// Transfertx performs money transfer from one account to another
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	// Create a variable to store the result
	var result TransferTxResult
	// Execute the transaction
	err := store.execTx(ctx, func(q *Queries) error {
		// Create a new entry for the sender
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//TODO update accounts balance
		return nil

	})
	// Return the result and error
	return result, err
}
