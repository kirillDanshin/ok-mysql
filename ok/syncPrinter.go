package ok

import "github.com/kirillDanshin/myutils"

func syncPrinter(c chan string) {
	for s := range c {
		// continue
		dlogClr.Ln(myutils.Concat(s, "\n\n\n\n"))
	}
}
