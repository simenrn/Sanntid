package main

import (
	"fmt"
	"elevator"
	//"time"
)

func main() {
	fmt.Println("Hellluuuuu!")

	driver.ElevInit()
	fmt.Println("Press STOP button to stop elevator and exit program\n")
	driver.ElevSetMotorDirection(1)

	for {
		if driver.ElevGetFloorSensorSignal() == 3 {
			
			driver.ElevSetMotorDirection(0)
			driver.ElevSetStopLamp(1)
			driver.ElevSetDoorOpenLamp(1)
			
			driver.ElevSetMotorDirection(-1)
			driver.ElevSetStopLamp(0)
			driver.ElevSetDoorOpenLamp(0)
		} else if driver.ElevGetFloorSensorSignal() == 0 {
			driver.ElevSetMotorDirection(0)
			driver.ElevSetStopLamp(1)
			
			driver.ElevSetMotorDirection(1)
			driver.ElevSetStopLamp(0)
		}

		if driver.ElevGetStopSignal() == 1 {
			driver.ElevSetMotorDirection(0)
		}
	}

}
