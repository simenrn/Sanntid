package statemachine

import (
	. ".././elev"
	. ".././definitions"
	"fmt"
)

func StateMachine(orderEventChannel chan int) {
	fmt.Println("Welcome to the eventmanager")
	Msg.State = IDLE
	for {
		switch Msg.State {
		case IDLE:
			//fmt.Println("State: Idle")
			if NextDirection() != NOTHING {
				Msg.Dirn = NextDirection()
				Msg.State = MOVING
			}

		case MOVING:
			//fmt.Println("State: Moving")
			ElevSetMotorDirection(Msg.Dirn)
			Msg.State = FloorReached()

		case DOOR_OPEN:
			//fmt.Println("State: Door open")
			ElevStopAtFloor(Msg.PrevFloor, orderEventChannel)
			Msg.State = IDLE
		}
	}
}
