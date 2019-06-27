package api

import (
	"fmt"

	db "../db"
	router "../router"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

type RouterCfg struct {
	Port string
	Host string
}

type Api struct {
	Router       *gin.Engine
	DBConfig     pgx.ConnConfig
	RouterConfig RouterCfg
}

func Initialize() (*Api, error) {

	var (
		api Api
		err error
	)

	err = db.Initialize("postgres://docker:1337@localhost:5432/forum")
	// err = db.Initialize("postgres://postgres:1337@localhost/forum")

	fmt.Println("------------------------------------------------\nDatabase initialize: OK")

	api.RouterConfig = RouterCfg{
		Host: "localhost",
		Port: "5000",
	}

	api.Router = router.Initialize()

	fmt.Println("Router initialize  : OK")

	if err != nil {
		return nil, err
	}

	return &api, nil
}

func (api *Api) Run() (err error) {
	gin.SetMode(gin.DebugMode)
	err = api.Router.Run(":5000")
	return
}
