package elevator

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

const (
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
	if IoInit() == 0 {
		errors.New("Error initializing elev_init..")
		return false
	}

	for i := 0; i < N_FLOORS; i++ {
		for b = 0; b < N_BUTTONS; b++ {
			ElevSetButtonLamp(b, i, 0)
		}
	}

}

func ElevSetMotorDirection(elev_motor_direction_t dirn) {
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

func ElevSetButtonLamp(elev_button_type_t button, floor, value int) {
	if floor < 0 || floor >= N_FLOORS {
		errors.New("Floor is out of range.")
	} else if button < 0 || button >= N_BUTTONS {
		errors.New("Button is out of range")
	} else if value {
		IoSetBit(lamp_channel_matrix[floor][button])
	} else {
		IoClearBit(lamp_channel_matrix[floor][button])
	}
}

func ElevFloorIndicator(int floor) {
	if floor < 0 || floor >= N_FLOORS {
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

func ElevSetDoorOpenLamp(int value) {
	if value {
		IoSetBit(LIGHT_DOOR_OPEN)
	} else {
		IoClearBit(LIGHT_DOOR_OPEN)
	}
}

func ElevSetStopLamp(int value) {
	if value {
		IoSetBit(LIGHT_STOP)
	} else {
		IoClearBit(LIGHT_STOP)
	}
}

func ElevGetButtonSignal(elev_button_type_t button, int floor) int {
	if floor < 0 || floor >= N_FLOORS {
		errors.New("Floor is out of range.")
	} else if button < 0 || button >= N_BUTTONS {
		errors.New("Button is out of range")
	} else {
		return IoReadBit(button_channel_matrix[floor][button])
	}
}

func ElevGetFloorSensorSignal() int {
	if IoReadBit(SENSOR_FLOOR1) {
		return 0
	} else if IoReadBit(SENSOR_FLOOR2) {
		return 1
	} else if IoReadBit(SENSOR_FLOOR3) {
		return 2
	} else if IoReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}
