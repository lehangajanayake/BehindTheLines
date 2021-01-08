package main

import (
	"bufio"
	"log"
	"net"

	//"time"

	"github.com/MissionImposible/server/game"
	"github.com/MissionImposible/server/network"
)



func main(){
	conns := make(chan *net.TCPConn)
	go network.StartServer(conns)
	done := make(chan bool)
	errchan := make(chan error)
	var que []*net.TCPConn
	go func() {
		for{
			if err := <- errchan; err != nil{
				log.Println("Error in the lobby")
				que = append(que[:0], que[1:]...)
			}
		}
	}()
	for v := range conns {
		go lobby(v, done, errchan)
		println(len(que))
		if len(que) == 1{
			que = append(que, v)
			g := game.NewGame(que)
			done <- true
			log.Println(g.Run())
			que = nil
		}else{
			que = append(que, v)
			
		}
		
	}


}

func lobby(conn *net.TCPConn, done chan bool, errs chan error){
	conn.SetKeepAlive(true)
	writer := bufio.NewWriter(conn)
	for {
		_, err := writer.WriteString("Lobby\n")
		if err != nil {
			errs <- err
			return
		}
		err = writer.Flush()
		if err != nil {
			errs <- err
			return
		}
		if done := <- done; done{
			break
		}
		//time.Sleep(time.Second *1)
	}
}
