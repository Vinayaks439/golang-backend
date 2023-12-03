package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTransfer(t *testing.T) {
	args := CreateTransferParams{
		FromAccountID: 1,
		ToAccountID:   1,
		Amount:        100,
	}
	transfer, err := testQueries.CreateTransfer(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)
	require.Equal(t, args.FromAccountID, transfer.FromAccountID)
	require.Equal(t, args.ToAccountID, transfer.ToAccountID)
	require.Equal(t, args.Amount, transfer.Amount)
	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)
}

func TestGetTransfer(t *testing.T) {
	getTransfer, err := testQueries.GetTransfer(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, getTransfer)
	require.Equal(t, getTransfer.ID, int64(1))
	require.Equal(t, getTransfer.FromAccountID, int64(1))
	require.Equal(t, getTransfer.ToAccountID, int64(1))
	require.Equal(t, getTransfer.Amount, int64(100))
	require.NotZero(t, getTransfer.ID)
	require.NotZero(t, getTransfer.CreatedAt)
}

func TestListTransfers(t *testing.T) {
	args := ListTransfersParams{
		FromAccountID: 1,
		ToAccountID:   1,
		Limit:         5,
		Offset:        0,
	}
	listTransfers, err := testQueries.ListTransfers(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, listTransfers)
	for _, transfer := range listTransfers {
		require.NotEmpty(t, transfer)
		require.Equal(t, int64(1), transfer.FromAccountID)
		require.Equal(t, int64(1), transfer.ToAccountID)
	}
}
