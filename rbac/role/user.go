package role

import (
	"github.com/traPtitech/traQ/rbac"
	"github.com/traPtitech/traQ/rbac/permission"
)

// User 一般ユーザーロール
const User = "user"

var userPerms = []rbac.Permission{
	// read, writeロールのパーミッションを全て含む
	permission.ChangeMyPassword,
	permission.GetMySessions,
	permission.DeleteMySessions,
	permission.GetMyTokens,
	permission.RevokeMyToken,
	permission.GetClients,
	permission.CreateClient,
	permission.EditMyClient,
	permission.DeleteMyClient,
	permission.CreateWebhook,
	permission.EditWebhook,
	permission.DeleteWebhook,
	permission.CreateBot,
	permission.EditBot,
	permission.DeleteBot,
	permission.ReissueBotToken,
}

func init() {
	userPerms = append(userPerms, readPerms...)
	userPerms = append(userPerms, writePerms...)
}
