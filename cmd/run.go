package cmd

import (
	"fishy/internal/vm"
	"fishy/pkg/log"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [file]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Run Fishy Bytecode file in the FishyVM",
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]
		inputData, err := os.ReadFile(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			log.Info("read bytecode from input", "file", inputFile, "bytes", len(inputData))
		}

		m := vm.New(inputData, memorySize)
		m.Run()

		m.DumpRegisters()
		fmt.Println()
		m.DumpMemory(0, 128)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&memorySize, "memory-size", "s", 1024, "total amount of memory to use")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
}
