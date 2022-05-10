/*
Copyright © 2022 Aseem Shrey

*/
package main

import (
	"os"

	"github.com/LuD1161/upi-recon-cli/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	// set log settings
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

var GitCommit string // set using go build ldflags "-X main.GitCommit"

func main() {
	cmd.Execute()
}
