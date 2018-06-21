package sdk

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/smallnest/goreq"
	"io"
	"strconv"
	"strings"
	"time"
)

var params map[string]string

const (
	DefaultAPIHost   = "http://api.polyv.net/v2"
	DefaultVideoHost = "http://api.polyv.net/v2/video"
	DefaultService   = "http://v.polyv.net/uc/services/rest?method="
)

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
	//http://v.polyv.net/uc/services/rest?method=getHostUrl
	var hostmsg HostMsg
	url := fmt.Sprintf("%sgetHostUrl", DefaultService)
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

// 上传视频的预览图
func (self *PolyvInfo) UpFirstImageByUrl(vid, img_url, img_name string) *UploadImgMsg {
	respmsg := self.UploadConverImageUrl(vid, "", img_url)
	if respmsg.Status_Code == 200 {
		return &UploadImgMsg{
			ErrorStr: "",
			Data:     true,
		}
	} else {
		return &UploadImgMsg{
			ErrorStr: respmsg.Message,
			Data:     false,
		}
	}
}

//上传多个视频的预览图URL
func (self *PolyvInfo) UploadConverImageUrl(vids, cataids, img_url string) *RespMsg {
	//http://api.polyv.net/v2/video/{userid}/uploadCoverImageUrl
	params = make(map[string]string)
	respmsg := RespMsg{}
	str := ""

	if img_url == "" {
		respmsg.Status_Code = 400
		respmsg.Message = "图片路径不能为空"
		return &respmsg
	}

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

	url := fmt.Sprintf("%s/%s/uploadCoverImageUrl", DefaultVideoHost, self.UserID)
	self.request(POST, url, str, params, &respmsg)
	return &respmsg

}

// 上传远程视频
func (self *PolyvInfo) UploadUrlFile(urlfileinfo *UrlFileInfo) *ReturnMsg {
	// http://v.polyv.net/uc/services/rest?method=uploadUrlFile
	var upvideomsg ReturnMsg
	var upasyncmsg UploadAsyncVideoMsg
	var err error
	if urlfileinfo.FileUrl == "" || urlfileinfo.Title == "" {
		return &ReturnMsg{Status_Code: 400}
	}
	//必须这样写，不要调整顺序！！！
	str := fmt.Sprintf("desc=%s&fileUrl=%s&tag=%s&title=%s&writetoken=%s%s",
		urlfileinfo.Desc, urlfileinfo.FileUrl, urlfileinfo.Tag, urlfileinfo.Title, self.WriteToken, self.SecretKey)

	sign := GetSign(str)

	url := fmt.Sprintf("%suploadUrlFile", DefaultService)
	req := goreq.New().Get(url)

	if urlfileinfo.CataId != "" {
		req.Query("cataid=" + urlfileinfo.CataId)
	}

	if urlfileinfo.Async {
		req.Query("async=true")
	}

	_, body, errs := req.
		Query("writetoken=" + self.WriteToken).
		Query("fileUrl=" + urlfileinfo.FileUrl).
		Query("title=" + urlfileinfo.Title).
		Query("desc=" + urlfileinfo.Desc).
		Query("tag=" + urlfileinfo.Tag).
		Query("cataid=" + urlfileinfo.CataId).
		Query("sign=" + sign).
		SetDebug(self.Verbose).
		SetCurlCommand(self.Verbose).
		End()

	if len(errs) > 0 {
		return &ReturnMsg{
			Status_Code: 400,
		}
	}

	if urlfileinfo.Async {
		err = json.Unmarshal([]byte(body), &upasyncmsg)
	} else {
		err = json.Unmarshal([]byte(body), &upvideomsg)
	}

	if err != nil {
		return &ReturnMsg{Status_Code: 400}
	}

	if urlfileinfo.Async {
		scode, _ := strconv.Atoi(upasyncmsg.Status_Code)
		upvideomsg = ReturnMsg{
			Status_Code: scode,
		}
	}
	return &upvideomsg
}

//获取用户空间及流量情况
func (self *PolyvInfo) GetUseInfo(query_date string) *UseMsg {
	//http://api.polyv.net/v2/user/{userid}/main
	var usemsg UseMsg
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("date=%s&ptime=%d%s", query_date, ptime, self.SecretKey)

	sign := GetSign(str)

	_, _, errs := goreq.New().Get(fmt.Sprintf("%s/user/%s/main", DefaultAPIHost, self.UserID)).
		BindBody(&usemsg).
		Query("ptime=" + fmt.Sprintf("%d", ptime)).
		Query("sign=" + sign).
		Query("date=" + query_date).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &UseMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &usemsg
}

