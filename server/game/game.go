package game

import (
	"log"
	"net"
	"sync"

	"github.com/lehangajanayake/BehindTheLine/Server/models"
)

//Game model game
type Game struct{
	PlayerNum int
	Players []*models.Player
}




//AddPlayer Adds a Player to the game 
//returns true if enough players are in the game
func (g *Game) AddPlayer(conn *net.TCPConn)bool{
	if len(g.Players) == g.PlayerNum  - 1{
		return true
	}
	p := &models.Player{
		ID: len(g.Players) + 1,
		Coords: &models.Coordinates{
			X: 100,
			Y: 200,
		},
		FacingFront: true,
		Guard: false,
		Conn: conn,
		InitialRead: make(chan string), UpdatePlayerCoordsRead: make(chan string), UpdatePlayerAnimationRead: make(chan string), UpdatePlayerFacingRead: make(chan string),
		InitialWrite: make(chan string), UpdatePlayerCoordsWrite: make(chan string), UpdatePlayerAnimationWrite: make(chan string), UpdatePlayerFacingWrite: make(chan string),
		Errchan: make(chan error),
	}
	g.Players = append(g.Players, p)
	go p.Read()
	return false

}

//Run starts the game and runs the game returns an error if the game fails to run
func (g *Game) Run(){
	var wg sync.WaitGroup
	for {
		for _, v := range g.Players {
			wg.Add(1)
			go func(v *models.Player, wg *sync.WaitGroup) {
				defer wg.Done()
				var str string
				var err error
				select{
				case str = <- v.UpdatePlayerCoordsRead:
					err = v.Coords.Update(str)
					if err != nil {
						log.Fatal("Error decoding the Player Coords: ", err)
					}
					for _, p := range g.Players{
						p.UpdatePlayerCoordsWrite <- v.Coords.String()
					}
				case err = <- v.Errchan:
					log.Println("Error getting data: ", err)
				}
			}(v, &wg)
			wg.Wait()
		}
	}
	
}