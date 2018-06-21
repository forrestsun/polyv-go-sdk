package sdk

import (
	"time"
)

type RespMsg struct {
	Status_Code int    `json:"code"`
	Status      string `json:"status"`
	Message     string `json:"message"`
}

type AddCataMsg struct {
	RespMsg
	Data AddCataInfo `json:"data"`
}

type DelCataMsg struct {
	RespMsg
	Data bool `json:"data"`
}

type AddCataInfo struct {
	CataTree string `json:"catatree"`
	CataID   uint64 `json:"cataid"`
}

type CataMsg struct {
	RespMsg
	Data []CataNodes `json:"data"`
}

type CataNodes struct {
	CataNode
	CataNodes []CataNode `json:"nodes"`
}

type CataNode struct {
	Text        string `json:"text"`
	Cataname    string `json:"cataname"`
	Catatree    string `json:"catatree"`
	Cataid      uint64 `json:"cataid"`
	Parentid    uint64 `json:"parentid"`
	CataProfile string `json:"cataProfile"`
	Videos      int    `json:"videos"`
}

type StandVideoList struct {
	RespMsg
	Data  []StandVideoInfo `json:"data"`
	Total int              `json:"total"`
}

type StandVideoInfo struct {
	Vid                 string   `json:"vid" bson:"_id"`
	Tag                 string   `json:"tag"`      //视频标签
	MP4                 string   `json:"mp4"`      //MP4源文件
	Title               string   `json:"title"`    //标题
	Df                  int      `json:"df"`       //视频码率数
	Times               string   `json:"times"`    //播放次数
	MP41                string   `json:"mp4_1"`    //流畅码率mp4格式视频地址
	MP42                string   `json:"mp4_2"`    //高清码率mp4格式视频地址
	MP43                string   `json:"mp4_3"`    //超清码率mp4格式视频地址
	CataId              string   `json:"cataid"`   //分类id， 如1为根目录
	SWF_Link            string   `json:"swf_link"` //返回flash连接
	Status              string   `json:"status"`
	Seed                int      `json:"seed"`                //加密视频为1，非加密为0
	PlayerWidth         string   `json:"playerwidth"`         //视频宽度
	Duration            string   `json:"duration"`            //时长
	FirstImage          string   `json:"first_image"`         //视频首图
	Original_definition string   `json:"original_definition"` //最佳分辨率
	Context             string   `json:"context"`             //视频描述
	PlayerHeight        string   `json:"playerheight"`        //视频高度
	PTime               string   `json:"ptime"`               //视频上传日期
	Source_FileSize     int64    `json:"source_filesize"`
	FileSzie            []int64  `json:"filesize"` //filesize
	MD5CheckSum         string   `json:"md5checksum"`
	HLS                 []string `json:"hls"` //索引文件，记录每个清晰度的m3u8的链接
}

type VideoMsg struct {
	RespMsg
	Data []VideoInfo `json:"data"`
}

type VideoInfo struct {
	Vid string `json:"vid" bson:"_id"`
	// Image_b             []string `json:"images_b"`  //视频截图大图地址
	// Images              []string `json:"images"`    //视频截图
	// ImageUrls           []string `json:"imageUrls"` //视频截图
	Tag      string `json:"tag"`      //视频标签
	MP4      string `json:"mp4"`      //MP4源文件
	Title    string `json:"title"`    //标题
	Df       int    `json:"df"`       //视频码率数
	Times    string `json:"times"`    //播放次数
	MP41     string `json:"mp4_1"`    //流畅码率mp4格式视频地址
	MP42     string `json:"mp4_2"`    //高清码率mp4格式视频地址
	MP43     string `json:"mp4_3"`    //超清码率mp4格式视频地址
	CataId   string `json:"cataid"`   //分类id， 如1为根目录
	SWF_Link string `json:"swf_link"` //返回flash连接
	Status   string `json:"status"`
	Seed     int    `json:"seed"` //加密视频为1，非加密为0
	// FLV1                string   `json:"flv1"` //流畅码率flv格式视频地址
	// FLV2                string   `json:"flv2"` //高清码率flv格式视频地址
	// FLV3                string   `json:"flv3"` //超清码率flv格式视频地址
	SourceFile          string   `json:"sourcefile"`
	PlayerWidth         string   `json:"playerwidth"`         //视频宽度
	Default_Video       string   `json:"default_video"`       //用户默认播放视频
	Duration            string   `json:"duration"`            //时长
	FirstImage          string   `json:"first_image"`         //视频首图
	Original_definition string   `json:"original_definition"` //最佳分辨率
	Context             string   `json:"context"`             //视频描述
	PlayerHeight        string   `json:"playerheight"`        //视频高度
	PTime               string   `json:"ptime"`               //视频上传日期
	Source_FileSize     int64    `json:"source_filesize"`
	FileSzie            []int64  `json:"filesize"` //filesize
	MD5CheckSum         string   `json:"md5checksum"`
	HLS                 []string `json:"hls"` //索引文件，记录每个清晰度的m3u8的链接
	// Tsfilesize1         string   `json:"tsfilesize1"`
	// Tsfilesize2         string   `json:"tsfilesize2"`
	// Tsfilesize3         string   `json:"tsfilesize3"`
}

