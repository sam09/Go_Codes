package main

import (
	"fmt"
	"net"
	"io"
	"os"
	"flag"
)
func read ( conn net.Conn ) {

	data := make([]byte , 1024)

	for {
		count , err := conn.Read (data)
		if err != nil {
			fmt.Println("Error " + err.Error() )
			os.Exit(1)
		}
		fmt.Println( string (data[0:count]) )
	}

}
func main () {

	flag.Parse()

	hostAndPort := fmt.Sprintf("%s:%s", flag.Arg(0), flag.Arg(1))

	conn, err := net.Dial ("tcp", hostAndPort)
	
	if err != nil {
		fmt.Println("Error Dialing")
		return
	}

	go read(conn)
	var msg string

	for {
		fmt.Scanf("%s", &msg)
		io.WriteString(conn, msg)
	}

}