package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	outputFile        string
	skipPreprocessing bool
	debugLexer        bool
	debugParser       bool
	debugRegisters    int
	debugMemory       bool
	verbose           bool
	memorySize        int
)

var rootCmd = &cobra.Command{
	Use:   "fishy",
	Short: "Fishy is a CLI tool used for the Fishy ecosystem",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
