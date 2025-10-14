package postgresdataaccess

import (
	"database/sql"
	"fmt"
	"log"
	"todo/models"

	_ "github.com/lib/pq"
)

func (da *PostgresTodoDataAccess) GetTodos(status models.TodoCompletionFilter) ([]models.Todo, error) {
	query := `
        SELECT
            t.id, t.completed, t.value, tl.id, tl.title, t.created_at, t.completed_at
        FROM todo t
        INNER JOIN todo_list tl on t.todo_list_id = tl.id
    `

	if status == models.TodoCompleted {
		query = query + " where t.completed = true;"
	} else if status == models.TodoIncomplete {
		query = query + " where t.completed = false;"
	}

	log.Println("Executing Query: ", query)

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

		var tlId sql.NullInt32
		var tlTitle sql.NullString

		var tCreatedAt sql.NullTime
		var tCompletedAt sql.NullTime

		err := rows.Scan(
			&tId, &tCompleted, &tValue, &tlId, &tlTitle, &tCreatedAt, &tCompletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		todo := models.Todo{
			Id:        uint(tId.Int32),
			Completed: tCompleted.Bool,
			Value:     tValue.String,

			TodoListId:    uint(tlId.Int32),
			TodoListTitle: tlTitle.String,

			CreatedAt:   tCreatedAt.Time,
			CompletedAt: tCompletedAt.Time,
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return todos, nil
}

func (da *PostgresTodoDataAccess) GetTodoById(id uint) (models.Todo, error) {
	query := `
        SELECT
            t.id, t.completed, t.value, tl.id, tl.title, t.created_at, t.completed_at
        FROM todo t
        INNER JOIN todo_list tl on t.todo_list_id = tl.id
        WHERE t.id = $1;
    `
	var tId sql.NullInt32
	var tCompleted sql.NullBool
	var tValue sql.NullString

	var tlId sql.NullInt32
	var tlTitle sql.NullString

	var tCreatedAt sql.NullTime
	var tCompletedAt sql.NullTime

	log.Println("Executing Query: ", query)
	err := da.db.QueryRow(query, id).Scan(&tId, &tCompleted, &tValue, &tlId, &tlTitle, &tCreatedAt, &tCompletedAt)
	if err != nil {
		return models.Todo{}, fmt.Errorf("unable to get todo with that id: %w", err)
	}

	todo := models.Todo{
		Id:        uint(tId.Int32),
		Completed: tCompleted.Bool,
		Value:     tValue.String,

		TodoListId:    uint(tlId.Int32),
		TodoListTitle: tlTitle.String,

		CreatedAt:   tCreatedAt.Time,
		CompletedAt: tCompletedAt.Time,
	}

	return todo, nil
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
	_, err := da.db.Exec(`UPDATE todo SET completed = $1, completed_at = $4, value = $2 WHERE id = $3;`, todo.Completed, todo.Value, todo.Id, todo.CompletedAt)
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

func (da *PostgresTodoDataAccess) DeleteHistoricalTodos() error {
	_, err := da.db.Exec(`DELETE FROM todo WHERE completed = true;`)
	if err != nil {
		return err
	}
	return nil
}
