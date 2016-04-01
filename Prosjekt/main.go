package main

import (
	"elevator"
	"fmt"
	"statemachine"
	//"time"
)

func main() {
	fmt.Println("Hellluuuuu!")

	elevator.ElevInit()
	fmt.Println("Press STOP button to stop elevator and exit program\n")
	elevator.ElevSetMotorDirection(1)

	for {
		if elevator.ElevGetFloorSensorSignal() == 3 {

			elevator.ElevSetMotorDirection(0)
			elevator.ElevSetStopLamp(1)
			elevator.ElevSetDoorOpenLamp(1)

			elevator.ElevSetMotorDirection(-1)
			elevator.ElevSetStopLamp(0)
			elevator.ElevSetDoorOpenLamp(0)
		} else if elevator.ElevGetFloorSensorSignal() == 0 {
			elevator.ElevSetMotorDirection(0)
			elevator.ElevSetStopLamp(1)

			elevator.ElevSetMotorDirection(1)
			elevator.ElevSetStopLamp(0)
		}

		if elevator.ElevGetStopSignal() == 1 {
			elevator.ElevSetMotorDirection(0)
		}
	}
}
