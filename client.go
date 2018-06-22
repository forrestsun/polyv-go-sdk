package sdk

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/smallnest/goreq"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	GET  = "get"
	POST = "post"
)

var params map[string]string

//echo -n 'mder.mder/1' | md5sum/sha1sum
// 为避免明文读取，参数passwd为SHA1生成密码生成后的参数
// 参数passwdmd5为密码的32位MD5校验码
func Login(email, passwd, passwdmd5 string) *PolyvUserInfo {
	var userinfo PolyvUserInfo
	if email == "" || passwd == "" || passwdmd5 == "" || len(passwd) < 20 {
		return &PolyvUserInfo{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     "参数不能为空或有错误",
			},
		}
	}

	passwd = strings.ToUpper(passwd[0:20])
	passwdmd5 = strings.ToLower(passwd)
	pwdmd5 := ""
	goreq.New().Post(fmt.Sprintf("http://api.polyv.net/v2/user/login?email=%s&password=%s&passwordMd5=%s",
		email, passwd, pwdmd5)).BindBody(&userinfo).
		// SetDebug(true).
		// SetCurlCommand(true).
		End()

	return &userinfo
}

func GetSign(value string) string {
	return strings.ToUpper(CryptoSHA1(value))
}

//对字符串进行SHA1哈希
func CryptoSHA1(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	return fmt.Sprintf("%x", t.Sum(nil))
}

//查询播放域名限制
func (self *PolyvInfo) GetHostUrl() *HostMsg {
	var hostmsg HostMsg
	url := "http://v.polyv.net/uc/services/rest?method=getHostUrl"
	_, _, errs := goreq.New().Post(url).
		BindBody(&hostmsg).
		Query("readtoken=" + self.ReadToken).
		SetDebug(self.Verbose).
		SetCurlCommand(self.Verbose).
		End()

	if len(errs) > 0 {
		return &HostMsg{
			Error:        errs[0].Error(),
			Host_Setting: "",
		}
	}
	return &hostmsg
}

//上传多个视频的预览图URL
func (self *PolyvInfo) UploadConverImageUrl(vids, cataids, img_url string) *RespMsg {
	params = make(map[string]string)
	respmsg := RespMsg{}
	str := ""

	ptime := time.Now().Unix() * 1000

	if vids != "" {
		str = fmt.Sprintf("fileUrl=%s&ptime=%d&vids=%s%s", img_url, ptime, vids, self.SecretKey)
		params["vids"] = vids

	} else if cataids != "" {
		str = fmt.Sprintf("fileUrl=%s&ptime=%d&cataids=%s%s", img_url, ptime, vids, self.SecretKey)
		params["cataids"] = cataids
	}

	params["fileUrl"] = img_url
	params["ptime"] = fmt.Sprintf("%d", ptime)

	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/uploadCoverImageUrl", self.UserID)
	self.request(POST, url, str, params, &respmsg)
	return &respmsg

}

func (self *PolyvInfo) UploadMultiUrlFile(title, file_url, cataid string) *RespMsg {
	params = make(map[string]string)
	respmsg := RespMsg{}
	ptime := time.Now().Unix() * 1000
	url := fmt.Sprintf("http://api.polyv.net/v2/video/grab/%s/upload/multi", self.UserID)
	str := fmt.Sprintf("cataid=%s&fileUrl=%s&ptime=%d&title=%s%s", cataid, file_url, ptime, title, self.SecretKey)

	self.request(POST, url, str, params, &respmsg)
	return &respmsg
}

//获取用户空间及流量情况
func (self *PolyvInfo) GetUseInfo(query_date string) *UseMsg {
	params = make(map[string]string)
	usemsg := UseMsg{}
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("date=%s&ptime=%d%s", query_date, ptime, self.SecretKey)
	url := fmt.Sprintf("http://api.polyv.net/v2/user/%s/main", self.UserID)
	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["date"] = query_date

	self.request(GET, url, str, params, &usemsg)
	return &usemsg
}

