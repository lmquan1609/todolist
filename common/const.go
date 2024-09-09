package common

const (
	CurrentUser = "current_user"
)

type Requester interface {
	GetUserId() int
	GetEmail() string
	GetRole() string
}

type DBType int

const (
	DBTypeItem DBType = 1
	DBTypeUser DBType = 2
)

const (
	PluginDBMain  = "mysql"
	PluginJWT     = "jwt"
	PluginPubsub  = "pubsub"
	PluginItemAPI = "item-api"

	TopicUserLikedItem   = "TopicUserLikedItem"
	TopicUserUnlikedItem = "TopicUserUnlikedItem"
)
