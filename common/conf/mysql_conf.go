package conf

/*
 *
 *	audit_status   		任务审核状态   notInvolved(数据库里面没有，后面新增的) 未参与  apply已报名 pending已参与 pass已通过 reject已驳回
 *  is_reward      		0 已发放奖励   1 未发放奖励
 *  language	  		en 英语   zh 中文  ar 阿拉伯语
 *  status        		0 未读   1 已读
 *  is_send_email       0 不发送  1 发送
 *	recipient_type		user 用户端  admin 管理端  all  全部
 *  message_type		notify  通知   announcement 公告
 *  status		        0 未发布  1 已发布  2 已移除 WonderfulWork 精彩作品
 *  is_draft		    0 非草稿  1 草稿   任务是否草稿
 */
const (
	NOTINVOLVED, APPLY, PENDING, PASS, REJECT = "notInvolved", "apply", "pending", "pass", "reject"
	YESREWARD, NOREWARD                       = 1, 0
	EN, ZH, AR                                = "en", "zh", "ar"
	UNREAD, READ                              = 0, 1
	SENDEMAIL, NOSENDEMAIL                    = 1, 0
	USER, ADMIN, ALl                          = "user", "admin", "all"
	NOTIFY, ANNOUNCEMENT                      = "notify", "announcement"
	WorkReleased, WorkPublished, WorkRemove   = 0, 1, 2
	DRAFT, NOTDRAFT                           = 1, 0
)
