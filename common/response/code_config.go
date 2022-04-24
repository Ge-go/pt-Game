package response

import "net/http"

/*
 * 	错误码范围分类
 *  后四位数为零的错误码不使用（保留）
 *  前期有部分全局错误码（已在代码中引用，一期先不调整），后期支持多语言需全部改到这里来
 */
const (

	////////////////(业务错误)/////////////////
	businessError = 400000 // 业务错误   		4xxxxx
	/*
	 *	 400001 (passwords are inconsistent)
	 *	 400002 (Verification code error)
	 *	 400003 (The password contains at least numbers and letters, but the length is between 8 and 20 digits)
	 */
	DifferentPassword = 400004 //密码不一致

	ErrEmailHasExisted = 400010 //邮箱已经存在
	RegisteredTask     = 400011 // 已报名任务
	ErrEmailCode       = 400012 // 非法的邮箱验证码
	NoAdult            = 400015 // 未成年

	//youtube
	NoBindingYoutubeAccount = 400110 // 没有绑定YouTube账户
	NoYoutubeAuth           = 400112 //没有youtube授权
	NoYourYoutubeChannel    = 400113 // 不是你的youtube频道
	NoEnoughChannelFans     = 400114 // 你的频道粉丝数不够

	//合规协议
	NoModerationUserName = 400200 // 用户名不合规
	NoAgreementSigned    = 400201 //未签署欧盟协议

	NoDocusigned        = 400300 //签名失败
	SignStatusDifferent = 400301 //签名状态不一致
	UserAlreadySigned   = 400302 //用户已经签名

	////////////////(权限错误)/////////////////
	permissionError = 410000 // 权限错误   		41xxxx
	/*   全局已用
	 *   410001（Unauthorized）
	 *	 410002（Invalid UserName or Password）
	 *	 410003 (Jwt Token Not Found)
	 *	 410004 (jwt token expired)
	 *	 410005 (account was locked)
	 */
	WrongUserNameOrPassword = 410006 //用户名或密码错误
	AccountHasBeenFrozen    = 410008 //账户被冻结

	////////////////(参数错误)/////////////////
	//parameterError 		= 	420000  // 参数错误   		42xxxx
	/* 全局已用
	 * 420001（Bad Request）
	 * 420002（Invalid Captcha）
	 * 420003（Parameter error）
	 * 420004（invalid email code）
	 * 420005（email has existed）
	 * 420006（email has not register）
	 */
	PasswordsDoNotMatch = 420007 //两次密码不一致

	////////////////(限制错误)/////////////////
	limitError = 430000 // 限制错误
	/* 全局已用
	 * 430000(Forbidden)
	 */

	////////////////(系统错误)/////////////////
	systemError = 500000 // 系统错误   		5xxxx
	/* 全局已用
	 * 500001(Internal Server Error)
	 * 500002(Role Has Existed)
	 */

	NoDataFound = 510002 //数据不存在

	////////////////(数据库错误)/////////////////
	databaseError = 510000 // 数据库错误			51xxxx
	NoDataChanges = 510001 //没有数据变动

	////////////////(文件错误)/////////////////
	fileIoError           = 520000 // 文件错误
	UnsupportedFileSuffix = 520001 //不支持的文件后缀

	////////////////(第三方调用错误)/////////////////
	thirdPartyError = 530000 // 第三方调用错误

	////////////////(网络错误)/////////////////
	networkError = 540000 // 网络错误

	////////////////(素材库下载错误)/////////////////
	material  = 550000 // 素材库下载错误
	RunOut    = 550001 //今日下载次数已用完
	OverLimit = 550002 //单个文件超过下载限制

	StartTIMEIsGreaterThanEndTIME        = 420009 // 开始时间大于结束时间
	StartAndEndTimeIsLessThanCurrentTime = 420010 // 开始时间和结束时间小于当前时间
	PublishTimeIsLessThanCurrentTime     = 420011 // 发布时间小于当前时间

	TheCurrentStatusCannotBeUpdate = 429001

	WorksHasBeenRepeated = 400401 // 上传的作品是重复的

	FailedToParsedTheRequestData    = 490001 //解析请求数据失败
	FailedToValidatedTheRequestData = 490002 //验证请求数据失败
	FailedToGetTheUserIdInTheToken  = 490003 //获取token中的userid失败!

	TheCurrentStatusCanNotBeUpdate = 491001 //当前状态不可以修改
	TheTaskHasExpired              = 491002 //任务已过期

	FailedToGetYoutubeFans = 492001 //获取youtube粉丝数量失败

	FailedToInsertTheDataToTheDatabase = 499001 //写入数据失败
	FailedToUpdateTheDataToTheDatabase = 499002 //更新数据失败
	FailedToDeleteTheDataToTheDatabase = 499003 //删除数据失败
	FailedToQueryTheDataToTheDatabase  = 499004 //查询数据失败
)

