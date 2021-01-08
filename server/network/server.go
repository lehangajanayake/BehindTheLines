package network

import (
	//"bufio"
	"log"
	"net"

	//"github.com/MissionImposible/server/game"
)

//StartServer starts the server
func StartServer(conns chan *net.TCPConn) {
	log.Println("Started the server")
	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	listener, err := net.ListenTCP("tcp",  addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Fatalf("Error Accepting the connections: %v", err)
		}
		log.Printf("Connection %+v", conn)
		conns <- conn
	}
}

// //HandleConn
// func HandleConn(conn net.Conn){
	
// 	defer conn.Close()
// 	for {
// 		data, err := bufio.NewReader(conn).ReadString('\n')
// 		// log.Println(string(data))
		
// 		if err != nil{
// 			if err.Error() == "EOF" {
// 				log.Println("Connection disconnected")
// 				return
// 			}
// 			log.Printf("Error reviving data: %v", err)
// 			return
// 		}
// 		_, err = conn.Write([]byte(data))
// 		if err != nil {
// 			log.Printf("Error sending data, %v", err)
// 		}
// 	}
// }

