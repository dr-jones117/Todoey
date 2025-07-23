package main

import (
	"log"
	"todo/dataaccess"
	"todo/dataaccess/postgresdataaccess"
	"todo/handlers"

	"github.com/gin-gonic/gin"
)

var (
	port = "8080"
)

type Dependencies struct {
	dataAccess dataaccess.TodoDataAccess
}

func injectDependencies() Dependencies {
	todoDataAccess := &postgresdataaccess.PostgresTodoDataAccess{}
	if err := todoDataAccess.ConnectDataAccess(); err != nil {
		log.Fatal(err.Error())
	}

	return Dependencies{
		dataAccess: todoDataAccess,
	}
}

func shutdownDependencies(dependencies Dependencies) {
	dependencies.dataAccess.DisconnectDataAccess()
}

func main() {
	dependencies := injectDependencies()
	defer shutdownDependencies(dependencies)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.SetupHTTPHandlers(router, dependencies.dataAccess)

	log.Println("Server listening on:", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
