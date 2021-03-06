package network

import (
	"bufio"
	"log"
	"net"

	"github.com/lehangajanayake/MissionImposible/frontend/models"
)

//Client model
type Client struct {
	Conn                                                                                                *net.TCPConn
	Players                                                                                             map[byte]*Player
	Bullets                                                                                             []Bullet
	InitialRead, InitialWrite, UpdatePlayerCoordsWrite, UpdatePlayerStateWrite, UpdatePlayerFacingWrite chan string //Writeing chans
}

//Connect connects to the server and returns the Client
func Connect(addr string, port string) (*Client, error) {
	TCPaddr, err := net.ResolveTCPAddr("tcp", addr+":"+port)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, TCPaddr)
	if err != nil {
		return nil, err
	}
	return &Client{
		Conn:        conn,
		Players:     make(map[byte]*Player),
		InitialRead: make(chan string), InitialWrite: make(chan string), UpdatePlayerCoordsWrite: make(chan string), UpdatePlayerStateWrite: make(chan string), UpdatePlayerFacingWrite: make(chan string),
	}, nil
}

const (
	newPlayer          = '0'
	updatePlayerCoords = '1'
	updatePlayerState  = '2'
	updatePlayerFacing = '3'
	initial         = '4'
)

//Read reads form the player client and push the data to the rele vant chan of string
func (c *Client) Read() {
	buff := bufio.NewReader(c.Conn)
	for {
		str, err := buff.ReadString('\n')
		if err != nil {
			log.Println("Error receiving data from the client,", c.Conn.RemoteAddr().String(), ": ", err)
			return
		}
		switch rune(str[0]) {
		case initial:
			c.InitialRead <- str[1 : len(str)-1]
		case newPlayer:
			log.Println("newPlayer")
			c.Players[str[1]], err = NewPlayer(str[2 : len(str)-1])
			if err != nil {
				log.Println("Error creating a newPlayer: ", err)
			}

		case updatePlayerCoords:
			err := c.Players[str[1]].Coords.Update(str[2 : len(str)-1])
			if err != nil {
				log.Println("Error updating the player coords, ", err, " ", str[2:len(str)-1])
			}
		case updatePlayerState:
			str = str[1 : len(str)-1] //trim the suffix and the prefix
			err = c.Players[str[0]].UpdatePlayerState(str[1:])
			if err != nil {
				log.Println("Error decdong the player animation")
			}
		case updatePlayerFacing:
			str = str[1 : len(str)-1] //trim the suffix and the prefix
			err = c.Players[str[0]].UpdatePlayerFacingFront(str[1:])
			if err != nil {
				log.Println("Error updating the player facing, ", err)
			}
		default:
			log.Println("Unknown data packet received")
		}
	}
}

func (c *Client) Write() {
	var str string
	for {
		select {
		case str = <-c.UpdatePlayerCoordsWrite:
			str = string(updatePlayerCoords) + str + "\n"
			n, err := c.Conn.Write([]byte(str))
			if err != nil || n != len(str) {
				log.Println("Error sending data to the server", err)
			}
		case str = <-c.UpdatePlayerStateWrite:
			str = string(updatePlayerState) + str + "\n"
			n, err := c.Conn.Write([]byte(str))
			if err != nil || n != len(str) {
				log.Println("Error sending data to the server", err)
			}
		case str = <-c.UpdatePlayerFacingWrite:
			str = string(updatePlayerFacing) + str + "\n"
			n, err := c.Conn.Write([]byte(str))
			if err != nil || n != len(str) {
				log.Println("Error sending data to the server", err)
			}
		}
	}
}

//Run runs the client forever
func (c *Client) Run(p *models.Player) error {
	go c.Read()
	go c.Write()
	log.Println("Waiting for other players to join")
	str := <-c.InitialRead
	result, err := NewPlayer(str)
	if err != nil {
		return err
	}
	p.Coords.X = result.Coords.X
	p.Coords.Y = result.Coords.Y
	p.FacingFront = result.FacingFront
	p.Guard = result.Guard
	p.Gun.Bullets = result.BulletsLeft
	close(c.InitialRead)
	log.Println(p, str)
	return nil
}
