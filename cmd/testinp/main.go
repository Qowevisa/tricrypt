package main

import (
	"bufio"
	"log"
	"os"

	"golang.org/x/term"
)

func main() {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	reader := bufio.NewReader(os.Stdin)
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			panic(err)
		}
		log.Printf("Read %c ; %d as rune\r\n", r, r)

		// CTRL + C
		if r == 3 {
			break
		}
	}
}
