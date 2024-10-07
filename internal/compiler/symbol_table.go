package compiler

import (
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/utils"
)

type Symbol struct {
	name     string
	dataType datatype.DataType
	addr     int
	section  Section
}

type SymbolTable struct {
	symbols map[string]*Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		symbols: make(map[string]*Symbol),
	}
}

func (s *SymbolTable) Set(name string, symbol *Symbol) {
	s.symbols[name] = symbol
}

func (s *SymbolTable) Get(name string) *Symbol {
	return s.symbols[name]
}

func (s *SymbolTable) Compile(name string, addr int) []byte {
	bytecode := []byte{}
	if symbol := s.Get(name); symbol != nil {
		bytecode = append(bytecode, utils.Bytes4(uint32(addr))...)
		bytecode = append(bytecode, utils.Bytes4(uint32(symbol.dataType))...)
	} else {
		log.Warn("tried to compile label that does not exist", "label", name)
	}
	return bytecode
}
