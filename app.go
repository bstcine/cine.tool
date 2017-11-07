package main

import (
	"os"
	"bufio"
)

func main() {
	MachiningImage(false,CheckHasMagick())

	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		if input.Text() == "end" {
			break
		}
	}
}