package server

const (
	ChatIDPrefix = "mundo_prd:" // Prefix for chat IDs
)

func GenerateChatID(ID string, uid string) string {
	return ChatIDPrefix + uid + ID
}
