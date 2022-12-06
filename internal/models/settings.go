package models

// Settings app settings
type Settings struct {
	AppParams         Params            `json:"app"`
	SendSMSURL        SendSMSSettings   `json:"sendSMS"`
	OracleCFTDbParams OracleCFTDbParams `json:"oracleCFTdb"`
	SecretKey         SecretKey         `json:"secretKey"`
	FileService       FileServiceParams `json:"fileService"`
}

// Params contains params of server meta data
type Params struct {
	ServerName string `json:"serverName"`
	PortRun    string `json:"portRun"`
	LogFile    string `json:"logFile"`
	ServerURL  string `json:"serverURL"`
}

// SendSMS
type SendSMSSettings struct {
	Url string `json:"url"`
}

// OracleCFTDbParams contains params of oracle db server
type OracleCFTDbParams struct {
	Server              string `json:"server"`
	User                string `json:"user"`
	Password            string `json:"password"`
	ConnectionCheckTime int    `json:"connectionCheckTime"`
}

// SecretKey for JWT
type SecretKey struct {
	Key string `json:"key"`
}

// FtpServiceParams contains params of sending file to ftp server
type FileServiceParams struct {
	BaseURL       string `json:"baseURL"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	TokenHeadName string `json:"tokenHeadName"`
	MegafonLife   string `json:"megafonLife"`
	HumoOnline    string `json:"humoOnline"`
}
