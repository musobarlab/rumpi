package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/musobarlab/rumpi/config"
	"github.com/musobarlab/rumpi/pkg/chathub"
)

func indexHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(res, req)
		return
	}
	fmt.Fprintf(res, "hell yeah!")
}

func main() {

	err := config.LoadEnv()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("starting application...")
	manager := chathub.NewManager(config.Env.WebsocketKey)

	// wsHandler := chathub.Handler{Manager: manager}

	// router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/", indexHandler).Methods("GET")
	// router.HandleFunc("/ws", wsHandler.WsHandler)

	// // start client manager
	// go manager.Run()

	// fmt.Println("server running on port 9000")
	// log.Fatal(http.ListenAndServe(":9000", router))

	//---------------------------- echo ---------------
	wsHandler := chathub.EchoHandler{Manager: manager}

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

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Env.HTTPPort)))
}
