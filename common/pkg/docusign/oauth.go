package docusign

import "time"

// OAuth Base path constants
const (
	// Production/Live server base path
	PRODUCTION_OAUTH_BASE_PATH = "account.docusign.com"
	// Demo/Sandbox server base path
	DEMO_OAUTH_BASE_PATH = "demo.docusign.net"
	// Stage server base path
	STAGE_OAUTH_BASE_PATH = "account-s.docusign.com"
	//授权范围
	SCOPE_SIGNATURE     = "signature"
	SCOPE_EXTENDED      = "extended"
	SCOPE_IMPERSONATION = "impersonation"
	//jwt type
	GRANT_TYPE_JWT = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	TIME_OUT       = 30 * time.Second //请求的默认时间
	JWT_EXP_TIME   = 24 * time.Hour   //jwt token的过期时间
	//ACCOUNT_EXP_TIME = 1 * time.Hour    //ACCOUNT INFO 的过期时间
)

type Docusign struct {
	Config  *Config
	account Account
}

//docusign config
type Config struct {
	Env             string `json:"env"`               //那个环境
	ClientId        string `json:"client_id"`         //客户端ID
	RedirectUri     string `json:"redirect_uri"`      //回调地址
	SignRedirectUri string `json:"sign_redirect_uri"` //签名的回调地址
	Sub             string `json:"sub"`               //用户的唯一标识
	PublicKey       string `json:"public_key"`        //公钥 RSA
	PrivateKey      string `json:"private_key"`       //私钥 RSA
}

func Client(config Config) *Docusign {
	return &Docusign{
		Config: &config,
	}
}

//获取service account 授权链接
func (d *Docusign) GetServiceAccountAuthUrl() (url string) {
	//resourcePath := "/oauth/auth"
	oAuthBasePath := d.GetOAuthBasePath()
	url = "https://" + oAuthBasePath + "?response_type=code&scope=signature+impersonation&client_id=" +
		d.Config.ClientId + "&redirect_uri=" + d.Config.RedirectUri

	return
}

//根据环境配置获取当前OAUTH的基本路径
func (d *Docusign) GetOAuthBasePath() (oAuthBasePath string) {
	if d.Config.Env == "demo" {
		oAuthBasePath = DEMO_OAUTH_BASE_PATH
	} else if d.Config.Env == "stage" {
		oAuthBasePath = STAGE_OAUTH_BASE_PATH
	} else {
		oAuthBasePath = PRODUCTION_OAUTH_BASE_PATH
	}
	return
}
