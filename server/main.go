package main

import (
	"log"
	"net"

	"github.com/lehangajanayake/BehindTheLine/Server/game"
)


func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	listener, err := startServer("0.0.0.0", "8080")
	if err != nil {
		log.Fatal("Error resolving the udp addr: ", err)
	}
	g := &game.Game{PlayerNum: 2}
	for {
		conn, err := listener.AcceptTCP()
		log.Println("Incoming connection, ", conn.RemoteAddr().String())
		if err != nil {
			log.Fatal("Error accepting connections: ", err)
		}
		if g.AddPlayer(conn){
			log.Println("Starting a game")
			go g.Run()
			g = &game.Game{PlayerNum: 2}
		}
	}
	

}

func startServer(addr string, port string)(*net.TCPListener, error){
	TCPaddr, err := net.ResolveTCPAddr("tcp", addr + ":" + port)
	if err != nil {
		return nil, err
	}
	log.Println("Starting the Server, ", TCPaddr)
	return net.ListenTCP("tcp", TCPaddr)
	
}

func listener(server *net.UDPConn, connchan chan *net.UDPAddr){
	for {
		buff := make([]byte, 1)
		_, remote, err := server.ReadFromUDP(buff)
		if err != nil {
			log.Fatal("Error Reading UDP Packets: ", err)
		}
		connchan <- remote
	}
}