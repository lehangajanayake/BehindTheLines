package models

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

//Player model
type Player struct {
	ID, BulletsLeft                                                                              int
	FacingFront, Guard                                                                           bool
	Animation                                                                                    string
	Coords                                                                                       *Coordinates
	Conn                                                                                         *net.TCPConn
	NewPlayerRead, UpdatePlayerCoordsRead, UpdatePlayerStateRead, UpdatePlayerFacingRead     chan string // Readign chans
	InitialWrite, NewPlayerWrite, UpdatePlayerCoordsWrite, UpdatePlayerStateWrite, UpdatePlayerFacingWrite chan string //Writeing chans
	doneRead, doneWrite                                                                          chan bool
	Errchan                                                                                      chan error
}

const (
	newPlayer             = '0'
	updatePlayerCoords    = '1'
	updatePlayerState = '2'
	updatePlayerFacing    = '3'
	initial = '4'
)

//NewPlayer creates a new player
func NewPlayer(id, x, y, bullets int, guard, facing bool, conn *net.TCPConn) *Player {
	return &Player{
		ID:                     id,
		Conn:                   conn,
		Coords:                 &Coordinates{X: x, Y: y},
		BulletsLeft:            bullets,
		Guard:                  guard,
		FacingFront:            facing,
		Animation: "Idle",
		UpdatePlayerCoordsRead: make(chan string), UpdatePlayerStateRead: make(chan string), UpdatePlayerFacingRead: make(chan string),
		InitialWrite: make(chan string), NewPlayerWrite: make(chan string), UpdatePlayerCoordsWrite: make(chan string), UpdatePlayerStateWrite: make(chan string), UpdatePlayerFacingWrite: make(chan string),
		Errchan:  make(chan error),
		doneRead: make(chan bool), doneWrite: make(chan bool),
	}
}

func (p *Player) String() string {
	return fmt.Sprintf("%d,%d,%t,%t,%s,%d", p.Coords.X, p.Coords.Y, p.FacingFront, p.Guard, p.Animation, p.BulletsLeft)
}

//Close closes all the channels and the goroutines releted to the player
func (p *Player) Close() {
	p.doneRead <- true
	p.doneWrite <- true
	close(p.doneRead)
	close(p.doneWrite)
	p.Conn.Close()
}

//Read reads form the player client and push the data to the relavant chan of string
func (p *Player) Read() {
	buff := bufio.NewReader(p.Conn)
	for {
		str, err := buff.ReadString('\n')
		if err != nil {
			log.Println("Error receiving data from the client,", p.Conn.RemoteAddr().String(), ": ", err)
			p.Errchan <- err
			return
		}
		switch rune(str[0]) {
		case newPlayer:
			log.Println("Sending a new player")
			str = str[1 : len(str)-1] //trim the suffix and the prefix
			p.NewPlayerRead <- str
		case updatePlayerCoords:
			str = str[1 : len(str)-1] //trim the suffix and the prefix
			p.UpdatePlayerCoordsRead <- str
		case updatePlayerState:
			str = str[1 : len(str)-1] //trim the suffix and the prefix
			p.UpdatePlayerStateRead <- str
		case updatePlayerFacing:
			str = str[1 : len(str)-1] //trim the suffix and the prefix
			p.UpdatePlayerFacingRead <- str
		}
		select {
		case <-p.doneRead:
			close(p.NewPlayerRead)
			close(p.UpdatePlayerFacingRead)
			close(p.UpdatePlayerCoordsRead)
			close(p.UpdatePlayerStateRead)
			return
		default:
			continue
		}
	}
}

func (p *Player) Write() {
	var str string
	dataLostErr := errors.New("Data lost writing to the client: " + p.Conn.RemoteAddr().String())
	for {
		//p.Conn.SetWriteDeadline(time.Now().Add(time.Second *1))
		select {
		case str = <- p.InitialWrite:
			data := string(initial) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {	
				p.Errchan <- err
			} else if n != len(data) {
				p.Errchan <- dataLostErr
			}
		case str = <-p.NewPlayerWrite:
			data := string(newPlayer) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {
				p.Errchan <- err
			} else if n != len(data) {
				p.Errchan <- dataLostErr
			}

		case str = <-p.UpdatePlayerCoordsWrite:
			data := string(updatePlayerCoords) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {
				p.Errchan <- err
			} else if n != len(data) {
				p.Errchan <- dataLostErr
			}
		case str = <-p.UpdatePlayerStateWrite:
			data := string(updatePlayerState) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {
				p.Errchan <- err
			} else if n != len(data) {
				p.Errchan <- dataLostErr
			}
		case str = <-p.UpdatePlayerFacingWrite:
			data := string(updatePlayerFacing) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {
				p.Errchan <- err
			} else if n != len(data) {
				p.Errchan <- dataLostErr
			}
		case <-p.doneWrite:
			close(p.NewPlayerWrite)
			close(p.UpdatePlayerFacingWrite)
			close(p.UpdatePlayerCoordsWrite)
			close(p.UpdatePlayerStateWrite)
			close(p.InitialWrite)
			return
		}

	}
}
