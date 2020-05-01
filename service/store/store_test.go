package store_test

import (
	"context"
	"github.com/barcodepro/go-service-template/service/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStore(t *testing.T) {
	ctx := context.Background()
	st, err := store.NewStore(ctx, store.TestStoreConfig(t))
	assert.NoError(t, err)
	assert.NotNil(t, st)
	assert.NotNil(t, st.Config)
	assert.NotNil(t, st.PgDB)
	//assert.NotNil(t, st.AccountRepository)
}

func TestNewPostgresStore(t *testing.T) {
	var testcases = []struct {
		dsn   string
		valid bool
	}{
		{dsn: "host=postgres dbname=billing_db_test user=postgres sslmode=disable", valid: true},
		{dsn: "host=postgres dbname=billing_db_test user=postgres sslmode=require", valid: true},
		{dsn: "host=postgres dbname=invalid user=postgres sslmode=disable", valid: false},
		{dsn: "invalid_string", valid: false},
	}

	for _, tc := range testcases {
		st, err := store.NewPostgresStore(tc.dsn)
		if tc.valid {
			assert.NoError(t, err)
			assert.NotNil(t, st)
		} else {
			assert.Error(t, err)
			assert.Nil(t, st)
		}
	}
}

func TestStore_Close(t *testing.T) {
	ctx := context.Background()
	sc := store.TestStoreConfig(t)
	st, err := store.NewStore(ctx, sc)
	assert.NoError(t, err)
	assert.NotNil(t, st)

	st.Close()
}
