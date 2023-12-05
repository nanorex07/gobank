package db

import (
	"context"
	"testing"

	"github.com/nanorex07/gobank/util"
	"github.com/stretchr/testify/assert"
)

func CreateRandomEntry(t *testing.T) Entry {
	account := CreateRandomAccount(t)
	arg := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomInt(0, account.Balance),
	}
	entry, err := testQueries.CreateEntry(context.Background(), arg)
	assert.NoError(t, err)
	assert.Equal(t, entry.AccountID, arg.AccountID)
	assert.Equal(t, entry.Amount, arg.Amount)

	return entry
}

func TestCreateEntry(t *testing.T) {
	CreateRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry_1 := CreateRandomEntry(t)

	entry_2, err := testQueries.GetEntry(context.Background(), entry_1.ID)
	assert.NoError(t, err)
	assert.Equal(t, entry_1.ID, entry_2.ID)
	assert.Equal(t, entry_1.AccountID, entry_2.AccountID)
	assert.Equal(t, entry_1.Amount, entry_2.Amount)
	assert.Equal(t, entry_1.CreatedAt, entry_2.CreatedAt)
}

func TestListEntries(t *testing.T) {
	account := CreateRandomAccount(t)
	for i := 0; i < 10; i++ {
		arg := CreateEntryParams{
			AccountID: account.ID,
			Amount:    util.RandomInt(0, account.Balance),
		}
		entry, err := testQueries.CreateEntry(context.Background(), arg)
		assert.NoError(t, err)
		assert.Equal(t, entry.AccountID, arg.AccountID)
		assert.Equal(t, entry.Amount, arg.Amount)
	}
	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}
	entires, err := testQueries.ListEntries(context.Background(), arg)
	assert.NoError(t, err)
	assert.Equal(t, len(entires), 5)
	for _, entry := range entires {
		assert.NotEmpty(t, entry)
	}
}
