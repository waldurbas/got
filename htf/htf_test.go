package htf_test

import (
	"fmt"
	"testing"

	"github.com/waldurbas/got/cnv"
	"github.com/waldurbas/got/htf"
)

func Test_checkDownLoad(t *testing.T) {
	url := "https://xxx..."

	di, err := htf.GetDownloadFilesInfo(url)

	if err != nil {
		t.Errorf("GetDownloadFilesInfo(%s):\n%v", url, err)
	} else {
		s := ""
		for i, fi := range di.List {
			ts := fi.Web.Time.Format("2006-01-02 15:04:05")
			s = s + "\n" + fmt.Sprintf("%3d. %-20s %12s  %s", i+1, fi.FileName, cnv.FormatInt64(int64(fi.Web.Size)), ts)

		}
		fmt.Println(s)
	}
}
