package dispatchtypes

type DispatchType string

const (
	SystemAnnouncement   DispatchType = "system.announcement"
	CreateEmote          DispatchType = "emote.create"
	UpdateEmote          DispatchType = "emote.update"
	DeleteEmote          DispatchType = "emote.delete"
	CreateEmoteSet       DispatchType = "emote_set.create"
	UpdateEmoteSet       DispatchType = "emote_set.update"
	DeleteEmoteSet       DispatchType = "emote_set.delete"
	CreateUser           DispatchType = "user.create"
	UpdateUser           DispatchType = "user.update"
	DeleteUser           DispatchType = "user.delete"
	AddUserConnection    DispatchType = "user.add_connection"
	UpdateUserConnection DispatchType = "user.update_connection"
	DeleteUserConnection DispatchType = "user.delete_connection"
	CreateCosmetic       DispatchType = "cosmetic.create"
	UpdateCosmetic       DispatchType = "cosmetic.update"
	DeleteCosmetic       DispatchType = "cosmetic.delete"
	CreateEntitlement    DispatchType = "entitlement.create"
	UpdateEntitlement    DispatchType = "entitlement.update"
	DeleteEntitlement    DispatchType = "entitlement.delete"
)
