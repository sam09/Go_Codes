package main

import (
	"fmt"
	"flag"
	"net"
	"io"
)

const maxConns = 5

func checkError (err error, info string ) {
	if err != nil {
		panic("Error " + info + " " + err.Error() )
	}
}

func initServer (hostAndPort string) *net.TCPListener {

	serverAddr, err := net.ResolveTCPAddr ( "tcp", hostAndPort)
	checkError (err, "Resolving Address failed")

	listener, err := net.ListenTCP("tcp", serverAddr )
	checkError (err, "Listen failed ")

	return listener
}

func handler (conn net.Conn , conns []net.Conn, client string , full chan int, empty chan int) {

	data := make([]byte , 1024)

	for {
		count , err := conn.Read (data)
		if err != nil {
			fmt.Printf("%s left\n", client )
			
			for i := range conns {
				if conns[i] == conn {
					conns[i] = nil
					break
				}
			}
			empty <- 1
			<- full
			return
		}
		msg := client + ": " + string (data[0:count])

		for _, i := range conns {
			if i != conn && i != nil {
				io.WriteString(i, msg)
			}
		}

	}
}


func main () {

	flag.Parse()
	if flag.NArg() != 2 {
		panic ("Usage : host port")
	}

	hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1) )
	listener := initServer(hostAndPort)

	conns := make ( []net.Conn, maxConns )
	count := 0
	
	empty := make(chan int, maxConns)
	full := make (chan int, maxConns)

	for {
		conn, err := listener.Accept()
		checkError( err, "Accept: ")
		flag := true
		full <- 1		
		for {
			select {
				case v := <- empty:
					count -= v
				default:
					flag = false
					break
			}
			if ! flag {
				break
			}
		}

		client := fmt.Sprintf("Client %d", count+1)
		conns[count] = conn
		count++
		fmt.Println(client + " joined")
		go handler(conn, conns, client, full, empty)
	}
}