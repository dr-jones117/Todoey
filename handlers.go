package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"todo/dataaccess"
	"todo/models"
	"todo/templates"
)

var (
	tmpl           *template.Template
	todoDataAccess dataaccess.TodoDataAccess
)

func loadTemplates() error {
	templatesFiles, err := filepath.Glob("templates/*.html")
	if err != nil {
		return err
	}

	iconFiles, err := filepath.Glob("templates/icons/*.html")
	if err != nil {
		return err
	}

	allFiles := append(templatesFiles, iconFiles...)

	tmpl, err = template.ParseFiles(allFiles...)
	if err != nil {
		return err
	}

	return nil
}

func SetupHTTPHandlers(tda dataaccess.TodoDataAccess) {
	todoDataAccess = tda

	if err := loadTemplates(); err != nil {
		log.Fatalf("Failed to load templates")
	}

	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoGet(w, r)
		case http.MethodPost:
			todoPost(w, r)
		case http.MethodPut:
			todoPut(w, r)
		case http.MethodDelete:
			todoDelete(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		}
	})

	http.HandleFunc("/todo-lists", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoListsGet(w, r)
		case http.MethodPost:
			todoListsPost(w, r)
		case http.MethodPut:
			todoListsPut(w, r)
		case http.MethodDelete:
			todoListsDelete(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusBadRequest)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.ExecuteTemplate(w, "index", nil)
	})

}

func todoListsGet(w http.ResponseWriter, r *http.Request) {
	todoLists, err := todoDataAccess.GetTodoLists()
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	if err = tmpl.ExecuteTemplate(w, "todoLists", MapTodoListsTemplate(todoLists)); err != nil {
		writeInternalServerError(w, err.Error())
		return
	}
}

func todoListsPost(w http.ResponseWriter, r *http.Request) {
	var todoList models.TodoList
	todoList, err := todoDataAccess.CreateTodoList(todoList)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	if err := tmpl.ExecuteTemplate(w, "todoList", MapTodoListTemplate(todoList)); err != nil {
		writeInternalServerError(w, "Unable to create todo list template")
		return
	}
}

func todoListsPut(w http.ResponseWriter, r *http.Request) {
	strId := r.FormValue("id")
	title := r.FormValue("title")

	if strId == "" {
		writeInternalServerError(w, "please provide a todolist id")
		return
	}

	id, err := strconv.Atoi(strId)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	todoDataAccess.UpdateTodoList(uint(id), title)

}

func todoListsDelete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	id, ok := params["todolistid"]
	if !ok {
		writeInternalServerError(w, "please provide a todolist id")
		return
	}

	idStr, err := strconv.Atoi(id[0])
	if err != nil {
		writeInternalServerError(w, "invalid id")
		return

	}
	err = todoDataAccess.DeleteTodoList(uint(idStr))
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}
}

func todoGet(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func todoPost(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	var err error
	var todoTemplateData templates.TodoTemplateData

	strTodoListId := r.FormValue("todolistid")
	if strTodoListId == "" {
		writeInternalServerError(w, "No todo list id was supplied")
		return
	}

	todoListId, err := strconv.Atoi(strTodoListId)
	if err != nil {
		writeInternalServerError(w, "invalid todo list id")
	}

	params := r.URL.Query()
	_, ok := params["setFocus"]
	if ok {
		todoTemplateData.FocusInput = true
	}

	todo.TodoListId = uint(todoListId)
	todo, err = todoDataAccess.CreateTodo(todo)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	todoTemplateData.Todo = todo
	tmpl.ExecuteTemplate(w, "todo", todoTemplateData)

}

func todoPut(w http.ResponseWriter, r *http.Request) {
	todo, err := MapTodoFromRequestForm(r)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	todo, err = todoDataAccess.UpdateTodo(todo)
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	tmpl.ExecuteTemplate(w, "todo", todo)
}

func todoDelete(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	idValues, ok := params["id"]
	if !ok {
		writeInternalServerError(w, "Please provide an id")
		return
	}

	paramId, err := strconv.Atoi(idValues[0])
	if err != nil {
		writeInternalServerError(w, err.Error())
		return
	}

	idTofind := uint(paramId)
	if err := todoDataAccess.DeleteTodo(idTofind); err != nil {
		writeInternalServerError(w, err.Error())
		return
	}
}

func writeInternalServerError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
	log.Println(msg)
}
