package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	args := CreateAccountParams{
		Owner:    "Test",
		Balance:  100,
		Currency: "USD",
	}
	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.Owner, account.Owner)
	require.Equal(t, args.Balance, account.Balance)
	require.Equal(t, args.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestGetAccount(t *testing.T) {
	getAccount, err := testQueries.GetAccount(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, getAccount)
	require.Equal(t, getAccount.ID, int64(1))
	require.Equal(t, getAccount.Owner, "Test")
	require.Equal(t, getAccount.Balance, int64(100))
	require.Equal(t, getAccount.Currency, "USD")
	require.NotZero(t, getAccount.ID)
	require.NotZero(t, getAccount.CreatedAt)
}

func TestListAccounts(t *testing.T) {
	arg := ListAccountsParams{
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		require.Equal(t, "Test", account.Owner)
	}
}

func TestUpdateAccount(t *testing.T) {
	args := UpdateAccountParams{
		ID:      1,
		Balance: 200,
	}
	account, err := testQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, args.ID, account.ID)
	require.Equal(t, args.Balance, account.Balance)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}

func TestDeleteAccount(t *testing.T) {
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Owner:    "Test_DELETE",
		Balance:  100,
		Currency: "USD",
	})
	getAccount, err := testQueries.GetAccount(context.Background(), account.ID)
	err = testQueries.DeleteAccount(context.Background(), getAccount.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.NotEmpty(t, getAccount)
	require.Equal(t, getAccount.ID, account.ID)
	require.Equal(t, getAccount.Owner, account.Owner)
	require.Equal(t, getAccount.Balance, account.Balance)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
