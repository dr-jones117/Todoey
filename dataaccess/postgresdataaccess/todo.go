package postgresdataaccess

import (
	"fmt"
	"todo/models"

	_ "github.com/lib/pq"
)

func (da *PostgresTodoDataAccess) CreateTodo(todo models.Todo) (models.Todo, error) {
	var id uint
	err := da.db.QueryRow(`INSERT INTO todo(completed, value, todo_list_id) VALUES($1, $2, $3) RETURNING Id;`, todo.Completed, todo.Value, todo.TodoListId).Scan(&id)
	todo.Id = id

	if err != nil {
		return todo, fmt.Errorf("unable to create todo: %w", err)
	}
	return todo, nil
}

func (da *PostgresTodoDataAccess) UpdateTodo(todo models.Todo) (models.Todo, error) {
	_, err := da.db.Exec(`UPDATE todo SET completed = $1, value = $2 WHERE id = $3;`, todo.Completed, todo.Value, todo.Id)
	if err != nil {
		return todo, fmt.Errorf("error")
	}

	return todo, nil
}
func (da *PostgresTodoDataAccess) DeleteTodo(id uint) error {
	_, err := da.db.Exec(`DELETE FROM todo WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("error")
	}
	return nil
}
