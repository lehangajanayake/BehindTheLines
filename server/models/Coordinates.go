package models

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//Coordinates hold a value of x, y of teh player
type Coordinates struct{
	X, Y int
}


//Update decodes the string and updates the coords
func (c *Coordinates)Update(str string)error{
	var err error
	fmt.Println("Read Coords, ", str)
	result := strings.Split(str, ",")
	if len(result) != 2 {
		return errors.New("Invalid string")
	}
	c.X, err = strconv.Atoi(result[0])
	if err != nil {
		return nil
	}
	c.Y, err = strconv.Atoi(result[1])
	if err != nil {
		return err
	}
	return nil
}

//String returns the string value of coords
func (c *Coordinates)String()string{
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}