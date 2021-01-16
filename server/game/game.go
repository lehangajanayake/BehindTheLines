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
	if len(g.Players) == g.PlayerNum  {
		return true
	}
	p := models.NewPlayer(len(g.Players), 100, 200, 60 , true, true, conn)
	g.Players = append(g.Players, p)
	log.Println(g.Players)
	go p.Read()
	go p.Write()
	if len(g.Players) == g.PlayerNum {
		return true
	} 
	return false

}

//Run starts the game and runs the game returns an error if the game fails to run
func (g *Game) Run(){
	var wg sync.WaitGroup
	done := make(chan bool)
	for _, p := range g.Players{
		for _, otherP := range g.Players{
			if p.ID == otherP.ID{
				continue
			}
			println("hello")
			otherP.NewPlayerWrite <- p.String()
		}
	}
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
						log.Println("Error decoding the Player Coords: ", err)
						return
					}
					for _, p := range g.Players{
						if v.ID == p.ID {
							continue
						}
						p.UpdatePlayerCoordsWrite <- v.Coords.String()
					}
				case str = <- v.UpdatePlayerAnimationRead:
					v.Animation = str
					for _, p := range g.Players{
						if v.ID == p.ID {
							continue
						}
						p.UpdatePlayerAnimationWrite <- str
					}

				case err = <- v.Errchan:
					log.Println("Error getting data: ", err)
					done <- true
					return
				default:	
					return
				}
			}(v, &wg)
			wg.Wait()
			select{
			case <-done:
				for _, p := range g.Players{
					p.Close()
				}
				close(done)
				break
			default:
				continue
			}
		}
	}
	
}