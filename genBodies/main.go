package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// This simple command line utility will fit entirely in the main() function.
func main() {
	// First, we just check if there are enough arguments.
	if len(os.Args) < 2 {
		// If not, print an error and exit.
		fmt.Println("genBodies requires at least one argument: the number of points to generate.")
		os.Exit(1)
	}

	// Then, we'll get the number to generate from the command line arguments.
	nBodies, err := strconv.Atoi(os.Args[1])
	// If the user didn't enter a number, exit.
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Now we'll seed the RNG with the current time...
	rand.Seed(time.Now().Unix())
	// and set the maximum deviation from the center in any axis, and the maximum mass
	posMax := 100
	massMax := 5

	// Now we just generate lines in a loop and print them.
	for i := 0; i < nBodies; i++ {
		// Each position is at most posMax away from the origin in that axis.
		// Go doesn't have functions to generate negative random integers,
		// so we generate a positive integer with twice the range and subtract.
		posX := rand.Intn(posMax*2) - posMax
		posY := rand.Intn(posMax*2) - posMax
		posZ := rand.Intn(posMax*2) - posMax
		// On the other hand, mass can't be negative (or zero), so this is easier.
		mass := rand.Intn(massMax-1) + 1
		// Now we print them out in a very simple format with colon seperation.
		fmt.Printf("%d:%d:%d:%d\n", posX, posY, posZ, mass)
	}
}

// Running this program with the command "go run main.go 5" you should see 5 lines with random bodies.
