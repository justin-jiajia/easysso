package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/justin-jiajia/easysso/api/config"
	"github.com/justin-jiajia/easysso/api/database"
	"github.com/justin-jiajia/easysso/api/router"

	"github.com/gin-gonic/gin"
)

//go:embed all:front_output/*
var front_fs_orgi embed.FS

func main() {
	log.SetFlags(log.Llongfile | log.Lmicroseconds | log.Ldate)
	log.Println("Reading Config...")
	config.ReadConfig()
	log.Println("Creating web server...")
	r := gin.New()
	if !config.Config.IsDev{
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	log.Println("Connecting to database...")
	database.InitDB()
	log.Println("Mirgrating database...")
	database.Migrate()
	log.Println("Loading routes")
	front_fs, _ := fs.Sub(front_fs_orgi, "front_output")
	router.InitApi(r)
	r.StaticFS("/ui", http.FS(front_fs))
	log.Println("Running...")
	r.Run(config.Config.Host)
}
