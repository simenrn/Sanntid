package json

import (
	. ".././definitions"
	"encoding/json"
	"fmt"
)

func EncodeMsg(msg MSG) []byte {
	encMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Error encoding msg: ", err)
	}
	return encMsg
}

func DecodeMsg(msg []byte, lenght int) MSG {
	var msg_rec MSG
	err := json.Unmarshal(msg[:lenght], &msg_rec)
	if err != nil {
		fmt.Println("Error decoding msg: ", err)
	}
	return msg_rec
}
