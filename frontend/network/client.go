package network

import (
	"bufio"
)


func read(strchan chan string, errchan chan error){
	for {
		str, err := bufio.NewReader(conn).ReadString('\n')
		str = str[:len(str)-1]
		if err != nil {
			errchan <- err
			return
		}
		strchan <- str
	}

}

func send(strchan chan string, errchan chan error){
	for {
		str := <- strchan
		println("sending data", str)
		_, err := conn.Write([]byte(str))
		println("sent data")
		if err != nil {
			errchan <- err
		}
	}
}



//CheckLobby checks of teh player is in the game
func checkLobby(str string)bool{
	if str == "Lobby"{
		return true
	}
	return false
}