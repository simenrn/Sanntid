package main

import (
	. "./elev"
	. "./statemachine"
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hellluuuuu!")

	ElevInit()

	fmt.Println("Press STOP button to stop elevator and exit program\n")
	//ElevSetMotorDirection(1)
	go GetOrders()
	go ElevLights()
	go EventManager()

	for {

		fmt.Println(Internal_orders)
		fmt.Println(External_orders)
		fmt.Println("lengden paa kooen: ", len(Que_Local))
		if len(Que_Local) > 1 {
			fmt.Println("ordre nummer 1: ", Que_Local[0])
			fmt.Println("ordre nummer 2: ", Que_Local[1])
		}
		fmt.Println("Current floor: ", Current_Floor)
		time.Sleep(time.Second * 3)
		/*
			if ElevGetFloorSensorSignal() == 3 {
				ElevSetMotorDirection(-1)
				ElevSetDoorOpenLamp(1)
			} else if ElevGetFloorSensorSignal() == 0 {

				ElevSetMotorDirection(1)

			}
		*/
		if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
		}
	}
}
