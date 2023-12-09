package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(conn)
	account1, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "Test_TransferTx_Account1",
		Balance:  1000,
		Currency: "USD",
	})
	require.NoError(t, err)
	require.NotEmpty(t, account1)
	require.Equal(t, account1.Balance, int64(1000))
	require.Equal(t, account1.Owner, "Test_TransferTx_Account1")
	account2, err := store.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "Test_TransferTx_Account2",
		Balance:  1000,
		Currency: "USD",
	})
	require.NoError(t, err)
	require.NotEmpty(t, account2)
	require.Equal(t, account2.Balance, int64(1000))
	require.Equal(t, account2.Owner, "Test_TransferTx_Account2")

	// Transfer 50 from account1 to account2
	amount := int64(50)
	n := int64(5)
	results := make(chan TransferTxResult)
	errs := make(chan error)
	for i := int64(0); i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	// check errors
	for i := int64(0); i < n; i++ {
		err := <-errs
		result := <-results
		require.NoError(t, err)
		require.NotEmpty(t, result)
		require.NotEmpty(t, result.Transfer)
		require.Equal(t, result.Transfer.FromAccountID, account1.ID)
		require.Equal(t, result.Transfer.ToAccountID, account2.ID)
		require.Equal(t, result.Transfer.Amount, amount)
		require.NotZero(t, result.Transfer.ID)

		_, err = store.GetTransfer(context.Background(), result.Transfer.ID)
		require.NoError(t, err)

		//check account entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.ID)
		require.Equal(t, fromEntry.Amount, -amount)
		require.NotZero(t, fromEntry.ID)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//to entry check
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.ID)
		require.Equal(t, toEntry.Amount, amount)
		require.NotZero(t, toEntry.ID)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//TODO: account balance
	}
}
