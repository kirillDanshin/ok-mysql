package ok

import (
	"github.com/kirillDanshin/dlog"
	"github.com/kirillDanshin/myutils"
)

func syncPrinter(c chan string) {
	for s := range c {
		// continue
		dlog.Ln(myutils.Concat(s, "\n\n\n\n"))
	}
}
