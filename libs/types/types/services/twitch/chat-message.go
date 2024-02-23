package twitch

const TOPIC_CHAT_MESSAGE = "twitch.chat_message"

type TwitchChatMessage struct {
	BroadcasterUserId           string
	BroadcasterUserName         string
	BroadcasterUserLogin        string
	ChatterUserId               string
	ChatterUserName             string
	ChatterUserLogin            string
	MessageId                   string
	Message                     *ChatMessageMessage
	Color                       string
	Badges                      []ChatMessageBadge
	MessageType                 string
	Cheer                       *ChatMessageCheer
	Reply                       *ChatMessageReply
	ChannelPointsCustomRewardId string
}

type FragmentType int32

const (
	FragmentType_TEXT      FragmentType = 0
	FragmentType_CHEERMOTE FragmentType = 1
	FragmentType_EMOTE     FragmentType = 2
	FragmentType_MENTION   FragmentType = 3
)

type ChatMessageMessageFragmentEmote struct {
	Id         string
	EmoteSetId string
	OwnerId    string
	Format     []string
}

type ChatMessageMessageFragmentMention struct {
	UserId    string
	UserName  string
	UserLogin string
}

type ChatMessageMessageFragmentCheermote struct {
	Prefix string
	Bits   int64
	Tier   int64
}

type ChatMessageMessageFragment struct {
	Type      FragmentType
	Text      string
	Cheermote *ChatMessageMessageFragmentCheermote
	Emote     *ChatMessageMessageFragmentEmote
	Mention   *ChatMessageMessageFragmentMention
}

type ChatMessageMessage struct {
	Text      string
	Fragments []ChatMessageMessageFragment
}

type ChatMessageBadge struct {
	Id    string
	SetId string
	Info  string
}

type ChatMessageCheer struct {
	Bits int64
}

type ChatMessageReply struct {
	ParentMessageId   string
	ParentMessageBody string
	ParentUserId      string
	ParentUserName    string
	ParentUserLogin   string
	ThreadMessageId   string
	ThreadUserId      string
	ThreadUserName    string
	ThreadUserLogin   string
}