var ConfigMessage = map[int]map[string]string{
	TheTaskHasExpired: map[string]string{
		"zh": "任务已过期!",
	},
	FailedToGetYoutubeFans: map[string]string{
		"zh": "获取youtube粉丝数量失败!",
	},
	FailedToQueryTheDataToTheDatabase: map[string]string{
		"zh": "查询数据失败!",
	},
	FailedToDeleteTheDataToTheDatabase: map[string]string{
		"zh": "更新数据失败!",
	},
	FailedToUpdateTheDataToTheDatabase: map[string]string{
		"zh": "更新数据失败!",
	},
	TheCurrentStatusCanNotBeUpdate: map[string]string{
		"zh": "当前状态不可以更新!",
	},
	FailedToGetTheUserIdInTheToken: map[string]string{
		"zh": "获取token中的userid失败!",
	},
	FailedToInsertTheDataToTheDatabase: map[string]string{
		"zh": "写入数据失败",
	},
	FailedToParsedTheRequestData: map[string]string{
		"zh": "解析请求数据失败",
	},
	FailedToValidatedTheRequestData: map[string]string{
		"zh": "验证请求数据失败",
	},

	NoDataFound: {
		"en": "No data found",
		"zh": "数据不存在",
		"ar": "لا توجد تغييرات في البيانات",
	},
	WorksHasBeenRepeated: {
		"en": "Works has been repeated",
		"zh": "作品被重复上传了",
	},
	NoAgreementSigned: map[string]string{
		"en": "No agreement signed",
		"zh": "未签署协议",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	NoModerationUserName: map[string]string{
		"en": "no Moderation UserName",
		"zh": "不合规的用户名",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	ErrEmailCode: map[string]string{
		"en": "invalid email code",
		"zh": "非法的邮箱验证码",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	// 400000~410000
	ErrEmailHasExisted: map[string]string{
		"en": "email has existed",
		"zh": "邮箱已经存在",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	DifferentPassword: map[string]string{
		"en": "passwords are inconsistent",
		"zh": "密码不一致",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	StartAndEndTimeIsLessThanCurrentTime: {
		"en": "Start time and end time are less than current time",
		"zh": "开始时间和结束时间小于当前时间",
		"ar": "وقت البدء ووقت الانتهاء أقل من الوقت الحالي",
	},
	PublishTimeIsLessThanCurrentTime: {
		"en": "Publish time is less than current time",
		"zh": "发布时间小于现在的时间",
		"ar": "وقت النشر أقل من الوقت الحالي",
	},
	NoBindingYoutubeAccount: map[string]string{
		"en": "No binding youtube account",
		"zh": "没有绑定youtube账户",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	RegisteredTask: map[string]string{
		"en": "Registered task",
		"zh": "已报名任务",
		"ar": "مهمة مسجلة",
	},
	StartTIMEIsGreaterThanEndTIME: {
		"en": "Start time is greater than end time",
		"zh": "开始时间大于结束时间",
		"ar": "وقت البدء أكبر من وقت الانتهاء",
	},
	NoYoutubeAuth: map[string]string{
		"en": "No youtube auth",
		"zh": "没有youtube授权",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	NoYourYoutubeChannel: map[string]string{
		"en": "No your youtube channel",
		"zh": "不是你的youtube 频道",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	NoEnoughChannelFans: map[string]string{
		"en": "your youtube channel fans is not Enough",
		"zh": "你的youtube 频道粉丝不足",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	NoAdult: map[string]string{
		"en": "No adult",
		"zh": "未成年",
		"ar": "لا يوجد حساب يوتيوب ملزم",
	},
	// 410000~420000
	WrongUserNameOrPassword: map[string]string{
		"en": "Wrong username or password",
		"zh": "密码错误或账号不存在",
		"ar": "اسم المستخدم خاطئ أو كلمة المرور خاطئة",
	},
	PasswordsDoNotMatch: {
		"en": "Passwords don't match",
		"zh": "两次输入的密码不同",
		"ar": "كلمتا المرور غير متطابقتين",
	},
	AccountHasBeenFrozen: {
		"en": "Your account has been frozen",
		"zh": "您的账户已被冻结",
		"ar": "تم تجميد حسابك",
	},
	// 510000~520000
	NoDataChanges: {
		"en": "No data changes",
		"zh": "数据无修改",
		"ar": "لا توجد تغييرات في البيانات",
	},
	// 520000~530000
	UnsupportedFileSuffix: {
		"en": "Unsupported file suffix",
		"zh": "不支持的文件后缀",
		"ar": "لاحقة ملف غير مدعومة",
	},
	// 550000~560000
	RunOut: {
		"en": "Today’s downloads have been exhausted",
		"zh": "今日下载次数已用完",
		"ar": "تم استنفاد التنزيلات اليوم",
	},
	OverLimit: {
		"en": "A single file exceeds the download limit of the day",
		"zh": "单个文件超过当日下载限制",
		"ar": "يتجاوز ملف واحد الحد الأقصى المسموح به للتنزيل في اليوم",
	},
}

func GetMessage(code int, language string) *ErrorNo {

	httpStatusCode := http.StatusOK
	if code >= 500000 {
		httpStatusCode = http.StatusInternalServerError
	}
	//switch {
	//case code >= 500000 && code < 540000:
	//	httpStatusCode = http.StatusInternalServerError
	//	break
	//case code >= 540000 && code < 550000:
	//	httpStatusCode = http.StatusGatewayTimeout
	//	break
	//}

	return &ErrorNo{
		HTTPStatusCode: httpStatusCode, //客户端错误同意返回 200
		ServiceCode:    code,
		Message:        ConfigMessage[code][language],
	}
}
