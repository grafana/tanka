package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/pkg/errors"
)

// Confirm asks the user for confirmation
func Confirm(msg, approval string) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(msg)
	fmt.Printf("Please type '%s' to confirm: ", approval)
	read, err := reader.ReadString('\n')
	if err != nil {
		return errors.Wrap(err, "reading from stdin")
	}
	if read != approval+"\n" {
		return errors.New("aborted by user")
	}
	return nil
}
