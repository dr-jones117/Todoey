package main

import (
	"log"
	"net/http"
	"todo/dataaccess"
)

var (
	port = "8080"
)

type Dependencies struct {
	dataAccess dataaccess.TodoDataAccess
}

func injectDependencies() Dependencies {
	// Change structs that implement the interfaces here!
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

	SetupHTTPHandlers(dependencies.dataAccess)

	log.Println("Server listening on:", port)
	err := http.ListenAndServe(":"+port, nil)

	if err != nil {
		log.Fatal(err)
	}
}
