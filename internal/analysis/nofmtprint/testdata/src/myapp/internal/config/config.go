package config

import (
	"fmt"
	"os"
)

func LoadFail(err error) {
	fmt.Fprintf(os.Stderr, "config: %v\n", err)
}
