package utils

import "fmt"

type RGB struct {
	R byte
	G byte
	B byte
}

func (c *RGB) String() string {
	return fmt.Sprintf("(%d,%d,%d)", c.R, c.G, c.B)
}

func (c *RGB) Hex() string {
	return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
} 