//获取用户空间及流量情况
func (self *PolyvInfo) GetTotalUseInfo() *UseMsg {
	params = make(map[string]string)
	usemsg := UseMsg{}
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("ptime=%d%s", ptime, self.SecretKey)
	url := fmt.Sprintf("http://api.polyv.net/v2/user/%s/main", self.UserID)
	params["ptime"] = fmt.Sprintf("%d", ptime)

	self.request(GET, url, str, params, &usemsg)

	return &usemsg
}

//获取单个视频的首图
func (self *PolyvInfo) GetVideoImage(vid, t string) *VideoImgMsg {
	params = make(map[string]string)
	vimgmsg := VideoImgMsg{}
	ptime := time.Now().Unix() * 1000

	if t != "1" {
		t = "2"
	}

	params["t"] = t
	params["vid"] = vid
	params["ptime"] = fmt.Sprintf("%d", ptime)

	str := fmt.Sprintf("ptime=%d&t=%s&vid=%s%s", ptime, t, vid, self.SecretKey)
	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/get-image", self.UserID)
	self.request(GET, url, str, params, &vimgmsg)

	return &vimgmsg
}

func (self *PolyvInfo) request(reqtype, url, str_sign string, param map[string]string, body interface{}) {
	var req *goreq.GoReq

	if reqtype == "post" {
		req = goreq.New().Post(url).BindBody(body)
	} else {
		req = goreq.New().Get(url).BindBody(body)
	}

	for k, v := range param {
		req.Query(fmt.Sprintf("%s=%s", k, v))
	}
	sign := GetSign(str_sign)
	req.Query("sign=" + sign)
	req.SetDebug(self.Verbose).SetCurlCommand(self.Verbose).End()
}

//获取单个视频信息
func (self *PolyvInfo) GetVideoInfo(vid string) *VideoMsg {
	params = make(map[string]string)
	videomsg := VideoMsg{}

	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("format=%s&ptime=%d&vid=%s%s", "json", ptime, vid, self.SecretKey)

	params["format"] = "json"
	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["vid"] = vid
	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/get-video-msg", self.UserID)

	self.request(GET, url, str, params, &videomsg)

	return &videomsg
}

//获取最新视频/全部视频列表
func (self *PolyvInfo) GetVideoList(catatree, pageSize, pageNum, startDate, endDate string) *VideoList {
	params = make(map[string]string)
	videolist := VideoList{}

	// self.loginfo(catatree)
	if catatree == "" {
		catatree = "1" //默认取根目录下所有资源
	}

	ptime := fmt.Sprintf("%d", time.Now().Unix()*1000)

	if pageSize == "" {
		pageSize = "10"
	}

	if pageNum == "" {
		pageNum = "1"
	}

	str := ""

	if startDate == "" && endDate == "" {
		str = fmt.Sprintf("catatree=%s&numPerPage=%s&pageNum=%s&ptime=%s", catatree, pageSize, pageNum, ptime)
	} else if endDate != "" && startDate == "" {
		str = fmt.Sprintf("catatree=%s&endDate=%s&numPerPage=%s&pageNum=%s&ptime=%s", catatree, endDate, pageSize, pageNum, ptime)
	} else if startDate != "" && endDate == "" {
		str = fmt.Sprintf("catatree=%s&numPerPage=%s&pageNum=%s&ptime=%s&startDate=%s", catatree, pageSize, pageNum, ptime, startDate)
	} else if startDate != "" && endDate != "" {
		str = fmt.Sprintf("catatree=%s&endDate=%s&numPerPage=%s&pageNum=%s&ptime=%s&startDate=%s", catatree, endDate, pageSize, pageNum, ptime, startDate)
	}

	param_str := str
	str = str + self.SecretKey

	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/get-new-list?%s", self.UserID, param_str)
	self.request(GET, url, str, nil, &videolist)

	return &videolist
}

