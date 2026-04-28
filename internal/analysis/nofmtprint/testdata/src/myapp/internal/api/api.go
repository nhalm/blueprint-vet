package api

import (
	"fmt"
	"os"
)

func Bad() {
	fmt.Println("user created")                  // want `use canonlog instead of fmt\.Println`
	fmt.Printf("got %d\n", 5)                    // want `use canonlog instead of fmt\.Printf`
	fmt.Fprintf(os.Stderr, "warn: %s\n", "x")    // want `use canonlog instead of fmt\.Fprintf`
	_ = fmt.Sprintf("user=%s", "alice")          // want `use canonlog instead of fmt\.Sprintf`
}
