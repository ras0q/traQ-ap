package role

import (
	"github.com/traPtitech/traQ/rbac"
	"github.com/traPtitech/traQ/rbac/permission"
)

// Read 読み取り専用ユーザーロール
const Read = "read"

var readPerms = []rbac.Permission{
	permission.GetChannel,
	permission.GetTopic,
	permission.GetMessage,
	permission.GetPin,
	permission.GetNotificationStatus,
	permission.ConnectNotificationStream,
	permission.GetUser,
	permission.GetMe,
	permission.GetClip,
	permission.GetClipFolder,
	permission.GetStar,
	permission.GetChannelVisibility,
	permission.GetUnread,
	permission.GetMutedChannels,
	permission.GetTag,
	permission.GetStamp,
	permission.GetMessageStamp,
	permission.GetMyStampHistory,
	permission.DownloadFile,
	permission.GetHeartbeat,
	permission.GetWebhook,
	permission.GetBot,
}
