package main

import (
	"fmt"
	"os"

	"github.com/MaikovskiyS/TraderBot/internal/app"
)

func main() {
	if err := app.Run(); err != nil {
		fmt.Printf("Service exited with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Service exited gracefully")
}
