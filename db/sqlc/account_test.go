package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/nanorex07/gobank/util"
	"github.com/stretchr/testify/assert"
)

func CreateRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:    util.RandomStr(7),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	acc, err := testQueries.CreateAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.NotEmpty(t, acc)

	assert.Equal(t, arg.Owner, acc.Owner)
	assert.Equal(t, arg.Balance, acc.Balance)
	assert.Equal(t, arg.Currency, acc.Currency)
	assert.NotZero(t, acc.ID)
	assert.NotZero(t, acc.CreatedAt)
	return acc
}

func TestCreateAccount(t *testing.T) {
	CreateRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account_1 := CreateRandomAccount(t)
	account_2, err := testQueries.GetAccount(context.Background(), account_1.ID)
	assert.NoError(t, err)
	assert.Equal(t, account_1.ID, account_2.ID)
	assert.Equal(t, account_1.Owner, account_2.Owner)
	assert.Equal(t, account_1.Balance, account_2.Balance)
	assert.Equal(t, account_1.Currency, account_2.Currency)
	assert.WithinDuration(t, account_1.CreatedAt, account_2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account_1 := CreateRandomAccount(t)

	arg := UpdateAccountParams{
		ID:      account_1.ID,
		Balance: util.RandomMoney(),
	}
	var err error
	var account_2 Account
	account_2, err = testQueries.UpdateAccount(context.Background(), arg)
	assert.NoError(t, err)

	assert.NoError(t, err)
	assert.NotEmpty(t, account_2)
	assert.Equal(t, arg.Balance, account_2.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account_1 := CreateRandomAccount(t)

	var err error
	var account_2 Account
	err = testQueries.DeleteAccount(context.Background(), account_1.ID)
	assert.NoError(t, err)

	account_2, err = testQueries.GetAccount(context.Background(), account_1.ID)
	assert.Error(t, err)
	assert.EqualError(t, err, sql.ErrNoRows.Error())
	assert.Empty(t, account_2)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAccount(t)
	}
	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}
	accounts, err := testQueries.ListAccount(context.Background(), arg)
	assert.NoError(t, err)
	assert.Len(t, accounts, 5)

	for _, account := range accounts {
		assert.NotEmpty(t, account)
	}
}
