package importer

type FailureReason string

const (
	FailureDuplicate            FailureReason = "DUPLICATE"
	FailurePlanLimit            FailureReason = "PLAN_LIMIT"
	FailureUnsupportedRole      FailureReason = "UNSUPPORTED_ROLE"
	FailureUnsupportedResponse  FailureReason = "UNSUPPORTED_RESPONSE_TYPE"
	FailureIncompatibleInterval FailureReason = "INCOMPATIBLE_INTERVALS"
	FailureInvalidRecord        FailureReason = "INVALID_RECORD"
)

type Failure struct {
	Name   string
	Reason FailureReason
}

type Report struct {
	ImportedCount int
	FailedCount   int
	Failures      []Failure
}

type RoleRequirement string

const (
	RoleEveryone    RoleRequirement = "EVERYONE"
	RoleSubscriber  RoleRequirement = "SUBSCRIBER"
	RoleVip         RoleRequirement = "VIP"
	RoleModerator   RoleRequirement = "MODERATOR"
	RoleBroadcaster RoleRequirement = "BROADCASTER"
)

type Command struct {
	Name, Response            string
	Enabled, Visible, IsReply bool
	Aliases                   []string
	Cooldown                  int
	Role                      RoleRequirement
	OnlineOnly, OfflineOnly   bool
}

type Timer struct {
	Name, Message                          string
	Enabled, OnlineEnabled, OfflineEnabled bool
	TimeInterval, MessageInterval          int
}
