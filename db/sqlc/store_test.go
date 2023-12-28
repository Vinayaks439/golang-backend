package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
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
	log.Println("Before transfer")
	log.Println("account1.Balance :", account1.Balance)
	log.Println("account2.Balance :", account2.Balance)
	// Transfer 50 from account1 to account2
	amount := int64(50)
	n := int64(2)
	results := make(chan TransferTxResult)
	errs := make(chan error)
	for i := int64(0); i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	// check errors
	existed := make(map[int]bool)
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
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		log.Println("During transfer :", i)
		log.Println("account1.Balance :", fromAccount.Balance)
		log.Println("account2.Balance :", toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)
		k := diff1 / amount
		require.True(t, k <= n && k >= 1)
		require.NotContains(t, existed, k)
		existed[int(k)] = true
	}
	// check account update balance

	updateAccount1, err := store.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	updateAccount2, err := store.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	log.Println("After transfer")
	log.Println("account1.Balance :", updateAccount1.Balance)
	log.Println("account2.Balance :", updateAccount2.Balance)
	require.Equal(t, updateAccount1.Balance, account1.Balance-n*amount)

	require.Equal(t, updateAccount2.Balance, account2.Balance+n*amount)
}
