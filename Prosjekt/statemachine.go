package statemachine

import (
	"fmt"
	. ".././elev"
)

/*type state int

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)*/

var state string

func EventManager() {

	state = "IDLE"
	for {
		switch state {
		case "IDLE":
			fmt.Println("State: Idle")
			if len(Que_Local) != 0 {
				state = "MOVING"
				}
			 // Legg til for eksterne ordre

		case "MOVING":
			fmt.Println("State: Moving")
			ExecuteOrder()
			state = "DOOR_OPEN"

		case "DOOR_OPEN":
			fmt.Println("State: Door open")
			ElevStopAtFloor(Current_Floor)
			state = "IDLE"
		}
	}
}
