package service

import (
	"github.com/katakuxiko/rest-api"
	"github.com/katakuxiko/rest-api/pkg/repository"
)
type Authorization interface{
	CreateUser(user todo.User) (int,error)
	GenerateToken(username string, password string) (string, error)
	ParseToken(token string) (int, error)
}

type TodoList interface{
	Create(userId int, list todo.TodoList)(int, error)
	GetAll(userId int)([]todo.TodoList, error)
	GetById(userId,listId int)(todo.TodoList, error)
	Delete(userId,listId int)error
	Update(userId, listId int, input todo.UpdateListInput) error
}

type TodoItem interface{

}

type Service struct {
	Authorization
	TodoList
	TodoItem
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList: NewTodoListService(repos.TodoList),
		
	}
}

