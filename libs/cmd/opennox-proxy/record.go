package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
)

type Record struct {
	SrcID uint32 `json:"src_id"`
	DstID uint32 `json:"dst_id"`
	Src   string `json:"src"`
	Dst   string `json:"dst"`
	Data  string `json:"data"`
}

func srcName(id uint32) string {
	if id == 0 {
		return "SRV"
	}
	return fmt.Sprintf("CLI%d", id)
}

func (p *Proxy) recordPacket(src, dst uint32, data []byte) {
	if *fFile == "" {
		return
	}
	p.emu.Lock()
	defer p.emu.Unlock()
	if p.enc == nil {
		f, err := os.Create(*fFile)
		if err != nil {
			panic(err)
		}
		p.efile = f
		p.enc = json.NewEncoder(f)
	}
	p.enc.Encode(&Record{
		SrcID: src, Src: srcName(src),
		DstID: dst, Dst: srcName(dst),
		Data: hex.EncodeToString(data),
	})
}
