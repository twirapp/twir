package permissions

import (
	"tsuwari/parser/internal/types"
	"tsuwari/parser/pkg/helpers"
)

func UserHasPermissionToCommand(badges []string, commandPermission string) bool {
	commandPermIndex := helpers.IndexOf(types.CommandPerms, commandPermission)

	res := false
	for _, b := range badges {
		idx := helpers.IndexOf(types.CommandPerms, b)
		if idx <= commandPermIndex {
			res = true
			break
		}
	}

	return res
}