type VideoList struct {
	RespMsg
	Data  []VideoInfo `json:"data"`
	Total int         `json:"total"`
}

type NoPassVideoList struct {
	Status_Code string            `json:"error"`
	Data        []NoPassVideoInfo `json:"data"`
}

type NoPassVideoInfo struct {
	Swf_Link    string `json:"swf_link"`
	Duration    string `json:"duration"`
	Title       string `json:"title"`
	First_Image string `json:"first_image"`
	Times       string `json:"times"`
	Tag         string `json:"tag"`
	Context     string `json:"context"`
	Ptime       string `json:"ptime"`
	Vid         string `json:"vid"`
}

type DelVideoMsg struct {
	RespMsg
	Data string `json:"data"`
}

type DelVideoList struct {
	RespMsg
	Data  []DelVideoInfo `json:"data"`
	Total int            `json:"total"`
}

type DelVideoInfo struct {
	Tag                 string   `json:"tag"`                 //视频标签
	MP4                 string   `json:"mp4"`                 //MP4源文件
	Title               string   `json:"title"`               //标题
	Df                  int      `json:"df"`                  //视频码率数
	Times               string   `json:"times"`               //播放次数
	VID                 string   `json:"vid"`                 //视频id
	MP4_1               string   `json:"mp4_1"`               //流畅码率mp4格式视频地址
	MP4_2               string   `json:"mp4_2"`               //高清码率mp4格式视频地址
	MP4_3               string   `json:"mp4_3"`               //超清码率mp4格式视频地址
	CataID              string   `json:"cataid"`              //分类id， 如1为根目录
	SWF_Link            string   `json:"swf_link"`            //返回视频flash链接
	Status              string   `json:"status"`              //视频状态码（data中的status）
	Seed                int      `json:"seed"`                //加密视频为1，非加密为0
	PlayWidth           string   `json:"playerwidth"`         //视频宽度
	Duration            string   `json:"duration"`            //时长
	First_Image         string   `json:"first_image"`         //视频首图
	Original_Definition string   `json:"original_definition"` //最佳分辨率
	Context             string   `json:"context"`             //视频描述
	PlayHeight          string   `json:"playerheight"`        //视频高度
	PTime               string   `json:"ptime"`               //视频上传日期
	Source_FileSize     int      `json:"source_filesize"`     //源视频文件大小
	MD5CheckSum         string   `json:"md5checksum"`         //上传到POLYV云平台的视频源文件的MD5值，可以用来校验是否上传错误或完整
	HLS                 []string `json:"hls"`                 //流畅、高清、超清清晰度的m3u8
}

type VideoViewMsg struct {
	RespMsg
	Data []VideoViewInfo `json:"data"`
}

type VideoViewInfo struct {
	CurrentTime     string `json:"currentTime"`     //日期
	PcVideoView     int    `json:"pcVideoView"`     //pc端播放量
	MobileVideoView int    `json:"mobileVideoView"` //移动端播放量
}

type OfflineRankList struct {
	DateID     int64    `json:"id" bson:"_id"`
	Data       RankData `json:"data"`
	CreateDate time.Time
}

type RankMsg struct {
	RespMsg
	Data RankData `json:"data"`
}

type RankData struct {
	TotalMoVideoView  int           `json:"totalMoVideoView"` //移动端总播放量
	TotalPcVideoView  int           `json:"totalPcVideoView"` //pc端总播放量
	PcVideoDailys     []VideoDailys `json:"pcVideoDailys"`    //pc端播放量排行列表
	MobileVideoDailys []VideoDailys `json:"moVideoDailys"`    //移动端播放量排行列表
}

type VideoDailys struct {
	VideoId         string `json:"videoId"`         //视频vid
	Title           string `json:"title"`           //视频标题
	Duration        string `json:"duration"`        //播放时长
	PcVideoView     int    `json:"pcVideoView"`     //pc端播放量
	MobileVideoView int    `json:"mobileVideoView"` //移动端播放量
}

type DomainMsg struct {
	RespMsg
	Data []DomainList `json:"data"`
}

type OfflineDomainList struct {
	DateID     int64        `json:"id" bson:"_id"`
	Data       []DomainList `json:"data"`
	CreateDate time.Time
}

type DomainList struct {
	Domain             string `json:"domain"`             //域名
	PcPlayDuration     int    `json:"pcPlayDuration"`     //PC端播放时长（单位：秒）
	PcFlowSize         int    `json:"pcFlowSize"`         //PC端消耗流量（单位：字节）
	PcVideoView        int    `json:"pcVideoView"`        //PC端总播放量
	PcUniqueViewer     int    `json:"pcUniqueViewer"`     //PC端唯一观众数
	MobilePlayDuration int    `json:"mobilePlayDuration"` //移动端播放时长
	MobileVideoView    int    `json:"mobileVideoView"`    //移动端播放量
	MobileUniqueViewer int    `json:"mobileUniqueViewer"` //移动端播放者数量
}

