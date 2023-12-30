package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
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

var txKey = struct{}{}

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
		txName := ctx.Value(txKey)
		log.Println("create transfer , txName:", txName)
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		log.Println("create entry 1 , txName:", txName)
		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		log.Println("create entry 2 , txName:", txName)
		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		if arg.FromAccountID < arg.ToAccountID {
			log.Println("update account 1&2 balance, txName:", txName)
			result.FromAccount, result.ToAccount, err = AddMoney(ctx, q, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			log.Println("update account 2&1 balance, txName:", txName)
			result.ToAccount, result.FromAccount, err = AddMoney(ctx, q, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount)
			if err != nil {
				return err
			}
		}
		return nil

	})
	// Return the result and error
	return result, err
}

func AddMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	accountID2 int64,
	amount1 int64,
	amount2 int64,
) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
