package driver

import (
	//"fmt"
	"errors"
)

const MOTOR_SPEED int = 2800

const (
	N_FLOORS  int = 4
	N_BUTTONS int = 3
)

type elev_motor_direction_t int

const (
	DIRN_DOWN elev_motor_direction_t = -1
	DIRN_STOP elev_motor_direction_t = 0
	DIRN_UP   elev_motor_direction_t = 1
)

type elev_button_type_t int

const (
	BUTTON_CALL_UP   elev_button_type_t = 0
	BUTTON_CALL_DOWN elev_button_type_t = 1
	BUTTON_COMMAND   elev_button_type_t = 2
)

var (
	lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
		[N_BUTTONS]int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
		[N_BUTTONS]int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
		[N_BUTTONS]int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
		[N_BUTTONS]int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}
	button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
		[N_BUTTONS]int{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
		[N_BUTTONS]int{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		[N_BUTTONS]int{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		[N_BUTTONS]int{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}
)

func ElevInit() bool {

	// Initiate hardware
	if IoInit()==0 {
		return false
	}
	// Zero all floor button lamps
	for i := 0; i < N_FLOORS; i++ {
		// Clearing all Call Down buttons
		if i!=0 {
			ElevSetButtonLamp(BUTTON_CALL_DOWN, i, 0)
		}

		// Clearing all Call UP buttons
		if (i!= N_FLOORS -1) {
			ElevSetButtonLamp(BUTTON_CALL_UP, i, 0)
		}

		ElevSetButtonLamp(BUTTON_COMMAND, i, 0)
		
	}
	ElevSetDoorOpenLamp(0)
    ElevSetFloorIndicator(0)
    return true
}

func ElevSetMotorDirection(dirn elev_motor_direction_t) {
	if dirn == 0 {
		IoWriteAnalog(MOTOR, 0)
	} else if dirn > 0 {
		IoClearBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	} else {
		IoSetBit(MOTORDIR)
		IoWriteAnalog(MOTOR, MOTOR_SPEED)
	}
}

func ElevSetButtonLamp(button elev_button_type_t, floor, value int) {
	if floor < 0 || floor >= N_FLOORS {
		errors.New("Floor is out of range.")
	} else if int(button) < 0 || int(button) >= N_BUTTONS {
		errors.New("Button is out of range")
	} else if value==1 {
		IoSetBit(lamp_channel_matrix[floor][button])
	} else {
		IoClearBit(lamp_channel_matrix[floor][button])
	}
}

func ElevSetFloorIndicator(floor int) {
	if (floor < 0) || (floor >= N_FLOORS) {
		errors.New("Floor is out of range.")
	}
	switch floor {
	case 0:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 1:
		IoClearBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	case 2:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoClearBit(LIGHT_FLOOR_IND2)
	case 3:
		IoSetBit(LIGHT_FLOOR_IND1)
		IoSetBit(LIGHT_FLOOR_IND2)
	}
}

func ElevSetDoorOpenLamp(value int) {
	if value==1 {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}

func ElevSetStopLamp(value int) {
	if value==1 {
		IoSetBit(LIGHT_STOP)
	} else {
		IoClearBit(LIGHT_STOP)
	}
}

func ElevGetButtonSignal(button elev_button_type_t, floor int) int {
	if IoReadBit(button_channel_matrix[floor][button]) == 1 {
		return 1
	}else{
		return 0
	}
}

func ElevGetFloorSensorSignal() int {
	if IoReadBit(SENSOR_FLOOR1)==1 {
		return 0
	} else if IoReadBit(SENSOR_FLOOR2)==1 {
		return 1
	} else if IoReadBit(SENSOR_FLOOR3)==1 {
		return 2
	} else if IoReadBit(SENSOR_FLOOR4)==1 {
		return 3
	} else {
		return -1
	}
}

func ElevGetStopSignal() int {
	if IoReadBit(STOP) == 1{
		return 1
	}else{
		return 0
	}
}
