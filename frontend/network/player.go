package network

import (
	//"errors"
	"fmt"
	//"log"
	"strconv"
	"strings"
	//"sync"
)

//Player is the network equivalent of the games player
type Player struct {
	ID, X, Y int
	Pos Coordinates
	FacingFront, Guard bool
}

//String returns the string value containing player data
func (p *Player) String()string{
	return fmt.Sprintf("%v,%v,%v,%v,%v", p.ID, p.X, p.Y, p.FacingFront, p.Guard)
}

//Decode decodes the player data from a string
func (p *Player) Decode(str string) error{
	var err error
	//log.Println(str)
	result := strings.Split(str, ",")
	p.ID, err = strconv.Atoi(result[0])
	if err != nil {
		return err
	}
	p.X, err = strconv.Atoi(result[1])
	if err != nil {
		return err
	}
	p.Y , err = strconv.Atoi(result[2])
	if err != nil {
		return err
	}
	if result[3] == "true"{
		p.FacingFront = true
	} 
	p.FacingFront = false
	if result[4] == "true"{
		p.Guard = true
	}
	p.Guard = false
	return nil
}
