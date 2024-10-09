package cmd

import (
	"fishy/internal/vm"
	"fishy/pkg/log"
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

		m := vm.New(inputData, memorySize, false)
		m.Run()

		if vomitRegisters {
			log.Info("vomiting registers 🤮")
			m.DumpRegisters()
		}

		if vomitMemory {
			log.Info("vomiting memory 🤮")
			m.DumpMemory(0, memorySize)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().IntVarP(&memorySize, "memory-size", "s", 1024, "total amount of memory to use")
	runCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	runCmd.Flags().BoolVarP(&vomitRegisters, "vomit-registers", "", false, "dump the registers when done")
	runCmd.Flags().BoolVarP(&vomitMemory, "vomit-memory", "", false, "dump the memory when done")
}
