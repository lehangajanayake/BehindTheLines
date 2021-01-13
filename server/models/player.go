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
type Player struct{
	ID, BulletsLeft int
	FacingFront, Guard bool
	Animation string
	Coords *Coordinates
	Conn *net.TCPConn
	NewPlayerRead, UpdatePlayerCoordsRead, UpdatePlayerAnimationRead, UpdatePlayerFacingRead  chan string // Readign chans
	NewPlayerWrite, UpdatePlayerCoordsWrite, UpdatePlayerAnimationWrite, UpdatePlayerFacingWrite  chan string //Writeing chans
	Errchan chan error
}

const(
	newPlayer = '0'
	updatePlayerCoords = '1'
	updatePlayerAnimation = '2'
	updatePlayerFacing = '3'
)

func (p *Player) String() string{
	return fmt.Sprintf("%d,%d,%t,%t,%s,%d", p.Coords.X, p.Coords.Y, p.FacingFront, p.Guard, p.Animation, p.BulletsLeft)
}


//Read reads form the player client and push the data to the relavant chan of string
func (p *Player) Read(){
	buff := bufio.NewReader(p.Conn)
	for {
		str, err := buff.ReadString('\n')
		if err != nil {
			log.Println("Error receiving data from the client,", p.Conn.RemoteAddr().String(), ": ", err)
			p.Errchan <- err
			return
		}
		log.Println(str)
		switch rune(str[0]){
		case newPlayer:
			str = str[1:len(str)-1] //trim the suffix and the prefix
			p.NewPlayerRead <- str
		case updatePlayerCoords:
			str = str[1:len(str)-1] //trim the suffix and the prefix
			p.UpdatePlayerCoordsRead <- str
		case updatePlayerAnimation:
			str = str[1:len(str)-1] //trim the suffix and the prefix
			p.UpdatePlayerAnimationRead <- str
		case updatePlayerFacing:
			str = str[1:len(str)-1] //trim the suffix and the prefix
			p.UpdatePlayerFacingRead <- str
		}
	}
}

func (p *Player) Write(){
	var str string
	dataLostErr := errors.New("Data lost writing to the client: " + p.Conn.RemoteAddr().String())
	for {
		select{
		case str = <- p.NewPlayerRead:
			data := string(newPlayer) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {
				p.Errchan <- err 
			}else if n != len(data){
				p.Errchan <- dataLostErr
			}			

		case str = <- p.UpdatePlayerCoordsWrite:
			data := string(updatePlayerCoords) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil {
				p.Errchan <- err
			}else if n != len(data){
				p.Errchan <- dataLostErr
			}
		case str = <- p.UpdatePlayerAnimationWrite:
			data := string(updatePlayerAnimation) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil{
				p.Errchan <- err
			}else if n != len(data){
				p.Errchan <- dataLostErr
			}
		case str = <- p.UpdatePlayerFacingWrite:
			data := string(updatePlayerFacing) + strconv.Itoa(p.ID) + str + "\n"
			n, err := p.Conn.Write([]byte(data))
			if err != nil  {
				p.Errchan <- err
			}else if n != len(data){
				p.Errchan <- dataLostErr
			}
		}
	}
}