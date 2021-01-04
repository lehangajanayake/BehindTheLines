package network

import (
	"log"
	"net"
	"strconv"
)

//Global connection to the server
var conn net.Conn

//Connect connects to the server
func Connect() {
        var err error
        conn, err = net.Dial("tcp", "localhost:8080")
        if err != nil {
                log.Fatal("Error connecting to the server")
        }
}

//Start starts the client
func Start(player *Player, players []*Player){
        readstrchan := make(chan string)
        sendstrchan := make(chan string)
        errchan := make(chan error)
        go read(readstrchan, errchan)
        go send(sendstrchan, errchan)
        first := true // is the handshake done
        for first{
                select{
                case str := <- readstrchan:
                        var err error
                        if checkLobby(str){
                                log.Println("Waiting in the Lobby")
                                continue
                        }
                        player.ID, err = strconv.Atoi(str)
                        if err != nil {
                                log.Fatalf("Error getting the player IDk, %s", err.Error())
                                break
                        }
                        first = false
                      
                case err := <- errchan:
                        log.Fatal(err)
                }
                //sendstrchan <- player.String()
        }
        for{
                log.Println(player.String())
                sendstrchan <- player.String()
                print("help")
                select{
                case str := <- readstrchan:
                        println("read info")
                        for _, v := range players{
                                err := v.Decode(str)
                                if err != nil {
                                        if err.Error() == "Wrong Player"{
                                                continue
                                        }
                                        log.Println("Error updating the Player")
                                }
                }
                case err := <- errchan:
                        log.Fatalf("Error getting informaition from the server, %v", err)
                }
               
                
        }
}