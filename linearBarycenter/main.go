package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// In this video, we'll implement a non-concurrent, non-linear version of the barycenter finder.
// It won't be particularly fast, but it will get the job done.
// ST
// First, though, we need some data to operate on, so we'll write a command line utility
// to generate random bodies. After that's done, we can write the actual barycenter program.
// ST
// Pull up your editor and let's get coding!
// ST (editor - genBodies)
// Once back here:
// Now it's time for the actual barycenter finder.
// The first thing to do is to create the struct we'll use to represent a body.

// We'll call it a MassPoint. It's a point in 3-space plus mass information.
// The MassPoint is the primary datastructure we'll pass around through the program.
type MassPoint struct {
	x, y, z, mass float64
}

// We can define a function that just adds them together.
// It'll add each coordinate, and the masses.
func addMassPoints(a MassPoint, b MassPoint) MassPoint {
	// We just return a new MassPoint where we've added in each coordinate
	return MassPoint{
		a.x + b.x,
		a.y + b.y,
		a.z + b.z,
		a.mass + b.mass,
	}
}

// Using that function, we can create a function that averages them.
func avgMassPoints(a MassPoint, b MassPoint) MassPoint {
	// All we have to do is add them together and divide by two.
	sum := addMassPoints(a, b)
	// So we divide by two in each coordinate
	return MassPoint{
		sum.x / 2,
		sum.y / 2,
		sum.z / 2,
		// But not in the mass
		sum.mass,
	}
}

// Then, we need a function that maps them to a different point in space by mass, as we discussed.
func toWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x * a.mass,
		a.y * a.mass,
		a.z * a.mass,
		a.mass,
	}
}

// And we need a function that takes them back.
func fromWeightedSubspace(a MassPoint) MassPoint {
	return MassPoint{
		a.x / a.mass,
		a.y / a.mass,
		a.z / a.mass,
		a.mass,
	}
}

// Finally, we'll write a function which takes a pair of mass points and returns the
// weighted average.
func avgMassPointsWeighted(a MassPoint, b MassPoint) MassPoint {
	// First we calculate the weighted version of both mass points
	aWeighted := toWeightedSubspace(a)
	bWeighted := toWeightedSubspace(b)
	return fromWeightedSubspace(avgMassPoints(aWeighted, bWeighted))
}

// Now, on to the actual application code. First, we'll define two useful helper functions.

// We'll create a function that just handles errors by aborting.
func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// And another that will close a file, so we can defer that operation.
func closeFile(fi *os.File) {
	err := fi.Close()
	handle(err)
}

// Now comes the actual bulk of our program, in the main function.
func main() {
	// Check arguments. We need exactly two (the executable name and one user-provided argument).
	if len(os.Args) != 2 {
		// If there are too many or not enough, abort.
		fmt.Println("Incorrect number of arguments!")
		os.Exit(1)
	}

	// Then, we'll open the input file with os.Open
	file, err := os.Open(os.Args[1])
	// Handle a possible error using our error handler
	handle(err)
	// And finally defer the closing of the file,
	// so even if the program aborts the file will still get closed.
	defer closeFile(file)

	// Now we need initial buffer for the MassPoints, which we'll load from the file.
	var masspoints []MassPoint

	// We'll time how long it takes to load them, just for comparison.
	startLoading := time.Now()
	// Let's make an infinite loop for loading, and then break out when there are no more points
	// to read.
	for {
		// We'll create a variable to hold the new point
		var newMassPoint MassPoint
		// Then we'll use fmt.Fscanf to load a single line from the file
		_, err = fmt.Fscanf(file, "%f:%f:%f:%f", &newMassPoint.x, &newMassPoint.y, &newMassPoint.z, &newMassPoint.mass)
		// If we got an EOF error, there are no more points to load
		if err == io.EOF {
			break
			// On other errors, we can just skip the line.
		} else if err != nil {
			continue
		}
		// Finally, we'll use append() to add the point into the list of points.
		masspoints = append(masspoints, newMassPoint)
	}

	// Now we'll report how many points we loaded
	fmt.Printf("Loaded %d values from file in %s.\n", len(masspoints), time.Since(startLoading))
	// And we should check that there are actually enough values.
	if len(masspoints) <= 1 {
		// If there aren't enough, we'll create an error and pass it to our error handler.
		handle(errors.New("Insufficient number of values; there must be at least one "))
	}

	// We also want to time the calculation itself, so we'll start a timer.
	startCalculation := time.Now()

	// Now, we'll make a loop. It'll run until there's exactly one point left.
	for len(masspoints) != 1 {
		// Each loop will need a new array of MassPoints
		var newMasspoints []MassPoint

		// We loop over the current list of bodies by twos
		for i := 0; i < len(masspoints)-1; i += 2 {
			// Adding the results of the averaging to the new array of mass points
			newMasspoints = append(newMasspoints, avgMassPointsWeighted(masspoints[i], masspoints[i+1]))
		}

		// Then we check to make sure we didn't leave off one
		if len(masspoints)%2 != 0 {
			newMasspoints = append(newMasspoints, masspoints[len(masspoints)-1])
		}

		// Finally, we need to switch out the old array with the new one
		masspoints = newMasspoints
	}

	// Once the loop is done we need the one remaining virtual body
	systemAverage := masspoints[0]

	// And then we'll print out the result in a pretty way.
	fmt.Printf("System barycenter is at (%f, %f, %f) and the system's mass is %f.\n",
		systemAverage.x,
		systemAverage.y,
		systemAverage.z,
		systemAverage.mass)
	// Finally, we just want to print out the time the calculation has taken.
	fmt.Printf("Calculation took %s.\n", time.Since(startCalculation))
}
