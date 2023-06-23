package client

type Database struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	// TODO: make an enum
	Adapter         string `json:"adapter"`
	Hostname        string `json:"hostname"`
	Database        string `json:"database"`
	ReviewsRequired int64  `json:"reviews_required"`
	Ssl             bool   `json:"ssl"`
	CaCertFile      string `json:"cacertfile"`
	KeyFile         string `json:"keyfile"`
	CertFile        string `json:"certfile"`
	RestrictAccess  bool   `json:"restrict_access"`
	AgentId         string `json:"agent_id"`
}
