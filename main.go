package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/musobarlab/go-websocket-chat/core"
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprintf(res, "hell yeah!")
}

func main() {
	fmt.Println("starting application...")
	manager := core.NewManager("555abcd")

	// wsHandler := core.Handler{Manager: manager}

	// router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/", indexHandler).Methods("GET")
	// router.HandleFunc("/ws", wsHandler.WsHandler)

	// // start client manager
	// go manager.Run()

	// fmt.Println("server running on port 9000")
	// log.Fatal(http.ListenAndServe(":9000", router))

	//---------------------------- echo ---------------
	wsHandler := core.EchoHandler{Manager: manager}

	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "hell yeah")
	})

	e.POST("/join", func(c echo.Context) error {
		return c.String(http.StatusOK, "success")
	})

	e.GET("/ws", wsHandler.WsHandler())

	go manager.Run()

	e.Logger.Fatal(e.Start(":9000"))
}
