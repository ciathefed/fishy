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

func (s *SymbolTable) GetSize() int {
	longestAddr := uint64(0)
	for _, sym := range s.symbols {
		if sym.addr > longestAddr {
			longestAddr = sym.addr
		}
	}

	switch {
	case longestAddr <= 0xFFFF:
		return 2
	case longestAddr <= 0xFFFFFFFF:
		return 4
	default:
		return 8
	}
}

func (s *SymbolTable) Compile(name string, addr uint64) []byte {
	bytecode := []byte{}
	size := s.GetSize()

	if symbol := s.Get(name); symbol != nil {
		switch size {
		case 2:
			bytecode = append(bytecode, utils.Bytes2(uint16(addr))...)
			bytecode = append(bytecode, byte(symbol.dataType))
		case 4:
			bytecode = append(bytecode, utils.Bytes4(uint32(addr))...)
			bytecode = append(bytecode, byte(symbol.dataType))
		default:
			bytecode = append(bytecode, utils.Bytes8(addr)...)
			bytecode = append(bytecode, byte(symbol.dataType))
		}
	} else {
		log.Warn("tried to compile label that does not exist", "label", name)
	}
	return bytecode
}
