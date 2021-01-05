package profari

import (
	"bufio"
	"fmt"
	"os"
)

// Pause is an helper that pauses the program until the user hit Enter.
// Pause is best used for debugging your test suite.
func Pause() {
	fmt.Println("Program Paused!")
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
