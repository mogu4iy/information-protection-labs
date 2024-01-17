package term

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"strings"
)

func ReadVar(key string, hide bool) (value string) {
	for {
		fmt.Printf("Input %s: ", key)
		if hide {
			password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", key, err)
				os.Exit(1)
			}
			value = string(password)
		}else {
			value = readInput()
		}
		if value == "" {
			fmt.Printf("\nInvalid %s. Try again.\n", key)
		} else {
			break
		}
	}
	fmt.Println()
	return
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
