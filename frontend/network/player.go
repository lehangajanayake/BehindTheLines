package network

import (
	"errors"
	"net"
	"strconv"
	"strings"
	"sync"
)

//Player model
type Player struct {
	mutex              sync.RWMutex
	BulletsLeft        int
	CurrentFrame       int
	FacingFront, Guard bool
	State          string
	Coords             *Coordinates
	Conn               *net.TCPConn
}

//NewPlayer decodes a the string and  a new player
func NewPlayer(str string) (*Player, error) {
	var err error
	p := &Player{mutex: sync.RWMutex{}, Coords: &Coordinates{mutex: sync.RWMutex{}}}
	result := strings.Split(str, ",")
	if len(result) != 6 {
		return nil, errors.New("Invalid string : " + str)
	}
	p.Coords.X, err = strconv.Atoi(result[0])
	if err != nil {
		return nil, err
	}
	p.Coords.Y, err = strconv.Atoi(result[1])
	if err != nil {
		return nil, err
	}
	if result[2] == "true" {
		p.FacingFront = true
	} else {
		p.FacingFront = false
	}
	if result[3] == "true" {
		p.Guard = true
	} else {
		p.Guard = false
	}
	p.State = result[4]
	p.BulletsLeft, err = strconv.Atoi(result[5])
	if err != nil {
		return nil, err
	}
	return p, nil

}

//UpdatePlayerFacingFront updates the facing of the player using a str
func (p *Player) UpdatePlayerFacingFront(str string) error {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return err
	}
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	if b {
		p.FacingFront = true
	} else {
		p.FacingFront = false
	}
	return nil
}

//UpdatePlayerState updates the animation type using a string
func (p *Player) UpdatePlayerState(str string) error {
	switch str {
	case "Idle", "Walking", "Shooting":
		p.mutex.RLock()
		defer p.mutex.RUnlock()
		p.State = str
		return nil
	default:
		return errors.New("Invalid String: " + str)
	}
}
