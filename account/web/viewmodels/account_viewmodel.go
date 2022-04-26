package viewmodels

// POST: /api/v1/user/register
type RegisterReq struct {
	Sub        string `json:"sub" valid:"required~sub is blank" example:"106268783650461104364"`
	Kind       int    `json:"kind" valid:"optional,range(1|1000000)" example:"2"`
	Gender     int    `json:"gender" valid:"optional,in(0|1|2)" example:"1"` //0男 1女 2其他
	Language   string `json:"language" valid:"required~Language is blank;optional,in(ar|de|en|es|fr|in|ma|po|py|ta|tu|yu|cn|cnh|cnt)" example:"en"`
	IsEuropean int    `json:"is_european" valid:"optional,in(0|1)" example:"1"` //勾选协议
}

// POST: /api/v1/user/verifyEmail
type VerifyEmailReq struct {
	Email     string `json:"email" valid:"required~email is blank,email" example:"1254@qq.com"`
	EmailCode string `json:"emailCode" valid:"optional,numeric,minstringlength(6)~email code must be 6-digit code" example:"123456"`
}

type VerifyEmailRsp struct {
	EmailCode string `json:"emailCode"`
}

type JwtToken struct {
	Token    string `json:"token"`
	UserName string `json:"userName"`
	IsLocked int    `json:"isLocked"`
}

// POST: /api/v1/user/Login
type LoginReq struct {
	Email        string `json:"email" valid:"required,email" example:"821230693@qq.com"`
	Password     string `json:"password" valid:"required,minstringlength(1),maxstringlength(30)" example:"123456Abc@123"`
	CaptchaId    string `json:"captcha_id" valid:"required" example:"1"`
	CaptchaValue string `json:"captcha_value" valid:"required" example:"xxx"`
}

// POST: /api/v1/user/resetPassword
type ResetPasswordReq struct {
	Email           string `json:"email" valid:"required~email is blank,email" example:"1254@qq.com"`
	EmailCode       string `json:"emailCode" valid:"required~emailCode is blank,numeric,minstringlength(6)~email code must be 6-digit code" example:"123456"`
	Password        string `json:"password" valid:"required~password is blank,minstringlength(8),maxstringlength(30)" example:"123456Abc@123"`
	ConfirmPassword string `json:"confirmPassword" valid:"required~confirmPassword is blank,minstringlength(8),maxstringlength(30)" example:"123456Abc@123"`
}

//POST register step 1
type MinorLimitReq struct {
	ShortName string `json:"shortName" valid:"required~shortName is blank" example:"Austria"`
	CertType  int    `json:"certType" valid:"required~certType is blank" example:"3"`
	Birthday  string `json:"birthday"  valid:"required~birthday is blank" example:"2012-12-23"`
	//Sub       string `json:"sub" valid:"required~sub is blank" example:"106268783650461104364"`
	//Channel   string `json:"channel" valid:"required~channel is blank" example:"106268783650461104364"`
	Region string `json:"region" valid:"required~Region is blank"  example:"Austria" `
}

//post register step Rsp
type MinorLimitRsp struct {
	Sub        string `json:"sub"`
	IsAdult    bool   `json:"isAdult" example:"true"`    //是否成年
	IsEuropean bool   `json:"isEuropean" example:"true"` //是否欧盟
	Age        int    `json:"age" example:"16"`          //年龄
}

type UserTag struct {
	TagId   int    `json:"value" gorm:"column:id"`
	TagName string `json:"text" gorm:"column:tag_name"`
}

//需要保存的注册信息
type RegisterInfo struct {
	Sub        string `json:"sub"`
	Country    string `json:"country" valid:"required~shortName is blank" example:"Austria"`
	Birthday   string `json:"birthday"  valid:"required~birthday is blank" example:"2012-12-23"`
	Region     string `json:"region"`
	Language   string `json:"language"`
	Password   string `json:"password" valid:"required~password is blank,minstringlength(8),maxstringlength(30)" example:"123456Abc@123"`
	Email      string `json:"email" valid:"required~email is blank,email" example:"123456@qq.com"`
	UserName   string `json:"username" valid:"required~userName is blank,maxstringlength(30)" example:"long"`
	IsAdult    int    `json:"isAdult" example:"1"` //是否成年
	Kind       int    `json:"kind" example:"3"`    //用户类型
	Gender     int    `json:"gender" example:"1"`  //性别
	IsEuropean bool   `json:"isEuropean" example:"true"`
	Level      string `json:"level" example:"1"`
}
