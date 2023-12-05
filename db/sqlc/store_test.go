package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	account_1 := CreateRandomAccount(t)
	account_2 := CreateRandomAccount(t)

	n := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account_1.ID,
				ToAccountID:   account_2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	var existed map[int]bool = map[int]bool{}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
		result := <-results
		assert.NotEmpty(t, result)

		transfer := result.Transfer
		assert.NotEmpty(t, transfer)
		assert.Equal(t, account_1.ID, transfer.FromAccountID)
		assert.Equal(t, account_2.ID, transfer.ToAccountID)
		assert.Equal(t, transfer.Amount, amount)
		assert.NotZero(t, transfer.ID)
		assert.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		// check entries
		fromEntry := result.FromEntry
		assert.NotEmpty(t, fromEntry)
		assert.Equal(t, account_1.ID, fromEntry.AccountID)
		assert.Equal(t, -amount, fromEntry.Amount)
		assert.NotZero(t, fromEntry.ID)
		assert.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		assert.NoError(t, err)

		toEntry := result.ToEntry
		assert.NotEmpty(t, toEntry)
		assert.Equal(t, account_2.ID, toEntry.AccountID)
		assert.Equal(t, amount, toEntry.Amount)
		assert.NotZero(t, toEntry.ID)
		assert.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		assert.NoError(t, err)

		fromAccount := result.FromAccount
		assert.NotEmpty(t, fromAccount)
		assert.Equal(t, fromAccount.ID, account_1.ID)

		toAccount := result.ToAccount
		assert.NotEmpty(t, toAccount)
		assert.Equal(t, toAccount.ID, account_2.ID)

		fmt.Println(">> Tx: ", fromAccount.Balance, toAccount.Balance)

		diff1 := account_1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account_2.Balance
		assert.Equal(t, diff1, diff2)
		assert.True(t, diff1 > 0)
		assert.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		assert.True(t, k >= 1 && k <= n)
		assert.NotContains(t, existed, k)
		existed[k] = true

	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account_1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account_2.ID)
	assert.NoError(t, err)

	fmt.Println(">> After: ", updateAccount1.Balance, updateAccount2.Balance)

	assert.Equal(t, account_1.Balance-int64(n)*amount, updateAccount1.Balance)
	assert.Equal(t, account_2.Balance+int64(n)*amount, updateAccount2.Balance)

}

func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account_1 := CreateRandomAccount(t)
	account_2 := CreateRandomAccount(t)

	fmt.Println(">> Before", account_1.Balance, account_2.Balance)

	n := 10
	amount := int64(10)

	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := account_1.ID
		toAccountID := account_2.ID

		if i%2 == 1 {
			fromAccountID = account_2.ID
			toAccountID = account_1.ID
		}

		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		assert.NoError(t, err)
	}

	updateAccount1, err := testQueries.GetAccount(context.Background(), account_1.ID)
	assert.NoError(t, err)

	updateAccount2, err := testQueries.GetAccount(context.Background(), account_2.ID)
	assert.NoError(t, err)

	fmt.Println(">> After: ", updateAccount1.Balance, updateAccount2.Balance)

	assert.Equal(t, account_1.Balance, updateAccount1.Balance)
	assert.Equal(t, account_2.Balance, updateAccount2.Balance)

}
