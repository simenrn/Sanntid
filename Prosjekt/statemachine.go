package statemachine

import (
	"fmt"
	"src/elevator"
)

type state int

const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)

func EventManager() {

	state = IDLE

	switch state {
	case IDLE:
		fmt.Println("State: Idle")

	case MOVING:
		fmt.Println("State: Moving")
		elevator.ElevetMotorDirection(dirn)
		for elevator.ElevGetFloorSensorSignal() != floor {
		}
		state = DOOR_OPEN

	case DOOR_OPEN:
		fmt.Println("State: Door open")
		elevator.ElevStopAtFloor(floor)
		state = IDLE

	}
}
