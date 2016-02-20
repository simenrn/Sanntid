package main

import(
	"fmt"
	//"time"
	"net"
	//"runtime"
)

const (
	host = "129.241.187.23"
	localIP = "129.241.187.153"
	TCPClientPort = "20021"
	TCPServPort_fix = "34933"
	TCPServPort_term = "33546"
	
)

func TCPconnect(done chan bool){
	
	laddr, err := net.ResolveTCPAddr("tcp", ":" + TCPClientPort)
	if err != nil{
		fmt.Println("Failed to resolve address for client")
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if err!=nil{
		fmt.Println("Failed to create listener on local")
	}
	
	conn, err := listener.AcceptTCP()
	if err!=nil{
		fmt.Println("Failed to set up set up new connection")
	}
	
	fmt.Println("connected to workspace 21")
	conn.Close()
	done<-true
	
}

func sendMsg(msg string, conn *net.TCPConn, bytes_sent int, err error){
	
	bytes_sent,err = conn.Write([]byte(msg))
	conn.Write([]byte{0})
	if err!=nil{
		fmt.Println("Failed to send msg2")
	}
}

func resvAndPrintMsg(conn *net.TCPConn, bytes_read int, err error, buffer []byte){
	bytes_read, err = conn.Read(buffer)
	fmt.Println("Read2", bytes_read, "bytes: ", string(buffer))
	
}

func main(){
	
	raddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(host,TCPServPort_term))
	if err!=nil{
		fmt.Println("Failed to resolve address for server")
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err!=nil{
		fmt.Println("Failed to connect to server")
	}
	buffer := make([]byte, 1024)
	bytes_read, err := conn.Read(buffer)
	if err!=nil{
		fmt.Println("Failed to read from buffer")
	}
	fmt.Println("Read", bytes_read, "bytes: ", string(buffer))

	done := make(chan bool)
	go TCPconnect(done)
	bytes_sent, err := conn.Write([]byte("Connect to: 129.241.187.153:20021"))
	conn.Write([]byte{0})
	if err!=nil{
		fmt.Println("Failed to send msg")
	}
	fmt.Println("Sent", bytes_sent, "bytes")
	
	sendMsg("Halla dude",conn, bytes_sent, err)
	resvAndPrintMsg(conn, bytes_read, err, buffer)

	sendMsg("Klar for deadlifting man?",conn, bytes_sent, err)
	resvAndPrintMsg(conn, bytes_read, err, buffer)
	sendMsg("det blir sick sant????",conn, bytes_sent, err)
	resvAndPrintMsg(conn, bytes_read, err, buffer)
	
	<-done
	err = conn.Close()
	if err != nil{
		fmt.Println("Failed to close connection")
	}
}
