package main

import (
	"bufio"
	"log"
	"net"
)


func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panicf("Error Accepting the connections: %w", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn){
	
	defer conn.Close()
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		log.Println(string(data))
		if err != nil{
			if err.Error() == "EOF" {
				log.Println("Connection disconnected")
				return
			}
			log.Printf("Error reviving data: %w", err)
			return
		}
	}
}