//按标题查找视频
func (self *PolyvInfo) SearchByTitle(title, pageSize, pageNum string) *StandVideoList {
	videolist := StandVideoList{}

	if pageSize == "" {
		pageSize = "10"
	}

	if pageNum == "" {
		pageNum = "1"
	}

	str := fmt.Sprintf("keyword=%s&numPerPage=%s&pageNum=%s&ptime=%s",
		title, pageSize, pageNum, fmt.Sprintf("%d", time.Now().Unix()*1000))
	param_str := str
	str = str + self.SecretKey

	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/search?%s", self.UserID, param_str)
	self.request(GET, url, str, nil, &videolist)

	return &videolist
}

func (self *PolyvInfo) AddCata(cata_name string) *AddCataMsg {
	resp := AddCataMsg{}
	params = make(map[string]string)
	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/addCata", self.UserID)
	ptime := time.Now().Unix() * 1000

	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["cataname"] = cata_name
	params["parentid"] = "1"

	str := fmt.Sprintf("cataname=%s&parentid=1&ptime=%d%s", cata_name, ptime, self.SecretKey)
	self.request(POST, url, str, params, &resp)

	return &resp
}

func (self *PolyvInfo) DelCata(cataid string) *DelCataMsg {
	resp := DelCataMsg{}
	params = make(map[string]string)
	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/deleteCata", self.UserID)
	ptime := time.Now().Unix() * 1000

	str := fmt.Sprintf("cataid=%s&ptime=%d&userid=%s%s", cataid, ptime, self.UserID, self.SecretKey)

	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["cataid"] = cataid
	params["userid"] = self.UserID

	self.request(POST, url, str, params, &resp)

	return &resp
}

// 获取视频分类目录
func (self PolyvInfo) CataJson() *CataMsg {
	catamsg := CataMsg{}
	params = make(map[string]string)
	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/cataJson", self.UserID)
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("ptime=%d&userid=%s%s", ptime, self.UserID, self.SecretKey)
	params["ptime"] = fmt.Sprintf("%d", ptime)
	self.request(GET, url, str, params, &catamsg)
	return &catamsg
}

// 获取视频回收站列表
func (self *PolyvInfo) GetDelList(pageNum, pageSize string) *DelVideoList {
	params = make(map[string]string)
	delvideolist := DelVideoList{}

	if pageSize == "" {
		pageSize = "10"
	}

	if pageNum == "" {
		pageNum = "1"
	}

	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/get-del-list", self.UserID)
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("format=json&numPerPage=%s&pageNum=%s&ptime=%d%s", pageSize, pageNum, ptime, self.SecretKey)

	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["pageNum"] = pageNum
	params["numPerPage"] = pageSize
	params["format"] = "json"

	self.request(GET, url, str, params, &delvideolist)

	return &delvideolist
}

//获取视频播放的加密串
//return sign,ts
func (self PolyvInfo) GetVideoPlaySign(vid string) (sign string, ts int64) {
	ptime := time.Now().Unix() * 1000
	str_sign := fmt.Sprintf("%s%s%d", self.SecretKey, vid, ptime)
	return getMD5sign(str_sign), ptime
}

func getMD5sign(value string) string {
	return strings.ToUpper(hex.EncodeToString(sumMD5([]byte(value))))
}

