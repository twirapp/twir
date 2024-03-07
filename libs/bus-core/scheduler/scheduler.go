package scheduler

const (
	CreateDefaultCommandsSubject = "scheduler.create_default_commands"
	CreateDefaultRolesSubject    = "scheduler.create_default_roles"
)

type CreateDefaultCommandsRequest struct {
	ChannelsIDs []string
}

type CreateDefaultRolesRequest struct {
	ChannelsIDs []string
}
