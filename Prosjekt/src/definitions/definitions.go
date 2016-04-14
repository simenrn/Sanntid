package definitions

import (
	. ".././driver"
	. ".././udp"
)


var (
	Lamp_channel_matrix = [N_FLOORS][N_BUTTONS]int{
		[N_BUTTONS]int{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
		[N_BUTTONS]int{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
		[N_BUTTONS]int{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
		[N_BUTTONS]int{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4},
	}
	Button_channel_matrix = [N_FLOORS][N_BUTTONS]int{
		[N_BUTTONS]int{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
		[N_BUTTONS]int{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
		[N_BUTTONS]int{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
		[N_BUTTONS]int{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4},
	}
)

const (
	MOTOR_SPEED int = 2800

	N_FLOORS  int = 4
	N_BUTTONS int = 3

	IDLE      int = 0
	MOVING    int = 1
	DOOR_OPEN int = 2

	DIRN_DOWN int = -1
	DIRN_STOP int = 0
	DIRN_UP   int = 1

	ADD_ORDER    int = 0
	REMOVE_ORDER int = 1
	NOTHING int = 5

	BUTTON_CALL_UP   int = 0
	BUTTON_CALL_DOWN int = 1
	BUTTON_COMMAND   int = 2

	ON  int = 1
	OFF int = 0
)

type MSG struct {
	State           int
	Dirn            int
	PrevFloor       int
	PrevDirn        int
	MessageType     int
	Que_Local       [N_FLOORS]int
	Que_Global_Up   [N_FLOORS]int
	Que_Global_Down [N_FLOORS]int
}

var Msg = MSG{}
var buff = make([]byte, 1024)
var Udp_Msg = Udp_message{"broadcast", buff, 1024}

var Current_Floor int

const (
	LOCAL_LISTEN_PORT		int = 20267
	BROADCAST_LISTEN_PORT	int = 30267
	MESSAGE_SIZE	 		int = 1024
)

