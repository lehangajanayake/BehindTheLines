package network

import (
        "net"
        "log"
        "bufio"
)

//Client model
type Client struct{
        Conn *net.TCPConn
        Players map[byte]*Player
        Bullets []Bullet
	InitialWrite, UpdatePlayerCoordsWrite, UpdatePlayerAnimationWrite, UpdatePlayerFacingWrite  chan string //Writeing chans
}

//Connect connects to the server and returns the Client
func Connect(addr string, port string)(*Client, error){
        TCPaddr, err := net.ResolveTCPAddr("tcp", addr + ":" + port)
        if err != nil{
                return nil, err
        }
        conn, err := net.DialTCP("tcp", nil, TCPaddr)
        if err != nil {
                return nil, err
        }
        return &Client{
                Conn: conn,
                Players: make(map[byte]*Player, 1),
                InitialWrite: make(chan string), UpdatePlayerCoordsWrite: make(chan string), UpdatePlayerAnimationWrite: make(chan string), UpdatePlayerFacingWrite: make(chan string),
        }, nil
}

const(
	initial = '0'
	updatePlayerCoords = '1'
	updatePlayerAnimation = '2'
	updatePlayerFacing = '3'
)


//Read reads form the player client and push the data to the rele vant chan of string
func (c *Client) Read(){
        buff := bufio.NewReader(c.Conn)
	for {
		str, err := buff.ReadString('\n')
		if err != nil {
			log.Println("Error receiving data from the client,", c.Conn.RemoteAddr().String(), ": ", err)
			return
		}
		switch rune(str[0]){
		case updatePlayerCoords:
                        str = str[1:len(str)-1] //trim the suffix and the prefix
                        err := c.Players[str[0]].Coords.Update(str[1:])
                        if err != nil {
                                log.Println("Error updating the player coords, ", err)
                        }
		case updatePlayerAnimation:
			str = str[1:len(str)-1] //trim the suffix and the prefix
			c.Players[str[0]].Animation = str[1:]
		case updatePlayerFacing:
			str = str[1:len(str)-1] //trim the suffix and the prefix
                        err = c.Players[str[0]].UpdatePlayerFacingFront(str[1:])
                        if err != nil {
                                log.Println("Error updating the player facing, ", err)
                        }
		}
	}
}

func (c *Client) Write(){
        var str string
	for {
		select{
                case str = <- c.UpdatePlayerCoordsWrite:
                        println("got data")
			data := string(updatePlayerCoords) + str
			n, err := c.Conn.Write([]byte(data))
			if err != nil || n != len(data) {
				return 
			}
		}
	}
}

//Run runs the client forever
func (c *Client) Run(){
        go c.Read()
        go c.Write()
}