package network

import (
	"bufio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	//"strconv"
	"sync"
)

//Client has the data of the client
type Client struct {
        conn *net.TCPConn
        errchan chan error
        readchan, writechan chan string
        readbuff *bufio.Reader
        writebuff *bufio.Writer
}

const (
	non = byte(0)
	updatePlayerPos = byte(1)
	updatePlayerAnimation = byte(2)
	newBullet = byte(3)
	updateBulletPos = byte(4)
)

//Connect connects to the server
func (c *Client) Connect() {
        log.SetFlags(log.Ltime | log.Lshortfile)
        var err error
        addr, err := net.ResolveTCPAddr("tcp", ":8080")
        c.conn, err = net.DialTCP("tcp", nil, addr)
        if err != nil {
                log.Fatal("Error connecting to the server")
        }
        c.errchan = make(chan error)
        c.readchan = make(chan string)
        c.writechan = make(chan string)
        c.readbuff = bufio.NewReader(c.conn)
        c.writebuff =  bufio.NewWriter(c.conn)
}

func (c *Client) read()(string, error){
        str, err := c.readbuff.ReadString('\n')
        str = strings.TrimSuffix(str, "\n")
        if err != nil {
                return "", err
        }
        return str, nil

}

func (c *Client) send(packetID byte, str string)error{
        _, err := c.writebuff.WriteString(string(packetID) + str + "\n")
        if err != nil {
                return err
        }
        return c.writebuff.Flush()
}



//CheckLobby checks of the player is in the game
func checkLobby(str string)bool{
	log.Printf("Waiting in the Lobby \r")
	if str == "Lobby"{
		return true
	}
	return false
}

//Start starts the client
func (c *Client) Start(player *Player, players *[]*Player){
        // done := make(chan bool)
        // defer func ()  {
        //         log.Println("Done")
        //         done <- true
        // }()
        // go c.read()
        // go c.send()
        // go func() {
        //         for {
        //                select {
        //                case err := <- c.errchan:
        //                        log.Println("Error: ", err.Error())
        //                case <-done:
        //                 break
        //                }
                        
        //         }
        // }()
        // for checkLobby(<-c.readchan){
        //         log.Println("IN the lobby")
        // }

        // err := player.Decode(<-c.readchan)
        // if err != nil {
        //         log.Fatal("Error decoding the player: ", err.Error())
        // }
        // len, err := strconv.Atoi(<-c.readchan)
        // if err != nil {
        //         log.Fatal("Error getting the player count: ", err.Error())
        // }
        // for i := 0; i < len - 1; i++{
        //         p := new(Player)
        //         err := p.Decode(<-c.readchan)
        //         if err != nil {
        //                 log.Fatal("Error getting the player data: ", err.Error())
        //         }
        //         *players = append(*players, p)
       // }
     str, err := c.read()
     if err != nil {
        log.Fatal("Error getting information form the server, ", err)
     }
     if checkLobby(str){
        for{
                str, err = c.read()
                if err != nil {
                        log.Fatal("Error getting information form the server, ", err)
                }
                if checkLobby(str){
                        continue
                }
                if str[1:] == "Starting"{
                        log.Println("starting the game")
                        break
                }
        }
     }
     str, err = c.read()
     if err != nil {
        log.Fatal("Error starting the game, ", err.Error())
     }
     err = player.Decode(str[1:])
     if err != nil {
             log.Fatal("Error decodeing the player, ", err.Error())
     }
     var num int
     str, err = c.read()
     if err != nil {
        log.Fatal("Error getting information from  the serve, ", err.Error())
     }
     num, err = strconv.Atoi(str[1:])
     if err != nil {
             log.Fatal("Error getting the number of players, ", err.Error())
     }
     log.Println("Got the number of players: ", num)
     for i := 0; i < num -1 ; i++{
        println(i)
        str, err = c.read()
        println("read")
        if err != nil {
                log.Fatal("Error getting information form the server, ", err.Error())
        }
        p := new(Player)
        err = p.Decode(str[1:])
        if err != nil{
                log.Fatal("Error decodeing the player, ", err.Error())
        }
        *players = append(*players, p)
     }
     //str, err = c.read()
//      if err != nil {
//              log.Fatal("The game didnt start")
//      }
//      log.Println(str)
     return
}

//SendAndGet send the player to server 
func (c *Client) SendAndGet(player *Player, players []*Player, wg *sync.WaitGroup){
        defer wg.Done()
        err := c.send(updatePlayerPos, player.Pos.String())
        if err !=  nil {
                log.Println("Error sending info to the server,", err.Error())
                return  
        }
        for _, v := range players{
                c.conn.SetReadDeadline(time.Now().Add(time.Second *1))
                str, err := c.read()
                if err != nil{
                        log.Println("Error getting the players, ", err.Error())
                        return 
                }
                switch str[0]{
                case updatePlayerPos:
                        if str[2] != byte(v.ID){
                                continue
                        }
                        str = strings.Trim(str, str[:2])
                        err = v.Pos.Update(str)
                        if err != nil {
                                log.Println("Error decoding the player, ", err)
                                return
                        }
                case non:
                        str = strings.TrimPrefix(str, string(non)) 
                        err = v.Decode(str)
                        if err != nil {
                               log.Println("Error decoding the player", err.Error())
                               return
                       }
                }
               
                
        }
}