package game

import (
	"log"
	"net"
	"strconv"
	"sync"
)

//Game is the main game
type Game struct{
	Players [2]*Player
}

//NewGame created a new game
func NewGame(conn []net.Conn)*Game{
	return &Game{
		[2]*Player{
			{ID: 1, X: 400, Y:400, FacingFront: true, Guard:true, Connection: conn[0]},
			{ID: 2,X:200, Y:400, FacingFront: true, Guard:false, Connection: conn[1]},
		},
	}

	
}

// //Lobby creates a lobby
// func (g *Game) Lobby(connchan chan net.Conn){

// }

//Run runs the game
func (g *Game) Run()error{
	for _, v := range g.Players{
		v.Send(strconv.Itoa(v.ID))
	}
	for {
		err := g.Update()
		if err != nil{
			return err
		}
	}
}
//Update updates the game
func (g *Game) Update()error{
	var wg sync.WaitGroup
	errchan := make(chan error, 2)
	for i, v := range g.Players{
		wg.Add(1)
		go func(i int, p *Player, wg *sync.WaitGroup ) {
			defer wg.Done()
			str, err := p.Read()
			if err != nil {
				return
			}
			err = p.Decode(str)
			if err != nil {
				log.Println("Error decoding the string, %w", err)
				return
			}
			for i2, v2 := range g.Players{ //Send the players information to each player connected
				if i == i2 {
					continue
				}
				p.Send(v2.String())
				
			}
		}(i, v, &wg)
	}
	wg.Wait()
	if err := <-errchan; err != nil { //Read twice from the channel
		return err
	}
	if err := <-errchan; err != nil {
		return err
	}
	return nil
	
}