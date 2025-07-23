package postgresdataaccess

import (
	"database/sql"
	"fmt"
	"sort"
	"time"
	"todo/models"

	_ "github.com/lib/pq"
)

func (da *PostgresTodoDataAccess) GetTodoLists() ([]models.TodoList, error) {
	rows, err := da.db.Query(`
        SELECT
            tl.id, tl.title, tl.created_at,
            t.id, t.completed, t.value, t.created_at, t.todo_list_id
        FROM todo_lists AS tl
        LEFT JOIN todos AS t ON t.todo_list_id = tl.id;
    `)
	if err != nil {
		return nil, fmt.Errorf("unable to get todo lists: %w", err)
	}
	defer rows.Close()

	todoListMap := make(map[int]*models.TodoList)

	for rows.Next() {
		var tlId int
		var tlTitle string
		var tlCreatedAt time.Time

		var tId sql.NullInt32
		var tCompleted sql.NullBool
		var tValue sql.NullString
		var tCreatedAt sql.NullTime
		var tTodo_list_Id sql.NullInt32

		err := rows.Scan(
			&tlId, &tlTitle, &tlCreatedAt,
			&tId, &tCompleted, &tValue, &tCreatedAt, &tTodo_list_Id,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// If the list doesn't exist in the map, add it
		if _, exists := todoListMap[tlId]; !exists {
			todoListMap[tlId] = &models.TodoList{
				Id:        uint(tlId),
				Title:     tlTitle,
				CreatedAt: tlCreatedAt,
				Todos:     []models.Todo{},
			}
		}

		if uint(tId.Int32) != 0 {
			todo := models.Todo{
				Id:         uint(tId.Int32),
				Completed:  tCompleted.Bool,
				Value:      tValue.String,
				CreatedAt:  tCreatedAt.Time,
				TodoListId: uint(tTodo_list_Id.Int32),
			}

			// Append the todo to the correct list
			todoListMap[tlId].Todos = append(todoListMap[tlId].Todos, todo)
		}

	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	// Convert map to slice
	todoLists := make([]models.TodoList, 0, len(todoListMap))
	for _, list := range todoListMap {
		todoLists = append(todoLists, *list)
	}

	// We have to do order bys in go since maps are intentionally randomized
	for i := range todoLists {
		sort.Slice(todoLists[i].Todos, func(j, k int) bool {
			return todoLists[i].Todos[j].Id < todoLists[i].Todos[k].Id
		})
	}

	sort.Slice(todoLists, func(i, j int) bool {
		return todoLists[i].Id < todoLists[j].Id
	})

	return todoLists, nil
}

func (da *PostgresTodoDataAccess) CreateTodoList(todoList models.TodoList) (models.TodoList, error) {
	var id uint

	// Correct SQL with VALUES() and RETURNING
	err := da.db.QueryRow(
		`INSERT INTO todo_lists(title) VALUES($1) RETURNING id`,
		todoList.Title,
	).Scan(&id)
	if err != nil {
		return todoList, fmt.Errorf("failed to create the todo list: %w", err)
	}

	todoList.Id = id
	return todoList, nil
}

func (da *PostgresTodoDataAccess) UpdateTodoList(todoListId uint, title string) error {
	_, err := da.db.Exec(`UPDATE todo_lists SET title = $1 WHERE id = $2;`, title, todoListId)
	if err != nil {
		return fmt.Errorf("failed to update the todo list: %w", err)
	}
	return nil
}

func (da *PostgresTodoDataAccess) DeleteTodoList(id uint) error {
	_, err := da.db.Exec(`DELETE FROM todo_lists WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("failed to delete the todolist: %w", err)
	}
	return nil
}
