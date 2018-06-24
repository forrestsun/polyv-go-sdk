package sdk

import (
	"github.com/jinzhu/now"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

func TestVideoInfo(t *testing.T) {
	var p PolyvInfo
	email := ""
	passwd := ""
	passwdmd5 := ""
	global_vid := ""

	Convey("test login", t, func() {
		pu := Login(email, passwd, passwdmd5, false)
		p = PolyvInfo{
			UserID:     pu.UserID,
			ReadToken:  pu.ReadToken,
			WriteToken: pu.WriteToken,
			SecretKey:  pu.SecretKey,
			Verbose:    true,
		}
		So(pu.Status_Code, ShouldEqual, 200)

		// update_cata_msg := p.UpdateCata("1529561645408", "y", "", false, false, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg = p.UpdateCata("1529561645408", "n", "", false, false, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg = p.UpdateCata("1529561645408", "y", "", false, false, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg = p.UpdateCata("1529561645408", "ccc", "", false, false, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg = p.UpdateCata("1529561645408", "y", "", true, true, true)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg = p.UpdateCata("1529561645408", "y", "", true, true, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg = p.UpdateCata("1529561645408", "y", "", false, true, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

		// update_cata_msg := p.UpdateCata("1529561645408", "y", "", false, false, false)
		// So(update_cata_msg.Status_Code, ShouldEqual, 200)

	})

	Convey("test polyv-go-sdk video", t, func() {
		use_msg := p.GetUseInfo(now.BeginningOfDay().Format("2006-01-02"))
		So(use_msg.Status_Code, ShouldEqual, 200)
		So(use_msg.Data.UserId, ShouldEqual, p.UserID)
		So(use_msg.Data.Email, ShouldEqual, email)

		msg := p.CataJson()
		So(msg.Status_Code, ShouldEqual, 200)
		org_cata_nodes := len(msg.Data[0].CataNodes)

		cata_msg := p.AddCata("a")
		So(cata_msg.Status_Code, ShouldEqual, 200)
		So(len(strings.Split(cata_msg.Data.CataTree, ",")), ShouldEqual, 2)

		msg = p.CataJson()
		So(msg.Status_Code, ShouldEqual, 200)
		So(len(msg.Data[0].CataNodes), ShouldEqual, org_cata_nodes+1)

		for _, v := range msg.Data[0].CataNodes {
			if v.Cataname == "a" {
				delmsg := p.DelCata(strconv.FormatUint(v.Cataid, 10))
				So(delmsg.Status_Code, ShouldEqual, 200)
			}
		}

		msg = &CataMsg{}
		msg = p.CataJson()
		So(msg.Status_Code, ShouldEqual, 200)
		So(len(msg.Data[0].CataNodes), ShouldEqual, org_cata_nodes)

		total_record := msg.Data[0].Videos

		vlist := p.GetVideoList("", "3", "1", "", "")
		So(vlist.Status_Code, ShouldEqual, 200)
		So(len(vlist.Data), ShouldEqual, 3)
		So(vlist.Total, ShouldEqual, total_record)

		vid := vlist.Data[0].Vid
		global_vid = vid
		vinfo := p.GetVideoInfo(vid)
		So(vinfo.Status_Code, ShouldEqual, 200)
		So(len(vinfo.Data), ShouldEqual, 1)

		img_url := "http://ok0jpejfs.bkt.clouddn.com/k2.jpg"
		up_msg := p.UploadConverImageUrl(vid, "", img_url)
		So(up_msg.Status_Code, ShouldEqual, 200)

		img_url = ""
		up_msg = p.UploadConverImageUrl(vid, "", img_url)
		So(up_msg.Status_Code, ShouldEqual, 400)

		img_msg := p.GetVideoImage(vid, "2")
		So(img_msg.Message, ShouldEqual, "success")
		So(img_msg.Status_Code, ShouldEqual, 200)
		So(img_msg.Status, ShouldEqual, "success")

		title := vlist.Data[0].Title
		for i := 0; i < 10; i++ {
			mlist := p.SearchByTitle(title, "100", "1")
			So(mlist.Status_Code, ShouldEqual, 200)
			So(mlist.Total, ShouldNotEqual, 0)
		}

		start_date := now.MustParse(vlist.Data[0].PTime)
		end_date := now.MustParse(vlist.Data[0].PTime)
		vlist = p.GetVideoList("", "9", "1", start_date.Format("2006-01-02"), end_date.Format("2006-01-02"))
		So(vlist.Status_Code, ShouldEqual, 200)
		So(vlist.Total, ShouldEqual, 0)

		vlist = p.GetVideoList("", "9", "1", start_date.Format("2006-01-02"), "")
		So(vlist.Status_Code, ShouldEqual, 200)
		So(vlist.Total, ShouldNotEqual, 0)

		vlist = p.GetVideoList("", "9", "1", "", end_date.Format("2006-01-02"))
		So(vlist.Status_Code, ShouldEqual, 200)
		So(vlist.Total, ShouldNotEqual, 0)

		for _, v := range msg.Data[0].CataNodes {
			vlist = p.GetVideoList(v.Catatree, "10", "1", "", "")
			So(vlist.Status_Code, ShouldNotEqual, 400)
			So(vlist.Total, ShouldEqual, v.Videos)
		}

		vid = "1"
		vinfo = p.GetVideoInfo(vid)
		So(vinfo.Status_Code, ShouldEqual, 400)

		//todo:上传一个资源
		//查回收站
		//删除上传资源
		//查回收站

		delmsg := p.GetDelList("1", "10")
		So(delmsg.Status_Code, ShouldEqual, 200)
		So(len(delmsg.Data), ShouldEqual, 0)

	})

	Convey("get videoview", t, func() {
		vs := p.VideoView("", "", "")
		So(vs.Status_Code, ShouldEqual, 400)

		dr := "" //默认值为7days
		//时间段，具体值为以下几个：today（今天），yesterday（昨天），this_week（本周），
		//last_week（上周），7days（最近7天），this_month（本月），last_month（上个月），
		//this_year（今年），last_year（去年），默认值为7days

		period := ""
		//显示周期，具体为以下几个值：daily（按日显示），weekly（按周显示），monthly（按月显示）。
		//默认值为daily。period的值受限于dr的值，当dr的值为today，yesterday，this_week，
		//last_week，7days时，period只能为daily，当dr的值为this_month，last_month时，
		//period只能为daily或者weekly
		vs = p.VideoView(global_vid, dr, period)
		So(vs.Status_Code, ShouldEqual, 200)
		So(len(vs.Data), ShouldEqual, 7)

		dr = "this_month"
		period = "daily"
		vs = p.VideoView(global_vid, dr, period)
		So(vs.Status_Code, ShouldEqual, 200)
		// So(len(vs.Data), ShouldEqual, 24)

	})

	Convey("get videoview ranklist", t, func() {
		vs := p.RankList("", now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02"))
		So(vs.Status_Code, ShouldEqual, 200)
	})

	Convey("get videoview domain list", t, func() {
		vs := p.DomainList("", now.BeginningOfMonth().Format("2006-01-02"), now.EndOfMonth().Format("2006-01-02"))
		So(vs.Status_Code, ShouldEqual, 200)

		vs = p.DomainList("this_year", now.BeginningOfDay().Format("2006-01-02"), now.BeginningOfDay().Format("2006-01-02"))
		So(vs.Status_Code, ShouldEqual, 200)
	})

	Convey("get videoview log list by month", t, func() {
		vs := p.MonthViewLog(now.BeginningOfMonth().Format("200601"), 50, 1)
		So(vs.Status_Code, ShouldEqual, 200)
	})

}
