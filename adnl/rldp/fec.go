package rldp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/xssnick/tonutils-go/tl"
)

func init() {
	tl.Register(FECRaptorQ{}, "fec.raptorQ data_size:int symbol_size:int symbols_count:int = fec.Type")
	tl.Register(FECRoundRobin{}, "fec.roundRobin data_size:int symbol_size:int symbols_count:int = fec.Type")
	tl.Register(FECOnline{}, "fec.online data_size:int symbol_size:int symbols_count:int = fec.Type")
}

type FECRaptorQ struct {
	DataSize     uint32 // `tl:"int"`
	SymbolSize   uint32 // `tl:"int"`
	SymbolsCount uint32 // `tl:"int"`
}

func (f *FECRaptorQ) Parse(data []byte) ([]byte, error) {
	if len(data) < 12 {
		return nil, fmt.Errorf("fec raptor data too short")
	}
	f.DataSize = binary.LittleEndian.Uint32(data[:4])
	f.SymbolSize = binary.LittleEndian.Uint32(data[4:8])
	f.SymbolsCount = binary.LittleEndian.Uint32(data[8:12])
	return data[12:], nil
}

func (f *FECRaptorQ) Serialize(buf *bytes.Buffer) error {
	tmp := make([]byte, 12)
	binary.LittleEndian.PutUint32(tmp[0:4], f.DataSize)
	binary.LittleEndian.PutUint32(tmp[4:8], f.SymbolSize)
	binary.LittleEndian.PutUint32(tmp[8:12], f.SymbolsCount)
	buf.Write(tmp)
	return nil
}

type FECRoundRobin struct {
	DataSize     uint32 `tl:"int"`
	SymbolSize   uint32 `tl:"int"`
	SymbolsCount uint32 `tl:"int"`
}

type FECOnline struct {
	DataSize     uint32 `tl:"int"`
	SymbolSize   uint32 `tl:"int"`
	SymbolsCount uint32 `tl:"int"`
}
