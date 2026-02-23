package models

type ModUsersViewModel struct {
	RegisteredUsers []User
}

type ModUserViewModel struct {
	EditableUser User
}

type ModBoardsViewModel struct {
	RegisteredBoards []Board
}

type ModBoardViewModel struct {
	EditableBoard BoardDto
}

func NewModUsersViewModel(registered_users []User) (m ModUsersViewModel) {
	m.RegisteredUsers = registered_users

	return m
}

func NewModUserViewModel(registered_user User) (m ModUserViewModel) {
	m.EditableUser = registered_user

	return m
}

func NewModBoardsViewModel(registered_boards []Board) (m ModBoardsViewModel) {
	m.RegisteredBoards = registered_boards

	return m
}

func NewModBoardViewModel(registered_board BoardDto) (m ModBoardViewModel) {
	m.EditableBoard = registered_board

	return m
}
