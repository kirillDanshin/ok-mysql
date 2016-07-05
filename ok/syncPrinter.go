package ok

import "github.com/kirillDanshin/dlog"

func syncPrinter(c chan string) {
	for s := range c {
		dlog.Ln(s)
		dlog.P("\n\n\n\n")
	}
}