func sumMD5(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

//移动视频到指定分类
func (self PolyvInfo) ChangeCata(vids, cataid string) *ChangeCataMsg {
	params = make(map[string]string)
	changcatamsg := ChangeCataMsg{}
	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/changeCata", self.UserID)
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("cataid=%s&ptime=%d&userid=%s&vids=%s%s", cataid, ptime, self.UserID, vids, self.SecretKey)

	params["format"] = "json"
	params["vids"] = vids
	params["cataid"] = cataid
	params["ptime"] = fmt.Sprintf("%d", ptime)

	self.request(GET, url, str, params, &changcatamsg)
	return &changcatamsg
}

//删除视频
func (self PolyvInfo) DelVideo(vid string) *DelVideoMsg {
	params = make(map[string]string)
	delmsg := DelVideoMsg{}

	url := fmt.Sprintf("http://api.polyv.net/v2/video/%s/del-video", self.UserID)
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("ptime=%d&vid=%s%s", ptime, vid, self.SecretKey)

	params["vid"] = vid
	params["ptime"] = fmt.Sprintf("%d", ptime)

	self.request(GET, url, str, params, &delmsg)

	return &delmsg
}

//查询视频播放量统计数据接口
func (self *PolyvInfo) VideoView(vid, dr, period string) *VideoViewMsg {
	params = make(map[string]string)
	videoviewmsg := VideoViewMsg{}

	ptime := time.Now().Unix() * 1000
	url := fmt.Sprintf("http://api.polyv.net/v2/videoview/%s", self.UserID)

	if dr == "" {
		dr = "7days"
	}

	if period == "" {
		period = "daily"
	}

	str := fmt.Sprintf("dr=%s&period=%s&ptime=%d&vid=%s%s", dr, period, ptime, vid, self.SecretKey)

	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["dr"] = dr
	params["period"] = period
	params["vid"] = vid

	self.request(GET, url, str, params, &videoviewmsg)

	return &videoviewmsg
}

//查询视频播放量排行接口
func (self *PolyvInfo) RankList(dr, start, end string) *RankMsg {
	params = make(map[string]string)
	rankmsg := RankMsg{}
	url := fmt.Sprintf("http://api.polyv.net/v2/videoview/%s/ranklist", self.UserID)
	ptime := time.Now().Unix() * 1000

	if dr == "" {
		dr = "7days"
	}

	str := fmt.Sprintf("dr=%s&end=%s&ptime=%d&start=%s%s", dr, end, ptime, start, self.SecretKey)

	params["start"] = start
	params["dr"] = dr
	params["end"] = end
	params["ptime"] = fmt.Sprintf("%d", ptime)
	self.request(GET, url, str, params, &rankmsg)
	return &rankmsg
}

//查询播放域名统计数据接口
func (self *PolyvInfo) DomainList(dr, start, end string) *DomainMsg {
	params = make(map[string]string)
	domainmsg := DomainMsg{}
	url := fmt.Sprintf("http://api.polyv.net/v2/domain/%s", self.UserID)
	ptime := time.Now().Unix() * 1000

	if dr == "" {
		dr = "7days"
	}

	str := fmt.Sprintf("dr=%s&end=%s&ptime=%d&start=%s%s", dr, end, ptime, start, self.SecretKey)
	params["start"] = start
	params["dr"] = dr
	params["end"] = end
	params["ptime"] = fmt.Sprintf("%d", ptime)

	self.request(GET, url, str, params, &domainmsg)
	return &domainmsg
}

//获取某一天视频日志
func (self *PolyvInfo) ViewLog(vid, day string) *VideoLogMsg {
	params = make(map[string]string)
	videologms := VideoLogMsg{}
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("day=%s&ptime=%d&userid=%s%s", day, ptime, self.UserID, self.SecretKey)
	url := fmt.Sprintf("http://api.polyv.net/v2/data/%s/viewlog", self.UserID)

	self.request(GET, url, str, params, &videologms)

	return &videologms
}

// 批量获取视频日志
func (self *PolyvInfo) MonthViewLog(month string, numPerPage, pageNum int) *VideoLogMsg {
	params = make(map[string]string)
	videologms := VideoLogMsg{}
	ptime := time.Now().Unix() * 1000
	url := fmt.Sprintf("http://api.polyv.net/v2/viewlog/%s/monthly/%s", self.UserID, month)
	str := fmt.Sprintf("month=%s&numPerPage=%d&pageNum=%d&ptime=%d%s", month, numPerPage, pageNum, ptime, self.SecretKey)
	if numPerPage == 0 {
		numPerPage = 99
	}

	if pageNum == 0 {
		pageNum = 1
	}

	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["month"] = month
	params["numPerPage"] = strconv.Itoa(numPerPage)
	params["pageNum"] = strconv.Itoa(pageNum)
	self.request(GET, url, str, params, &videologms)

	return &videologms
}
