package udp

import (
	. ".././definitions"
	"fmt"
	"net"
	"strconv"
)

var Laddr *net.UDPAddr //Local address
var Baddr *net.UDPAddr //Broadcast address

func Udp_init(localListenPort, broadcastListenPort, message_size int, send_ch, receive_ch chan Udp_message) (err error) {
	//Generating broadcast address
	//fmt.Println("Generating broadcast address")
	Baddr, err = net.ResolveUDPAddr("udp4", "255.255.255.255:"+strconv.Itoa(broadcastListenPort))
	if err != nil {
		return err
	}

	//Generating locaLaddress
	//fmt.Println("Generating local address")
	tempConn, err := net.DialUDP("udp4", nil, Baddr)
	defer tempConn.Close()
	tempAddr := tempConn.LocalAddr()
	Laddr, err = net.ResolveUDPAddr("udp4", tempAddr.String())
	if len(ElevatorList)>0{
		ElevatorList[0].MyIP = Laddr.IP
	} else {
		Msg.MyIP = Laddr.IP
	}
	Laddr.Port = localListenPort

	//Creating local listening connections
	localListenConn, err := net.ListenUDP("udp4", Laddr)
	if err != nil {
		return err
	}

	//Creating listener on broadcast connection
	broadcastListenConn, err := net.ListenUDP("udp", Baddr)
	if err != nil {
		localListenConn.Close()
		return err
	}

	go udp_receive_server(localListenConn, broadcastListenConn, message_size, receive_ch)
	go udp_transmit_server(localListenConn, broadcastListenConn, send_ch)

	//	fmt.Printf("Generating local address: \t Network(): %s \t String(): %s \n", Laddr.Network(), Laddr.String())
	//	fmt.Printf("Generating broadcast address: \t Network(): %s \t String(): %s \n", Baddr.Network(), Baddr.String())
	return err
}

func udp_transmit_server(lconn, bconn *net.UDPConn, send_ch chan Udp_message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_transmit_server: %s \n Closing connection.", r)
			lconn.Close()
			bconn.Close()
		}
	}()

	var err error
	var n int

	for {
		msg := <-send_ch
		//		fmt.Printf("Writing %s \n", msg.Data)
		if msg.Raddr == "broadcast" {
			n, err = lconn.WriteToUDP([]byte(msg.Data), Baddr)
		} else {
			raddr, err := net.ResolveUDPAddr("udp", msg.Raddr+strconv.Itoa(LOCAL_LISTEN_PORT))
			if err != nil {
				fmt.Printf("Error: udp_transmit_server: could not resolve raddr\n")
				panic(err)
			}
			n, err = lconn.WriteToUDP([]byte(msg.Data), raddr)
		}
		if err != nil || n < 0 {
			fmt.Printf("Error: udp_transmit_server: writing\n")
			panic(err)
		}
		//		fmt.Printf("udp_transmit_server: Sent %s to %s \n", msg.Data, msg.Raddr)
	}
}

func udp_receive_server(lconn, bconn *net.UDPConn, message_size int, receive_ch chan Udp_message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_receive_server: %s \n Closing connection.", r)
			lconn.Close()
			bconn.Close()
		}
	}()

	bconn_rcv_ch := make(chan Udp_message)
	lconn_rcv_ch := make(chan Udp_message)

	go udp_connection_reader(lconn, message_size, lconn_rcv_ch)
	go udp_connection_reader(bconn, message_size, bconn_rcv_ch)

	for {
		select {

		case buf := <-bconn_rcv_ch:
			receive_ch <- buf

		case buf := <-lconn_rcv_ch:
			receive_ch <- buf
		}
	}
}

func udp_connection_reader(conn *net.UDPConn, message_size int, rcv_ch chan Udp_message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("ERROR in udp_connection_reader: %s \n Closing connection.", r)
			conn.Close()
		}
	}()

	for {
		buf := make([]byte, message_size)
		//fmt.Printf("udp_connection_reader: Waiting on data from UDPConn\n")
		n, raddr, err := conn.ReadFromUDP(buf)
		//		fmt.Printf("udp_connection_reader: Received %s from %s \n", string(buf), raddr.String())
		if err != nil || n < 0 {
			fmt.Printf("Error: udp_connection_reader: reading\n")
			panic(err)
		}
		rcv_ch <- Udp_message{Raddr: raddr.String(), Data: buf, Length: n}
	}
}
