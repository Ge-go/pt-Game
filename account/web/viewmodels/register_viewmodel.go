package viewmodels

//POST register step 2
type CheckEmailAndPasswordReq struct {
	Sub             string `json:"sub" valid:"required~sub is blank" example:"106268783650461104364"`
	UserName        string `json:"username" valid:"required~username is blank,maxstringlength(30)" example:"long"`
	Password        string `json:"password" valid:"required~password is blank,minstringlength(8),maxstringlength(30)" example:"123456Abc@123"`
	ConfirmPassword string `json:"confirmPassword" valid:"required~confirmPassword is blank,minstringlength(8),maxstringlength(30)" example:"123456Abc@123"`
	Email           string `json:"email" valid:"required~email is blank,email" example:"123456@qq.com"`
	EmailCode       string `json:"emailCode" valid:"required~emailCode is blank,minstringlength(6),maxstringlength(6),numeric" example:"123456"`
}