//获取用户空间及流量情况
func (self *PolyvInfo) GetTotalUseInfo() *UseMsg {
	//http://api.polyv.net/v2/user/{userid}/main
	var usemsg UseMsg
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("ptime=%d%s", ptime, self.SecretKey)

	sign := GetSign(str)
	_, _, errs := goreq.New().Get(fmt.Sprintf("%s/user/%s/main", DefaultAPIHost, self.UserID)).
		BindBody(&usemsg).
		Query("ptime=" + fmt.Sprintf("%d", ptime)).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &UseMsg{
			RespMsg: RespMsg{

				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &usemsg
}

//获取单个视频的首图
func (self *PolyvInfo) GetVideoImage(vid, t string) *VideoImgMsg {
	//http://api.polyv.net/v2/video/{userid}/get-image
	var vimgmsg VideoImgMsg
	ptime := time.Now().Unix() * 1000

	if t != "1" {
		t = "2"
	}

	str := fmt.Sprintf("ptime=%d&t=%s&vid=%s%s", ptime, t, vid, self.SecretKey)
	sign := GetSign(str)

	url := fmt.Sprintf("%s/video/%s/get-image", DefaultAPIHost, self.UserID)

	_, _, errs := goreq.New().Get(url).
		BindBody(&vimgmsg).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("t=" + t).
		Query("vid=" + vid).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &VideoImgMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &vimgmsg
}

const (
	GET  = "get"
	POST = "post"
)

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
	//http://api.polyv.net/v2/video/{userid}/get-video-msg
	params = make(map[string]string)
	videomsg := VideoMsg{}

	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("format=%s&ptime=%d&vid=%s%s", "json", ptime, vid, self.SecretKey)

	params["format"] = "json"
	params["ptime"] = fmt.Sprintf("%d", ptime)
	params["vid"] = vid
	url := fmt.Sprintf("%s/video/%s/get-video-msg", DefaultAPIHost, self.UserID)

	self.request(GET, url, str, params, &videomsg)

	return &videomsg
}

//获取最新视频/全部视频列表
func (self *PolyvInfo) GetVideoList(catatree, pageSize, pageNum, startDate, endDate string) *VideoList {
	//http://api.polyv.net/v2/video/{userid}/get-new-list
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

	url := fmt.Sprintf("%s/video/%s/get-new-list?%s", DefaultAPIHost, self.UserID, param_str)
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

	url := fmt.Sprintf("%s/video/%s/search?%s", DefaultAPIHost, self.UserID, param_str)
	self.request(GET, url, str, nil, &videolist)

	return &videolist
}

func (self *PolyvInfo) AddCata(cata_name string) *AddCataMsg {
	resp := AddCataMsg{}
	params = make(map[string]string)
	url := fmt.Sprintf("%s/%s/addCata", DefaultVideoHost, self.UserID)
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
	url := fmt.Sprintf("%s/%s/deleteCata", DefaultVideoHost, self.UserID)
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
	url := fmt.Sprintf("%s/%s/cataJson", DefaultVideoHost, self.UserID)
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

	url := fmt.Sprintf("%s/%s/get-del-list", DefaultVideoHost, self.UserID)
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
	// http://api.polyv.net/v2/video/{userid}/changeCata
	var changcatamsg ChangeCataMsg
	url := fmt.Sprintf("%s/%s/changeCata", DefaultVideoHost, self.UserID)
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("cataid=%s&ptime=%d&userid=%s&vids=%s%s", cataid, ptime, self.UserID, vids, self.SecretKey)
	sign := GetSign(str)

	_, _, errs := goreq.New().Get(url).
		BindBody(&changcatamsg).
		Query("format=json").
		Query("vids=" + vids).
		Query("cataid=" + cataid).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &ChangeCataMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &changcatamsg
}

//删除视频
func (self PolyvInfo) DelVideo(vid string) *DelVideoMsg {
	//http://api.polyv.net/v2/video/{userid}/del-video
	var delmsg DelVideoMsg
	url := fmt.Sprintf("%s/%s/del-video", DefaultVideoHost, self.UserID)
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("ptime=%d&vid=%s%s", ptime, vid, self.SecretKey)
	sign := GetSign(str)

	_, _, errs := goreq.New().Get(url).BindBody(&delmsg).
		Query("vid=" + vid).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &DelVideoMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &delmsg
}

//查询视频播放量统计数据接口
func (self *PolyvInfo) VideoView(vid, dr, period string) *VideoViewMsg {
	var videoviewmsg VideoViewMsg
	req := goreq.New().Get(fmt.Sprintf("%s/videoview/%s", DefaultAPIHost, self.UserID))
	ptime := time.Now().Unix() * 1000

	if dr == "" {
		dr = "7days"
	}

	if period == "" {
		period = "daily"
	}

	if vid == "" {
		return &VideoViewMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     "VID不能为空",
			},
		}
	}

	str := fmt.Sprintf("dr=%s&period=%s&ptime=%d&vid=%s%s", dr, period, ptime, vid, self.SecretKey)

	sign := GetSign(str)

	req.Query("vid=" + vid).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("dr=" + dr).
		Query("period=" + period).
		Query("sign=" + sign)

	_, _, errs := req.
		BindBody(&videoviewmsg).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &VideoViewMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &videoviewmsg
}

