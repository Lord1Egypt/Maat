package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/Lord1Egypt/Maat/app"
	"github.com/Lord1Egypt/Maat/cmd/maatd/cmd"
)

func main() {
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount("maat", "maatpub")
	config.SetBech32PrefixForValidator("maatvaloper", "maatvaloperpub")
	config.SetBech32PrefixForConsensusNode("maatvalcons", "maatvalconspub")
	config.Seal()

	rootCmd := cmd.NewRootCmd()

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
