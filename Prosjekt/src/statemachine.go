package statemachine

import (
	. ".././definitions"
	. ".././elev"
	"fmt"
)

func StateMachine(orderEventChannel chan int) {
	fmt.Println("Welcome to the eventmanager")
	ElevatorList[0].State = IDLE
	for {
		switch ElevatorList[0].State {
		case IDLE:

			//fmt.Println("State: Idle")
			if NextDirection() == DIRN_UP || NextDirection() == DIRN_DOWN {
				ElevatorList[0].Dirn = NextDirection()
				ElevatorList[0].State = MOVING
			}
			if NextDirection() == DIRN_STOP {
				ElevatorList[0].Dirn = NextDirection()
				ElevatorList[0].State = DOOR_OPEN
			}

		case MOVING:
			//mt.Println("State: Moving")
			ElevSetMotorDirection(ElevatorList[0].Dirn)
			ElevatorList[0].State = FloorReached()

		case DOOR_OPEN:
			//fmt.Println("State: Door open")
			ElevStopAtFloor(ElevatorList[0].CurrentFloor, orderEventChannel)
			ElevatorList[0].State = IDLE
		}
	}
}
