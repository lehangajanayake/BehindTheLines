package network

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

//Coordinates hold the x and y value
type Coordinates struct{
	X, Y int
}

//String retruns a string with x and y values
func (c *Coordinates) String()string{
	return fmt.Sprintf("%d,%d", c.X, c.Y)
}

//Update updates the coordinates using a string
func (c *Coordinates) Update(str string)error{
	var err error
	str = str[0:]
	strslice := strings.Split(str, ",")
	if len(strslice) != 2 {
		return errors.New("Invalid var length")
	}
	c.X, err = strconv.Atoi(strslice[0])
	if err != nil {
		return err
	}
	c.Y, err = strconv.Atoi(strslice[1])
	if err != nil {
		return err
	}

	return nil
}