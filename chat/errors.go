package chat

import "errors"

var UserNotFoundEroor = errors.New("User не найден")
var UserNotCreated = errors.New("User не создан")
var UserAlreadyExists = errors.New("User с таким именем уже есть")
var UsersNotFound = errors.New("Users пуст")

var MessageNotSended = errors.New("Сообщение не отправлено")
var MessageIsEmpty = errors.New("Сообщение пустое")
var MessageNotFound = errors.New("Сообщения не сущетсвует")
