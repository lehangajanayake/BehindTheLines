package game

import (
	"bufio"
	"log"
	"strconv"
	"sync"

	"errors"
	// "fmt"
	//"log"
	"net"
	// "strconv"
	// "strings"
	// "sync"
)

//Game is the main game
type Game struct{
	Players []*Player
}


//NewGame created a new game
func NewGame(conn []*net.TCPConn)*Game{
	return &Game{
		Players: []*Player{
			{ID: 1, X:200, Y:200, FacingFront: true, Guard: true, Connection: conn[0], ReadBuff: bufio.NewReader(conn[0]), WriteBuff: bufio.NewWriter(conn[0])},
			{ID: 2, X:200, Y:400, FacingFront: true, Guard: false, Connection: conn[1], ReadBuff: bufio.NewReader(conn[1]), WriteBuff: bufio.NewWriter(conn[1])},
		},
	}

	
}

// //Lobby creates a lobby
// func (g *Game) Lobby(connchan chan net.Conn){

// }

//Run runs the game
func (g *Game) Run()error{
	log.SetFlags(log.Ltime | log.Lshortfile)
	for _, p := range g.Players{
		defer p.Connection.Close()
		err := p.Send("Starting")
		if err != nil {
			log.Println("Error sending inti data to the client: ", p.Connection.RemoteAddr().String())
		}
		err = p.Send(p.String())
		if err != nil {
			log.Println("Error sending inti data to the client: ", p.Connection.RemoteAddr().String())
		}
		err = p.Send(strconv.Itoa(len(g.Players)))
		if err != nil {
			log.Println("Error sending inti data to the client: ", p.Connection.RemoteAddr().String())
		}
		for _, v := range g.Players{
			if p == v {
				continue
			}
			err = p.Send(v.String())
			if err != nil {
				log.Println("Error sending inti data to the client: ", p.Connection.RemoteAddr().String())
			}
		}
	}
	for {
		err := g.Update()
		if err != nil {
			return err
		}
	}
}
//Update updates the game
func (g *Game) Update()error{
	var wg sync.WaitGroup
	errchan := make(chan error)
	
	// go func() {
	// 	for {
	// 		select{
	// 		case str := <- g.readchan:
	// 			for _, v := range g.Players{
	// 				err := v.Decode(str)
	// 				if err != nil {
	// 					continue
	// 				}
	// 			}
	// 		case err := <-g.errchan:
	// 			log.Println("Error reading info from the clients, ", err.Error())
	// 		}
	// 	}
	// }()
	for _, v := range g.Players{
		wg.Add(1)
		go func(p *Player, players []*Player) {
			defer wg.Done()
			//p.Connection.SetDeadline(time.Now().Add(time.Millisecond *100))
			str, err := p.Read()
			if err != nil {
				log.Println("Error getting info from the client, ", err.Error())
				errchan <- err
				return
			}
			err = p.Decode(str)
			if err != nil {
				log.Println("Error decoding the player, ", err.Error())
				errchan <-err
				return
			}
			for _, v := range players{
				if p == v {
					continue
				}
				err := p.Send(v.String())
				if err != nil {
					log.Println("Error sending info to the client")
					errchan <- err
					return
				}
			}
		}(v, g.Players)
		
	}
	done := make(chan bool)
	var errs []error
	go func() {
		for v := range errchan{
			errs = append(errs, v)
		}
		done <- true
	}() 
	wg.Wait()
	close(errchan)
	<-done
	if len(errs) != 0 {
		for _, v := range errs{
			log.Println(v)
			if neterr, ok := v.(net.Error); ok && neterr.Timeout(){
				continue
			}else if errors.Is(v, errors.New("Wrong Player")){
				continue
			}
			return v
		}
	}
	return nil
}

// //Send sends data to all the clien ts
// func (g *Game)Write(conn *net.TCPConn, done chan bool){
// 	buff := bufio.NewWriter(conn)
// 	for {
// 		str := <- g.writechan
// 		println("Writing to the connection, " + str)
// 		_, err := buff.WriteString(str + "\n")
// 		//println("done writeng")
// 		if err != nil {
// 			g.errchan <- err
// 		}
// 		err = buff.Flush()
// 		if err != nil {
// 			g.errchan <- err
// 		}
// 		done <- true
// 	}
// }

// func (g *Game)Read(conn *net.TCPConn){
// 	buff := bufio.NewReader(conn) 
// 	for {
// 		str, err := buff.ReadString('\n')
// 		if err != nil {
// 			g.errchan <- err
// 		}
// 		g.readchan <- str
// 	}
// }