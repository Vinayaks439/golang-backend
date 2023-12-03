package db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateEntries(t *testing.T) {
	args := CreateEntriesParams{
		Amount:    100,
		AccountID: int64(1),
	}
	entry, err := testQueries.CreateEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.Amount, entry.Amount)
	require.Equal(t, args.AccountID, entry.AccountID)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)
}

func TestGetEntry(t *testing.T) {
	getEntry, err := testQueries.GetEntry(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, getEntry)
	require.Equal(t, getEntry.ID, int64(1))
	require.Equal(t, getEntry.AccountID, int64(1))
	require.Equal(t, getEntry.Amount, int64(100))
	require.NotZero(t, getEntry.ID)
	require.NotZero(t, getEntry.CreatedAt)
}

func TestListEntries(t *testing.T) {
	entries, err := testQueries.ListEntries(context.Background(), ListEntriesParams{
		AccountID: int64(1),
		Limit:     5,
		Offset:    0,
	})
	require.NoError(t, err)
	require.NotEmpty(t, entries)
	for _, entry := range entries {
		require.NotEmpty(t, entry)
		require.Equal(t, entry.AccountID, int64(1))
		require.Equal(t, entry.Amount, int64(100))
		require.NotZero(t, entry.ID)
		require.NotZero(t, entry.CreatedAt)
	}
}

func TestUpdateEntry(t *testing.T) {
	args := UpdateEntryParams{
		ID:     1,
		Amount: 200,
	}
	entry, err := testQueries.UpdateEntry(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)
	require.Equal(t, args.ID, entry.ID)
	require.Equal(t, args.Amount, entry.Amount)
	require.NotZero(t, entry.ID)
	require.NotZero(t, entry.CreatedAt)

}
func TestDeleteEntry(t *testing.T) {
	entry, err := testQueries.CreateEntries(context.Background(), CreateEntriesParams{
		Amount:    100,
		AccountID: int64(2),
	})
	getEntry, err := testQueries.GetEntry(context.Background(), entry.ID)
	err = testQueries.DeleteEntry(context.Background(), 1)
	require.NoError(t, err)
	require.NotEmpty(t, getEntry)
	require.Equal(t, entry.ID, getEntry.ID)
	require.Equal(t, entry.AccountID, getEntry.AccountID)
	require.Equal(t, entry.Amount, getEntry.Amount)
}
