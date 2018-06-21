package sdk

import (
	"github.com/jinzhu/now"
	. "github.com/smartystreets/goconvey/convey"
	"strconv"
	"strings"
	"testing"
)

func TestVideoInfo(t *testing.T) {

	p := PolyvInfo{
		UserID:     "",
		WriteToken: "",
		ReadToken:  "",
		SecretKey:  "",
		Verbose:    false,
	}

	Convey("test polyv-go-sdk video", t, func() {
		msg := &CataMsg{}
		msg = p.CataJson()
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
		vinfo := p.GetVideoInfo(vid)
		So(vinfo.Status_Code, ShouldEqual, 200)
		So(len(vinfo.Data), ShouldEqual, 1)

		img_url := ""
		up_msg := p.UploadConverImageUrl(vid, "", img_url)
		So(up_msg.Status_Code, ShouldEqual, 200)

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
	})

}
