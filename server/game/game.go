package game

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/lehangajanayake/BehindTheLine/Server/models"
)

//Game model game
type Game struct {
	PlayerNum int
	//Players   []*models.Player
	Players   map[byte]*models.Player
}

var PlayerArray [3]struct {
	ID, X, Y, Bullets int
	Guard, Facing     bool
} = [3]struct {
	ID, X, Y, Bullets int
	Guard, Facing     bool
}{
	struct {
		ID      int
		X       int
		Y       int
		Bullets int
		Guard   bool
		Facing  bool
	}{0, 100, 200, 60, true, true},
	struct {
		ID      int
		X       int
		Y       int
		Bullets int
		Guard   bool
		Facing  bool
	}{1, 200, 100, 60, false, false},
	struct {
		ID      int
		X       int
		Y       int
		Bullets int
		Guard   bool
		Facing  bool
	}{2, 100, 200, 60, false, true},
}

func NewGame(pNum int) *Game {
	return &Game{
		PlayerNum: pNum,
		Players:   make(map[byte]*models.Player, pNum),
	}

}

//AddPlayer Adds a Player to the game
//returns true if enough players are in the game
func (g *Game) AddPlayer(conn *net.TCPConn) bool {
	if len(g.Players) == g.PlayerNum {
		return true
	}
	pl := PlayerArray[len(g.Players)]
	g.Players[byte(pl.ID)] = models.NewPlayer(pl.ID, pl.X, pl.Y, pl.Bullets, pl.Guard, pl.Facing, conn)
	log.Println(g.Players)
	go g.Players[byte(pl.ID)].Read()
	go g.Players[byte(pl.ID)].Write()
	// go func() {
	// 	for {
	// 		err := <-p.Errchan
	// 		log.Println("Error in Client, ", p.Conn.RemoteAddr().String(), ":", err)
	// 		//g.Players = append(g.Players[:pl.ID], g.Players[pl.ID + 1:]...)
	// 	}
	// }()
	return len(g.Players) == g.PlayerNum

}

//Run starts the game and runs the game returns an error if the game fails to run
func (g *Game) Run() {
	tick := time.NewTicker(time.Millisecond*16)
	defer tick.Stop()
	var wg sync.WaitGroup
	done := make(chan bool)
	for k, p := range g.Players {
		p.InitialWrite <- p.String()
		for otherK, otherP := range g.Players {
			if k == otherK {
				continue
			}
			log.Println("sent player info to", otherP.ID)
			otherP.NewPlayerWrite <- p.String()
		}
	}
	for {
		// if len(g.Players) >= 1 {
		// 	log.Println("Not Enough Players to play")
		// 	return 
		// }
		for k, v := range g.Players {
			wg.Add(1)
			go func(v *models.Player, wg *sync.WaitGroup) {
				defer wg.Done()
				var str string
				var err error
				select {
				case str = <- v.UpdatePlayerCoordsRead:
					err = v.Coords.Update(str)
					if err != nil {
						log.Println("Error decoding the Player Coords: ", err)
						return
					}
					for _, p := range g.Players {
						if v.ID == p.ID {
							continue
						}
						p.UpdatePlayerCoordsWrite <- v.Coords.String()
					}
				case str = <-v.UpdatePlayerStateRead:
					v.Animation = str
					for _, p := range g.Players {
						if v.ID == p.ID {
							continue
						}
						p.UpdatePlayerStateWrite <- str
					}
				case str = <-v.UpdatePlayerFacingRead:
					switch str {
					case "true":
						v.FacingFront = true
					case "false":
						v.FacingFront = false
					}
					for _, p := range g.Players {
						if v.ID == p.ID {
							continue
						}
						p.UpdatePlayerFacingWrite <- str
					}

				case err = <-v.Errchan:
					log.Println("Error getting data: ", err)
					delete(g.Players, k)
					done <- true
					return
				// case <- tick.C:
				// 	log.Println("ticked")
				// 	return
				default:
					return
				}
			}(v, &wg)
			wg.Wait()
			select {
			case <-done:
				for _, p := range g.Players {
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
