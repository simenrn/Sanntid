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

		fmt.Println("Interne ordre: ", Que_Local)
		fmt.Println("Opp ordre: ", Que_Global_Up)
		fmt.Println("Ned ordre: ", Que_Global_Down)
		fmt.Println("Msg.PrevFloor: ", Msg.PrevFloor)
		fmt.Println("Msg.Dirn: ", Msg.Dirn)
		fmt.Println("Msg.State", Msg.State)
		time.Sleep(time.Second * 3)

		if ElevGetStopSignal() == 1 {
			ElevSetMotorDirection(0)
		}
	}
}
