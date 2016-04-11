package statemachine

import (
	"fmt"
	. ".././elevator"
)

type state int

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)

var DestinationFloor int

func EventManager() {

	state = IDLE
	for {
		switch state {
		case IDLE:
			fmt.Println("State: Idle")
			if NewInternalOrder = true {
				for i:= 0; i<N_FLOORS; i++{
					if Internal_orders[i] == 1{
						DestinationFloor = i
						state = MOVING
						break
					}

				}
			} // Legg til for eksterne ordre

		case MOVING:
			fmt.Println("State: Moving")
			ExecuteOrder(DestinationFloor)
			state = DOOR_OPEN

		case DOOR_OPEN:
			fmt.Println("State: Door open")
			ElevStopAtFloor(DestinationFloor)
			state = IDLE
		}
	}
}
