package htf

import (
	"os"
)

func GetExecutable(url string, dir string) error {
	xFile, _ := os.Executable()

	RemoveOldFile(xFile)

	f, err := GetDownloadFilesInfo(url + "/" + dir)

	if err != nil {
		return err
	}

	fi, err := f.GetFileInfo(xFile)
	if err != nil {
		return err
	}

	if !fi.Changed {
		return nil
	}

	oFile := RemoveOldFile(xFile)

	err = os.Rename(xFile, oFile)
	if err != nil {
		return err
	}

	err = fi.Download(xFile)
	if err != nil {
		return err
	}

	Chmod(xFile, 0755)
	if FileExists(oFile) {
		os.Remove(oFile)
	}

	return nil
}
