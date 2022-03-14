package model

import "github.com/labstack/echo"

// EchoGroup to store routes group
type EchoGroup struct {
	Admin *echo.Group
	API   *echo.Group
	Token *echo.Group
}
