package ok

import (
	"fmt"
	"net"

	"github.com/kirillDanshin/dlog"
	"github.com/kirillDanshin/ok-mysql/defaults"
)

func newInst(cfg *Config) (*Instance, error) {
	if cfg.Address == "" {
		return nil, fmt.Errorf("Can't use empty address")
	}
	var (
		snaplen = cfg.SnapshotLength
		Lazy    = cfg.Lazy
	)
	if snaplen == 0 {
		dlog.F("using default SnapshotLength %v for %v", defaults.SnapLen, cfg.Address)
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
		Lazy:    Lazy,
	}, nil
}
