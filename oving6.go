package main

import(
	"fmt"
	"net"
	"os/exec"
	//"strings"
	"strconv"
	"time"
	//"encoding/binary"
)

var baddr string = "129.241.187.255"
var port string = "52134"



func spawnProcess(){
	cmd := exec.Command("gnome-terminal", "-x", "go", "run", "oving6hoved.go")
	out, err := cmd.Output()
	checkErr(err, "spawning process?")
	fmt.Println(string(out))
}


func checkErr(err error, loc string){
	if err != nil{
		fmt.Println("Error " + loc)
	}
}

func backup(count *int){
	fmt.Println("----Backup process running----")
	
	buff := make([]byte, 64)
	sAddr,err := net.ResolveUDPAddr("udp", ":" + port)
	checkErr(err,"resolving UDP address for backup")
	sock, err := net.ListenUDP("udp", sAddr)
	defer sock.Close()
	checkErr(err,"error")
	fmt.Println("Count in backup", *count)
	for true {
		
		sock.SetReadDeadline(time.Now().Add(time.Second*2))
		length,_,err := sock.ReadFromUDP(buff[0:])
		fmt.Println(string(buff))
		if (err!= nil){
			return
		} else{
			*count ,err= strconv.Atoi(string(buff[0:length]))
			checkErr(err, "converting string to int in backup")
			fmt.Println("Slave: read from master = ", *count)
		}
	}
}

func master(count *int){
	fmt.Println("----Master process running----")
	spawnProcess()
	for *count<1000{
		fmt.Println("Master count: ", *count)
		UDPSendAlive(count)
		*count = *count + 1
		time.Sleep(500*time.Millisecond)
	}
}

func UDPSendAlive(count *int){
	mAddr, err := net.ResolveUDPAddr("udp",net.JoinHostPort(baddr,port))
	checkErr(err,"resolving UDP address")
	conn, err := net.DialUDP("udp", nil, mAddr)
	checkErr(err,"setting up UDP connection")
	msg:=strconv.Itoa(*count)
	_,err=conn.Write([]byte(msg))
	checkErr(err,"sending alive signal")
}

func main(){
	count:= 0	
	backup(&count)
	master(&count)
}


