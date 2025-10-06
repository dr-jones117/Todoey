package main

import (
	"log"
	"todo/dataaccess"
	"todo/dataaccess/postgresdataaccess"
	"todo/handlers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	port = "60235"
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

	store := cookie.NewStore([]byte("secret-key-23l4lkjlkj"))
	store.Options(sessions.Options{
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})

	router.Use(sessions.Sessions("todoey-session", store))
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlers.SetupHTTPHandlers(router, dependencies.dataAccess)

	log.Println("Server listening on:", port)
	err := router.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
}
