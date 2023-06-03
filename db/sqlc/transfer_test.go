package db

import (
	"context"
	"github.com/astronely/simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomTransfer(t *testing.T, accountFrom, accountTo Account) Transfer {
	arg := CreateTransferParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)

	require.NotZero(t, transfer.ID)
	require.NotZero(t, transfer.CreatedAt)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	createRandomTransfer(t, accountFrom, accountTo)
}

func TestGetTransfer(t *testing.T) {
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, accountFrom, accountTo)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer2.CreatedAt, time.Second)
}

func TestListTransfers(t *testing.T) {
	var n = 5
	accountFrom := createRandomAccount(t)
	accountTo := createRandomAccount(t)
	for i := 0; i < n*2; i++ {
		createRandomTransfer(t, accountFrom, accountTo)
	}

	arg := ListTransfersParams{
		FromAccountID: accountFrom.ID,
		ToAccountID:   accountTo.ID,
		Limit:         int32(n),
		Offset:        int32(n),
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, n)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
		require.True(t, transfer.FromAccountID == accountFrom.ID && transfer.ToAccountID == accountTo.ID)
	}
}
