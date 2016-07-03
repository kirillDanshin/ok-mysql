package ok

import (
	"fmt"
	"log"
	"net"

	"github.com/kirillDanshin/ok-mysql/defaults"
)

func newInst(cfg *Config) (*Instance, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("Can't use empty address")
	}
	var (
		snaplen = cfg.SnapshotLength
	)
	if snaplen == 0 {
		log.Printf("using default SnapshotLength %v for %v", defaults.SnapLen, cfg.Address)
		if defaults.SnapLen == 0 {
			return nil, fmt.Errorf("config.SnapshotLength equals to zero and no valid default value provided")
		}
		snaplen = defaults.SnapLen
	}

	addr, err := net.ResolveTCPAddr(defaults.Net, cfg.Address)
	if err != nil {
		return nil, err
	}

	return &Instance{
		Addr:    addr,
		SnapLen: snaplen,
	}, nil
}
