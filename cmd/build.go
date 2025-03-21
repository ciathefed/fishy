package cmd

import (
	"fishy/internal/compiler"
	"fishy/internal/lexer"
	"fishy/internal/parser"
	"fishy/internal/preprocessor"
	"fishy/pkg/ast"
	"fishy/pkg/log"
	"fishy/pkg/token"
	"fishy/pkg/utils"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build [file]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Compile a FishyASM file to Fishy Bytecode",
	Run: func(cmd *cobra.Command, args []string) {
		inputFile := args[0]

		var source string
		if !skipPreprocessing {
			pp, err := preprocessor.New(inputFile, false)
			if err != nil {
				log.Fatal(err)
			}
			err = pp.Process()
			if err != nil {
				log.Fatal(err)
			}
			source = pp.Output()
		} else {
			sourceBytes, err := os.ReadFile(inputFile)
			if err != nil {
				log.Fatal(err)
			}
			source = string(sourceBytes)

			if verbose {
				log.Infof("skipping pre-processing")
			}
		}

		l := lexer.New(source)
		if debugLexer {
			log.Info("debugging lexer")
			for {
				t := l.NextToken()
				fmt.Printf("Token(kind=%#v, value=%#v, pos=(%d, %d))\n", t.Kind, t.Value, t.Start, t.End)
				if t.Kind == token.EOF {
					break
				}
			}
			os.Exit(0)
		}

		p := parser.New(l)
		statements, err := p.Parse()
		if err != nil {
			log.Fatal(err)
		}
		if debugParser {
			log.Info("debugging parser")
			printAST(statements)
			os.Exit(0)
		}

		c := compiler.New(statements)
		bytecode, err := c.Compile()
		if err != nil {
			log.Fatal(err)
		}

		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		n, err := file.Write(bytecode)
		if err != nil {
			log.Fatal(err)
		}

		if verbose {
			log.Info("wrote bytecode to output", "file", outputFile, "bytes", n)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	buildCmd.Flags().StringVarP(&outputFile, "output", "o", "out.fbc", "output file")
	buildCmd.Flags().BoolVarP(&skipPreprocessing, "skip-pre-processing", "", false, "skip pre-processing stage")
	buildCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	buildCmd.Flags().BoolVarP(&debugLexer, "debug-lexer", "", false, "only dump the lexer output")
	buildCmd.Flags().BoolVarP(&debugParser, "debug-parser", "", false, "only dump the parser output")
}

func printAST(stmts []ast.Statement) {
	for _, stmt := range stmts {
		switch s := stmt.(type) {
		case *ast.Label:
			fmt.Printf("%s(name=%#v)\n", s.String(), s.Name)
		case *ast.Instruction:
			fmt.Printf("%s(name=%#v, args=[", s.String(), s.Name)
			sv := []string{}
			for _, v := range s.Args {
				sv = append(sv, stringValue(v))
			}
			fmt.Print(strings.Join(sv, ", "))
			fmt.Println("])")
		case *ast.Sequence:
			fmt.Printf("%s(name=%#v, values=[", s.String(), s.Name)
			sv := []string{}
			for _, v := range s.Values {
				sv = append(sv, stringValue(v))
			}
			fmt.Print(strings.Join(sv, ", "))
			fmt.Println("])")
		}
	}
}

func stringValue(value ast.Value) string {
	switch v := value.(type) {
	case *ast.NumberLiteral:
		return v.Value
	case *ast.StringLiteral:
		return fmt.Sprintf("%#v", v.Value)
	case *ast.Register:
		return utils.Registers[v.Value]
	case *ast.Identifier:
		return v.Value
	case *ast.AddressOf:
		return fmt.Sprintf("[%s]", stringValue(v.Value))
	case *ast.RegisterOffsetNumber:
		return fmt.Sprintf("%s %s %s", stringValue(&v.Left), v.Operator.String(), stringValue(&v.Right))
	case *ast.RegisterOffsetRegister:
		return fmt.Sprintf("%s %s %s", stringValue(&v.Left), v.Operator.String(), stringValue(&v.Right))
	case *ast.LabelOffsetNumber:
		return fmt.Sprintf("%s %s %s", stringValue(v.Left), v.Operator.String(), stringValue(&v.Right))
	case *ast.LabelOffsetRegister:
		return fmt.Sprintf("%s %s %s", stringValue(v.Left), v.Operator.String(), stringValue(&v.Right))
	}

	return fmt.Sprintf("%#v", value)
}
