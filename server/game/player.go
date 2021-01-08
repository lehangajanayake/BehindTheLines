package game

import (
	"fmt"
	"net"
	//"errors"
	"strings"
	"strconv"
	"bufio"
)

//Player is the network player
type Player struct{
	ID, X, Y int
	FacingFront, Guard bool
	Connection *net.TCPConn
	ReadBuff *bufio.Reader
	WriteBuff *bufio.Writer
}

func (p *Player) String()string{
	return fmt.Sprintf("%v,%v,%v,%v,%v", p.ID, p.X, p.Y, p.FacingFront, p.Guard)
}

//Decode decodes the player data fomr a string
func (p *Player) Decode(str string) error{
	var err error
	if str == ""{
		return nil
	}
	result := strings.Split(str, ",")

	p.ID, err = strconv.Atoi(result[0])
	if err != nil {
		panic("1")
	}
	// if id != p.ID{
	// 	return errors.New("Wrong Player")
	// }
	p.X, err = strconv.Atoi(result[1])
	if err != nil {
		panic("2")
	}
	p.Y , err = strconv.Atoi(result[2])
	if err != nil {
		panic("3")
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

//Read reads form a given connection
func (p *Player) Read()(string, error){
	str, err := p.ReadBuff.ReadString('\n')
	if err != nil {
		return "", err
	}
	str = strings.TrimSuffix(str, "\n")
	//log.Printf("got data form client , %v \n", p.ID)
	return str, nil

}

//Send writes the player data to the clients
func (p *Player) Send(str string)error{
	_, err := p.WriteBuff.WriteString(str + "\n")
	if err != nil {
		return err
	}
	return  p.WriteBuff.Flush()
}
