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
		NewPlayerRead: make(chan string), UpdatePlayerCoordsRead: make(chan string), UpdatePlayerAnimationRead: make(chan string), UpdatePlayerFacingRead: make(chan string),
		NewPlayerWrite: make(chan string), UpdatePlayerCoordsWrite: make(chan string), UpdatePlayerAnimationWrite: make(chan string), UpdatePlayerFacingWrite: make(chan string),
		Errchan: make(chan error),
	}
	go p.Read()
	for _, otherP := range g.Players{
		println("Sending new player")
		otherP.NewPlayerWrite <- p.String()
	}
	g.Players = append(g.Players, p)
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
						log.Println("Sending the coord to the other players")
						p.UpdatePlayerCoordsWrite <- v.Coords.String()
					}
				case err = <- v.Errchan:
					log.Println("Error getting data: ", err)
				default:	
					return
				}
			}(v, &wg)
			wg.Wait()
		}
	}
	
}