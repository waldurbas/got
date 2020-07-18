package htp_test

import (
	"fmt"
	"testing"

	"github.com/waldurbas/got/cnv"
	h "github.com/waldurbas/got/htp"
)

func Test_checkDownLoad(t *testing.T) {
	url := "http://xxx..."
	di, err := h.GetDownloadFilesInfo(url)

	if err != nil {
		t.Errorf("GetDownloadFilesInfo(%s):\n%v", url, err)
	} else {
		s := ""
		for _, fi := range di.List {
			//			Size uint64
			//			Time time.Time
			ts := fi.Web.Time.Format("2006-01-02 15:04:05")
			s = s + "\n" + fmt.Sprintf("%-20s %12s  %s", fi.FileName, cnv.FormatInt64(int64(fi.Web.Size)), ts)

		}
		t.Errorf("noError: GetDownloadFilesInfo(%s):%s", url, s)
	}
}
