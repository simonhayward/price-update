package main

import (
	"log"

	"github.com/simonhayward/price-update/priceupdate"
)

func main() {
	if err := priceupdate.Run(); err != nil {
		log.Fatalf("error: %s\n", err)
	}
}
