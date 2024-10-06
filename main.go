package main

import (
	"fishy/cmd"
	"fishy/internal/lexer"
	"fishy/pkg/ast"
	"fishy/pkg/token"
	"fmt"
	"strings"
)

func main() {
	cmd.Execute()

	// data, err := os.ReadFile("test.fi")
	// if err != nil {
	// 	panic(err)
	// }

	// pp := preprocessor.New(data)
	// newData := pp.Process()

	// // fmt.Println(newData)

	// p := parser.New(lexer.New(newData))
	// stmts, err := p.Parse()
	// if err != nil {
	// 	panic(err)
	// }

	// printAST(stmts)
	// println()

	// c := compiler.New(stmts)
	// bytecode, err := c.Compile()
	// if err != nil {
	// 	panic(err)
	// }

	// printBytecode(bytecode)
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
		return lexer.Registers[v.Value]
	case *ast.Identifier:
		return v.Value
	case *ast.AddressOf:
		return fmt.Sprintf("[%s]", stringValue(v.Value))
	case *ast.BinaryExpression:
		op := "???"
		switch v.Operator {
		case token.PLUS:
			op = "+"
		case token.MINUS:
			op = "-"
		case token.STAR:
			op = "*"
		case token.SLASH:
			op = "/"
		}
		return fmt.Sprintf("(%s %s %s)", stringValue(v.Left), op, stringValue(v.Right))
	}

	return fmt.Sprintf("%#v", value)
}

func printBytecode(bytecode []byte) {
	numLines := (len(bytecode) + 15) / 16

	for line := 0; line < numLines; line++ {
		startIndex := line * 16
		endIndex := startIndex + 16
		if endIndex > len(bytecode) {
			endIndex = len(bytecode)
		}
		lineBytes := bytecode[startIndex:endIndex]

		fmt.Printf("0x%04X: ", startIndex)

		for _, b := range lineBytes {
			fmt.Printf("%02X ", b)
		}
		fmt.Println()
	}
}
