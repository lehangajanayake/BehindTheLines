package main

import (
	"log"
	"net"

	"github.com/lehangajanayake/BehindTheLine/Server/game"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	listener, err := startServer("0.0.0.0", "8080")
	log.Println("Starting the server")
	if err != nil {
		log.Fatal("Error startign the server: ", err)
	}
	g := game.NewGame(1)
	for {
		conn, err := listener.AcceptTCP()
		log.Println("Incoming connection, ", conn.RemoteAddr().String())
		if err != nil {
			log.Fatal("Error accepting connections: ", err)
		}
		if g.AddPlayer(conn) {
			log.Println("Starting a game")
			go g.Run()
			g = game.NewGame(1)
		}
	}

}

func startServer(addr string, port string) (*net.TCPListener, error) {
	TCPaddr, err := net.ResolveTCPAddr("tcp", addr+":"+port)
	if err != nil {
		return nil, err
	}
	return net.ListenTCP("tcp", TCPaddr)

}

func listener(server *net.UDPConn, connchan chan *net.UDPAddr) {
	for {
		buff := make([]byte, 1)
		_, remote, err := server.ReadFromUDP(buff)
		if err != nil {
			log.Fatal("Error Reading UDP Packets: ", err)
		}
		connchan <- remote
	}
}
