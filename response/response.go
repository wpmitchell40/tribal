package response

type TribalResponse struct {
	Username               string
	RespondingUserUsername string
	Response               string
	OrgName                string
	Email                  string
}

type SlackSlashData struct {
	Token string
	TeamId string
	TeamDomain string
	EnterpriseId string
	EnterpriseName string
	ChannelId string
	ChannelName string
	UserId string
	UserName string
	Command string
	Text string
	ResponseUrl string
	TriggerId string
}
