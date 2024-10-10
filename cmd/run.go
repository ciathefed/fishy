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

		m := vm.New(inputData, memorySize, false)
		m.Run()

		if vomitRegisters > -1 {
			msg := "main thread"
			if vomitRegisters > 0 {
				msg = fmt.Sprintf("thread %d", vomitRegisters)
			}

			if _, ok := m.GetThread(vomitRegisters); !ok {
				log.Errorf("%s does not exist", msg)
			} else {
				log.Infof("vomiting %s registers 🤮", msg)
				m.DumpRegisters(vomitRegisters)
			}
		} else if vomitRegisters > -2 {
			log.Info("vomiting registers 🤮")
			i := 0
			for {
				_, ok := m.GetThread(i)
				if !ok {
					break
				}

				if i == 0 {
					log.Info("main thread:")
				} else {
					log.Infof("thread %d", i)
				}

				m.DumpRegisters(i)

				i++
			}
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
	runCmd.Flags().IntVarP(&vomitRegisters, "vomit-registers", "", -2, "dump the registers at the index when done (-1 = all)")
	runCmd.Flags().BoolVarP(&vomitMemory, "vomit-memory", "", false, "dump the memory when done")
}
