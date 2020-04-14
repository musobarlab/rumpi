package core

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

// EchoHandler struct
type EchoHandler struct {
	Manager *Manager
}

func (h *EchoHandler) WsHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		sock, err := h.Manager.Upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return echo.NewHTTPError(400, err.Error())
		}

		id := c.Request().Header.Get("Sec-Websocket-Key")
		fmt.Println(id)

		var client Client
		client.ID = id
		client.Conn = sock
		client.MsgChan = make(chan []byte)
		client.Room = make(map[string]bool)
		client.Manager = h.Manager

		h.Manager.Register <- &client

		// Read message
		go client.Read()

		// Write message
		go client.Write()

		return nil
	}
}
