package compiler

import (
	"fishy/pkg/datatype"
	"fishy/pkg/log"
	"fishy/pkg/utils"
)

type Symbol struct {
	name     string
	dataType datatype.DataType
	addr     uint64
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

func (s *SymbolTable) Compile(name string, addr uint64) []byte {
	bytecode := []byte{}
	if symbol := s.Get(name); symbol != nil {
		bytecode = append(bytecode, utils.Bytes8(addr)...)
		bytecode = append(bytecode, byte(symbol.dataType))
	} else {
		log.Warn("tried to compile label that does not exist", "label", name)
	}
	return bytecode
}
