package viewmodels

// /api/v2/material/mateHome [get]
type MateHomeReq struct {
	Page       int    `json:"page"  valid:"required~page is blank,minstringlength(1)" example:"1"`
	PageSize   int    `json:"pageSize"  valid:"required~pageSize is blank" example:"10"`
	MateType   string `json:"mateType" valid:"required~mateType is blank" example:"All"`
	MateModule string `json:"mateModule" valid:"required~mateModule is blank" example:"All"`
}

type MateHomeRsp struct {
	Total        int64          `json:"total"  example:"100"`
	MateHomeList []MateHomeItem `json:"mateHomeList"`
}

type MateHomeItem struct {
	MateId          int64  `json:"mateId"  example:"1"`
	MateFileName    string `json:"mateFileName"  example:"test01"`
	MateFileFormat  string `json:"mateFileFormat"  example:"jpg"`
	ThumFileFormat  string `json:"thumFileFormat" example:"jpg"`
	MateFileSize    string `json:"mateFileSize"  example:"0.1M"`
	MateType        string `json:"mateType"  example:"载具"`
	MateModule      string `json:"mateModule"  example:"基础2D/音轨"`
	MatePermission  int    `json:"matePermission"  example:"0"`
	MateVersion     string `json:"mateVersion" example:"1.0.0"`
	MatePath        string `json:"matePath"  example:"material/Car/2D/v1.1.0/426aca11-d5bb-4c81-b61c-96da64669ff3.jpg"`
	ThumPath        string `json:"thumPath" example:"thumbnail/Car/2D/v1.1.0/426aca11-d5bb-4c81-b61c-96da64669ff3.jpg"`
	MateUuid        string `json:"mateUuid"  example:"426aca11-d5bb-4c81-b61c-96da64669ff3"`
	Language        string `json:"language"  example:"en"`
	MateStatus      int    `json:"mateStatus"  example:"0"`
	Uploader        string `json:"uploader"  example:"超级无敌管理员"`
	CreateTime      int64  `json:"createTime"  example:"1628925569000"`
	UpdateTime      int64  `json:"updateTime"  example:"1628925569000"`
	ThumDownloadUrl string `json:"thumdownloadUrl" example:"http://34.95.81.236/material/temp/5c218844-30c5-4c0c-a4fe-2cc6aa158ecd.png"`
	MateFavoriteNum int64  `json:"mateFavoriteNum" example:"10"`
	MateDownloadNum int64  `json:"mateDownloadNum" example:"23"`
	FavoriteStatus  int    `json:"favoriteStatus" example:"0"`
}
