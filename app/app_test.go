package app_test

import (
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"cosmossdk.io/log"
	"github.com/stretchr/testify/require"

	"github.com/Lord1Egypt/Maat/app"
)

func TestAppInitialization(t *testing.T) {
	db := dbm.NewMemDB()
	logger := log.NewNopLogger()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("application initialization panicked: %v", r)
		}
	}()

	maatApp := app.New(logger, db, nil, true, nil)
	require.NotNil(t, maatApp)
}
