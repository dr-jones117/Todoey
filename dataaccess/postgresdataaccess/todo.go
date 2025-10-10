package postgresdataaccess

import (
	"database/sql"
	"fmt"
	"todo/models"

	_ "github.com/lib/pq"
)

func (da *PostgresTodoDataAccess) GetTodos(status string) ([]models.Todo, error) {
	query := `
        SELECT
            t.id, t.completed, t.value, t.created_at, t.todo_list_id
        FROM todo AS t
    `

	if status == "completed" {
		query = query + " where t.completed = true"
	} else if status == "incomplete" {
		query = query + " where t.completed = false"
	}

	query = query + `;`

	rows, err := da.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("unable to get todo lists: %w", err)
	}
	defer rows.Close()

	var todos []models.Todo

	for rows.Next() {
		var tId sql.NullInt32
		var tCompleted sql.NullBool
		var tValue sql.NullString
		var tCreatedAt sql.NullTime
		var tTodo_list_Id sql.NullInt32

		err := rows.Scan(
			&tId, &tCompleted, &tValue, &tCreatedAt, &tTodo_list_Id,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		todo := models.Todo{
			Id:         uint(tId.Int32),
			Completed:  tCompleted.Bool,
			Value:      tValue.String,
			CreatedAt:  tCreatedAt.Time,
			TodoListId: uint(tTodo_list_Id.Int32),
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return todos, nil
}

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
