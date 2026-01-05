package models

type ModUsersViewModel struct {
	RegisteredUsers []User
}

type ModUserViewModel struct {
	EditableUser User
}

func NewModUsersViewModel(registered_users []User) (m ModUsersViewModel) {
	m.RegisteredUsers = registered_users

	return m
}

func NewModUserViewModel(registered_user User) (m ModUserViewModel) {
	m.EditableUser = registered_user

	return m
}
