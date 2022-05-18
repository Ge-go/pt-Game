package docusign

/**
// 获取当前账号信息
*/
type UserInfo struct {
	Sub        string    `json:"sub"`   //用户唯一标识
	Email      string    `json:"email"` //用户邮箱
	Name       string    `json:"name"`  //用户名称
	GivenName  string    `json:"given_name"`
	FamilyName string    `json:"family_name"`
	Created    string    `json:"created"`
	Accounts   []Account `json:"accounts"`
}

// 一个用户可以拥有多个账号
type Account struct {
	AccountId   string `json:"account_id"`   //账号
	IsDefault   bool   `json:"is_default"`   //是否默认
	AccountName string `json:"account_name"` //账号名称
	BaseUri     string `json:"base_uri"`     //配置的url
}
