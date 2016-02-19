package main

import(
	"fmt"
	"net"
	"time"
	//"runtime"
)



const (
	host = "129.241.187.23"
	udpPort = "20021"
	port_t = "30000"
)

func udpReceive(port string) {
	buff := make([]byte, 1024)
	addr, _ := net.ResolveUDPAddr("udp", ":" + port)
	fmt.Println(addr)
	sock, _ := net.ListenUDP("udp", addr)
	fmt.Println(sock)
	for {
		_,_, err:=sock.ReadFromUDP(buff)
		if err != nil {
			fmt.Println(err)
		} 
		fmt.Println(string(buff)+"lolz")
	}
}

func udpSend(){
	raddr,err:= net.ResolveUDPAddr("udp",net.JoinHostPort(host,udpPort))
	if err!=nil{
		fmt.Println("Failed to resolve adress for: " + udpPort)
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err!=nil{
		fmt.Println("Failed to connect")
	}
	go udpReceive(udpPort)
	
	for{
		time.Sleep(1000*time.Millisecond)
		conn.Write([]byte("heisann"))
		fmt.Println("Msg sent")
	}

	
}

func main(){
	udpSend()

}
