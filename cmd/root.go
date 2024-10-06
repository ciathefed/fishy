package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFile        string
	skipPreprocessing bool
	vomitLexer        bool
	vomitParser       bool
	vomitRegisters    bool
	vomitMemory       bool
	verbose           bool
	memorySize        int
)

var rootCmd = &cobra.Command{
	Use:   "fishy",
	Short: "Fishy is a CLI tool used for the Fishy echo system",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
