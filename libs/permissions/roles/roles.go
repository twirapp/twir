package roles

import "fmt"

type Permission int

const (
	ADMINISTRATOR Permission = 1 << iota

	UpdateChannelTitle
	UpdateChannelCategory

	ViewCommands
	ManageCommands

	ViewKeywords
	ManageKeywords

	ViewTimers
	ManageTimers

	ViewIntegrations
	ManageIntegrations

	ViewSongRequests
	ManageSongRequests

	ViewModeration
	ManageModeration

	ViewVariables
	ManageVariables

	ViewGreetings
	ManageGreetings
)

type Role struct {
	name        string
	permissions Permission
}

type User struct {
	name  string
	roles []Role
}

func (u User) HasPermission(p Permission) bool {
	for _, r := range u.roles {
		if r.permissions&p == p {
			return true
		}
	}
	return false
}

func createRole(name string, permissions Permission) Role {
	return Role{name, permissions}
}

func Test() {
	broadcaster := createRole("BROADCASTER", UpdateChannelTitle|ManageKeywords|ManageGreetings)

	user := User{"Alice", []Role{broadcaster}}

	fmt.Println(user.HasPermission(UpdateChannelTitle))
	fmt.Println(user.HasPermission(UpdateChannelCategory))
}
