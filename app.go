package main

import (
	"os"
	"bufio"
)

func main() {
	MachiningImage(true)

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if input.Text() == "end" {
			break
		}
	}
}