//查询视频播放量排行接口
func (self *PolyvInfo) RankList(dr, start, end string) *RankMsg {
	var rankmsg RankMsg
	req := goreq.New().Get(fmt.Sprintf("%s/videoview/%s/ranklist", DefaultAPIHost, self.UserID))
	jsonp := "" //todo
	ptime := time.Now().Unix() * 1000

	if dr == "" {
		dr = "7days"
	}

	str := ""
	if jsonp == "" {
		str = fmt.Sprintf("dr=%s&end=%s&ptime=%d&start=%s%s", dr, end, ptime, start, self.SecretKey)
	} else {
		str = fmt.Sprintf("dr=%s&jsonp=%s&end=%s&ptime=%d&start=%s%s", dr, jsonp, end, ptime, start, self.SecretKey)
	}

	sign := GetSign(str)

	_, _, errs := req.
		BindBody(&rankmsg).
		Query("start=" + start).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("dr=" + dr).
		Query("end=" + end).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &RankMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &rankmsg
}

//查询播放域名统计数据接口
func (self *PolyvInfo) DomainList(dr, start, end string) *DomainMsg {
	var domainmsg DomainMsg
	req := goreq.New().Get(fmt.Sprintf("%s/domain/%s", DefaultAPIHost, self.UserID))
	jsonp := "" //todo
	ptime := time.Now().Unix() * 1000

	if dr == "" {
		dr = "7days"
	}

	str := ""
	if jsonp == "" {
		str = fmt.Sprintf("dr=%s&end=%s&ptime=%d&start=%s%s", dr, end, ptime, start, self.SecretKey)
	} else {
		str = fmt.Sprintf("dr=%s&jsonp=%s&end=%s&ptime=%d&start=%s%s", dr, jsonp, end, ptime, start, self.SecretKey)
	}

	sign := GetSign(str)

	_, _, errs := req.
		BindBody(&domainmsg).
		Query("start=" + start).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("dr=" + dr).
		Query("end=" + end).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &DomainMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &domainmsg
}

//获取某一天视频日志
func (self *PolyvInfo) ViewLog(vid, day string) *VideoLogMsg {
	var videologms VideoLogMsg
	req := goreq.New().Get(fmt.Sprintf("%s/data/%s/viewlog", DefaultAPIHost, self.UserID))
	ptime := time.Now().Unix() * 1000

	str := fmt.Sprintf("day=%s&ptime=%d&userid=%s%s", day, ptime, self.UserID, self.SecretKey)

	sign := GetSign(str)

	_, _, errs := req.
		BindBody(&videologms).
		Query("day=" + day).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("vid=" + vid).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &VideoLogMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &videologms
}

// 批量获取视频日志
func (self *PolyvInfo) MonthViewLog(month string, numPerPage, pageNum int) *VideoLogMsg {
	var videologms VideoLogMsg
	req := goreq.New().Get(fmt.Sprintf("%s/viewlog/%s/monthly/%s", DefaultAPIHost, self.UserID, month))
	ptime := time.Now().Unix() * 1000
	str := fmt.Sprintf("month=%s&numPerPage=%d&pageNum=%d&ptime=%d%s", month, numPerPage, pageNum, ptime, self.SecretKey)

	sign := GetSign(str)
	if numPerPage == 0 {
		numPerPage = 99
	}

	if pageNum == 0 {
		pageNum = 1
	}

	_, _, errs := req.
		BindBody(&videologms).
		Query("month=" + month).
		Query(fmt.Sprintf("ptime=%d", ptime)).
		Query("numPerPage=" + strconv.Itoa(numPerPage)).
		Query("pageNum=" + strconv.Itoa(pageNum)).
		Query("sign=" + sign).
		SetCurlCommand(self.Verbose).
		SetDebug(self.Verbose).
		End()

	if len(errs) > 0 {
		return &VideoLogMsg{
			RespMsg: RespMsg{
				Status_Code: 400,
				Status:      "error",
				Message:     errs[0].Error(),
			},
		}
	}

	return &videologms
}
