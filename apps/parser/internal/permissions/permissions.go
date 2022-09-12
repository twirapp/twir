package permissions

import (
	"tsuwari/parser/pkg/helpers"
)

var CommandPerms = []string{"BROADCASTER", "MODERATOR", "SUBSCRIBER", "VIP", "FOLLOWER", "VIEWER"}

func UserHasPermissionToCommand(badges []string, commandPermission string) bool {
	commandPermIndex := helpers.IndexOf(CommandPerms, commandPermission)

	res := false
	for _, b := range badges {
		idx := helpers.IndexOf(CommandPerms, b)
		if idx <= commandPermIndex {
			res = true
			break
		}
	}

	return res
}
