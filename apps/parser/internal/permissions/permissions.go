package permissions

import (
	"fmt"

	"github.com/satont/tsuwari/apps/parser/pkg/helpers"
)

var CommandPerms = []string{"BROADCASTER", "MODERATOR", "VIP", "SUBSCRIBER", "FOLLOWER", "VIEWER"}

func UserHasPermissionToCommand(badges []string, commandPermission string) bool {
	commandPermIndex := helpers.IndexOf(CommandPerms, commandPermission)

	res := false
	for _, b := range badges {
		idx := helpers.IndexOf(CommandPerms, b)
		fmt.Println(idx, commandPermIndex)
		if idx == -1 {
			continue
		}
		if idx <= commandPermIndex {
			res = true
			break
		}
	}

	return res
}
