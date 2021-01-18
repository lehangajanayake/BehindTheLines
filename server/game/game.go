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
	Players   []*models.Player
}

var PlayerArray [2]struct {
	ID, X, Y, Bullets int
	Guard, Facing     bool
} = [2]struct {
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
}

func NewGame(pNum int) *Game {
	return &Game{
		PlayerNum: pNum,
		Players:   make([]*models.Player, 0),
	}

}

//AddPlayer Adds a Player to the game
//returns true if enough players are in the game
func (g *Game) AddPlayer(conn *net.TCPConn) bool {
	if len(g.Players) == g.PlayerNum {
		return true
	}
	pl := PlayerArray[len(g.Players)]
	p := models.NewPlayer(pl.ID, pl.X, pl.Y, pl.Bullets, pl.Guard, pl.Facing, conn)
	g.Players = append(g.Players, p)
	log.Println(g.Players)
	go p.Read()
	go p.Write()
	// go func() {
	// 	for {
	// 		err := <-p.Errchan
	// 		log.Println("Error in Client, ", p.Conn.RemoteAddr().String(), ":", err)
	// 		//g.Players = append(g.Players[:pl.ID], g.Players[pl.ID + 1:]...)
	// 	}
	// }()
	m := make(map[string]string, 2)
	println(m)
	return len(g.Players) == g.PlayerNum

}

//Run starts the game and runs the game returns an error if the game fails to run
func (g *Game) Run() {
	tick := time.NewTicker(time.Millisecond*16)
	defer tick.Stop()
	var wg sync.WaitGroup
	done := make(chan bool)
	for _, p := range g.Players {
		p.InitialWrite <- p.String()
		for _, otherP := range g.Players {
			if p.ID == otherP.ID {
				continue
			}
			otherP.NewPlayerWrite <- p.String()
		}
	}
	for {
		if len(g.Players) >= 1 {
			return 
		}
		for _, v := range g.Players {
			wg.Add(1)
			go func(v *models.Player, wg *sync.WaitGroup) {
				defer wg.Done()
				var str string
				var err error
				select {
				case str = <-v.UpdatePlayerCoordsRead:
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
					done <- true
					return
				case <- tick.C:
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
