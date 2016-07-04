package ok

import "fmt"

var (
	syncPrint = make(chan string, 128)
)

func syncPrinter(c chan string) {
	for s := range c {
		fmt.Println(s)
		fmt.Print("\n\n\n\n")
	}
}
