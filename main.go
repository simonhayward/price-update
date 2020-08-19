package main

import (
	"fmt"

	"github.com/simonhayward/price-update/priceupdate"
)

func main() {
	if err := priceupdate.Run(); err != nil {
		fmt.Printf("error: %s\n", err)
	}
}
