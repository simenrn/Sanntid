package main

import (
	. "./elev"
	"fmt"
	"time"
	. "./statemachine"
	//"time"
)



func main() {
	fmt.Println("Hellluuuuu!")

	ElevInit()

	fmt.Println("Press STOP button to stop elevator and exit program\n")
	//ElevSetMotorDirection(1)
	go GetOrders()
	go ElevLights()
	go EventManager()

	/*
	for {
		fmt.Println("Jeg er her modder")
		fmt.Println(Internal_orders)
		fmt.Println(External_orders)
		fmt.Println("lokalkoooo: ", Que_Local)
		time.Sleep(time.Second *3)
		if ElevGetFloorSensorSignal() == 3 {
			ElevSetMotorDirection(-1)
			ElevSetDoorOpenLamp(1)
		} else if ElevGetFloorSensorSignal() == 0 {

			ElevSetMotorDirection(1)

		}

		if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
		}
	}*/
	if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
}