type VideoLogMsg struct {
	RespMsg
	Data []VideoLogList `json:"data"`
}

type VideoLogList struct {
	PlayId          string `json:"playId"`       //表示此次播放动作的ID
	UserId          string `json:"userId"`       //用户ID
	VideoId         string `json:"videoId"`      //视频ID
	PlayDuration    int    `json:"playDuration"` //播放时长 (用户观看的总时间 ，例如：18：00开始看一个视频，看到了18：30，这30分钟就是播放时长)
	StayDuration    int    `json:"stayDuration"` //缓存时长
	CurrentTimes    int    `json:"currentTimes"` //播放时间 （用户观看的最后时间，例如：停止观看视频的时候，进度条最后的分钟数为35分钟，播放时间就是35分钟）
	Duration        int    `json:"duration"`     //视频总时长
	FlowSize        int    `json:"flowSize"`     //流量大小
	SessionId       string `json:"sessionId"//`  //用户自定义参数，如学员ID等
	Param1          string `json:"param1"`       //POLYV系统参数
	Param2          string `json:"param2"`
	Param3          string `json:"param3"`
	Param4          string `json:"param4"`
	Param5          string `json:"param5"`
	IpAddress       string `json:"ipAddress"`       // IP地址
	Country         string `json:"country"`         // 国家
	Province        string `json:"province"`        // 省份
	City            string `json:"city"`            // 城市
	Isp             string `json:"isp"`             // ISP运营商
	Referer         string `json:"referer"`         // 播放视频页面地址
	UserAgent       string `json:"userAgent"`       // 用户设备
	OperatingSystem string `json:"operatingSystem"` // 操作系统
	Browser         string `json:"browser"`         // 浏览器
	IsMobile        string `json:"isMobile"`        // 是否为移动端
	CurrentDay      string `json:"currentDay"`      // 日志查询日期 (格式为：yyyy-MM-dd)
	CurrentHour     int    `json:"currentHour"`     // 日志查看时间
	CreatedTime     int    `json:"createdTime"`     // 播放开始时间 (格式为13位的时间戳)
	LastModified    int    `json:"lastModified"`    // 日志更新日期 (格式为13位的时间戳)
}

type UrlFileInfo struct {
	FileUrl string
	CataId  string //设定上传视频的分类
	Async   bool
	FileInfo
}

type LocalFileInfo struct {
	FileName  string //文件名称
	WaterMark string //水印图片的URL，图片格式必须是png格式
	CataId    string //设定上传视频的分类
	FileInfo
}

type FileInfo struct {
	Title string `json:"title"`
	Tag   string `json:"tag"`
	Desc  string `json:"desc"`
}

type ReturnMsg struct {
	Status_Code int         `json:"error"`
	Data        []VideoInfo `json:"data"`
}

type UploadImgMsg struct {
	ErrorStr string `json:"error"`
	Data     bool   `json:"data"`
}

type UploadVideoMsg struct {
	Status_Code string      `json:"error"`
	Data        []VideoInfo `json:"data"`
}

type UploadAsyncVideoMsg struct {
	Status_Code string `json:"error"`
}

type HostMsg struct {
	Error        string `json:"error"`
	Host_Setting string `json:"Host_setting"`
}

/*setting_type
  0:无域名限制
  1:播放器启动禁止域名黑名单
  2:播放器启动允许域名白名单
  3:启动允许以及禁止播放域 (注：先判断允许播放域，再判断禁止播放域)
*/
type HostSetting struct {
	Disable_Host string `json:"disable_host"` //禁止播放的域名
	Enable_Host  string `json:"enable_host"`  //允许播放域名
	Setting_Type int    `json:"setting_type"` //域名设置类型
	UserID       string `json:"userid"`
}

type UseMsg struct {
	RespMsg
	Data UseInfo `json:"data"`
}

type UseInfo struct {
	TotalFlow  uint64 `json:"totalFlow"`  //用户总流量
	UsedSpace  uint64 `json:"usedSpace"`  //已用空间
	UsedFlow   uint64 `json:"usedFlow"`   //已用流量
	TotalSpace uint64 `json:"totalSpace"` //用户总空间
	UserId     string `json:"-"`          //POLYV用户ID
	Email      string `json:"-"`          //POLYV用户邮箱
}

type VideoImgMsg struct {
	RespMsg
	Data string `json:"data"`
}

type ChangeCataMsg struct {
	RespMsg
	Data bool `json:"data"`
}

type PolyvUserInfo struct {
	RespMsg
	PolyvInfo   `json:"data"`
	RemainBytes uint64 `json:"remainBytes"`
	LeftDay     int    `json:"leftday"`
}

type PolyvInfo struct {
	UserID     string
	WriteToken string
	ReadToken  string
	SecretKey  string
	Verbose    bool
	StartDate  string //传入保利威视第一个资源的日期
	CataList   map[string]string
}
