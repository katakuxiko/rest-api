package repository

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	todo "github.com/katakuxiko/rest-api"
	"github.com/sirupsen/logrus"
)

type TodoListPostgress struct {
	db *sqlx.DB
}

func NewTodoListPostgress(db *sqlx.DB) *TodoListPostgress {
	return &TodoListPostgress{db: db}
}

func (r *TodoListPostgress) Create(userId int,list todo.TodoList)(int ,error){
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title,description) VALUES ($1,$2) RETURNING id",todoListTable)
	row := tx.QueryRow(createListQuery,list.Title,list.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	createUsersListQuery := fmt.Sprintf("INSERT INTO %s (user_id,list_id) VALUES ($1,$2) ",userListTable)
	_,err = tx.Exec(createUsersListQuery,userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func(r *TodoListPostgress) GetAll(userId int)([]todo.TodoList, error){
	var lists []todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id ",
		todoListTable, userListTable)

	err := r.db.Select(&lists, query)
	
	return lists, err
}

func (r *TodoListPostgress)GetById(userId,listId int)(todo.TodoList, error){
	var list todo.TodoList

	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.list_id = $1",
		todoListTable, userListTable)

	err := r.db.Get(&list, query,listId)
	
	return list, err

}
func (r *TodoListPostgress)Delete(userId,listId int)error{
	query := fmt.Sprintf("delete from %s tl USING %s ul WHERE tl.id = ul.list_id AND ul.user_id = $1 AND ul.list_id = $2",
	todoListTable, userListTable)
	_, err := r.db.Exec(query, userId, listId)
	return err
}
func (r *TodoListPostgress) Update(userId int, listId int, input todo.UpdateListInput) error{
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	// title = "update title"
	// description = "update description"
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s tl SET %s from %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
	todoListTable, setQuery, userListTable, argId, argId+1)

	args = append(args, listId, userId)
	logrus.Debugf("Update query: %s", query)
	logrus.Debugf("Update query: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
