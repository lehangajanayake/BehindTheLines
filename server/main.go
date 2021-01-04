package main

import (
	//"log"
	"net"
	"time"

	"github.com/MissionImposible/server/game"
	"github.com/MissionImposible/server/network"
)



func main(){
	conns := make(chan net.Conn)
	go network.StartServer(conns)
	done := make(chan bool)
	errchan := make(chan error)
	var que []net.Conn
	for v := range conns {
		//v.SetDeadline(time.Now().Add(time.Second))
		println(len(que))
		if len(que) == 1{
			println("hello")
			que = append(que, v)
			g := game.NewGame(que)
			done <- true
			g.Run()
			que = nil
		}else{
			que = append(que, v)
			go lobby(v, done, errchan)
		}
		// if err := <- errchan; err != nil{
		// 	log.Println("Error in the lobby")
		// 	que = append(que[:0], que[1:]...)
		// }
	}


}

func lobby(conn net.Conn, done chan bool, errs chan error){
	for {
		_, err := conn.Write([]byte("Lobby\n"))
		if err != nil {
			errs <- err
		}
		if done := <- done; done{
			break
		}
		time.Sleep(1000)
	}
}
