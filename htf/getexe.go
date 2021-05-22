package htf

import (
	"os"
)

func GetExecutable(url string, dir string) (bool, error) {
	xFile, _ := os.Executable()

	RemoveOldFile(xFile)

	f, err := GetDownloadFilesInfo(url + "/" + dir)

	if err != nil {
		return false, err
	}

	fi, err := f.GetFileInfo(xFile)
	if err != nil {
		return false, err
	}

	if !fi.Changed {
		return false, nil
	}

	oFile := RemoveOldFile(xFile)

	err = os.Rename(xFile, oFile)
	if err != nil {
		return false, err
	}

	err = fi.Download(xFile)
	if err != nil {
		return false, err
	}

	if FileExists(oFile) {
		os.Remove(oFile)
	}

	return true, nil
}
