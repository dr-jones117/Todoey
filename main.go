package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"todo/dataaccess"
)

var (
	port = "8080"
)

type Dependencies struct {
	dataAccess dataaccess.TodoDataAccess
}

func injectDependencies() Dependencies {
	todoDataAccess := &dataaccess.PostgresTodoDataAccess{}
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

	router := gin.Default()
	SetupHTTPHandlers(router, dependencies.dataAccess)

	log.Println("Server listening on:", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
