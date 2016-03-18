package main

import (
	"fmt"
)

func main() {
	elevator.ElevInit()
	fmt.Println("Press STOP button to stop elevator and exit program\n")
	elevator.ElevSetMotorDirection(1)

	for {
		if elevator.ElevGetFloorSensorSignal() == 3 {
			elevator.ElevSetMotorDirection(-1)
		} else if elevator.ElevGetFloorSensorSignal() == 0 {
			elevator.ElevSetMotorDirection(1)
		}

		if elevator.ElevGetStopSignal() {
			elevator.ElevSetMotorDirection(0)
		}
	}
}
