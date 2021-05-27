package htf

// ----------------------------------------------------------------------------------
// getFile.go (https://github.com/waldurbas/got/htf)
// Copyright 2018,2021 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2018.12.11 (wu) Init
//-----------------------------------------------------------------------------------

import (
	"os"
	"path"
)

func GetFile(url string, dir string, xFile string, locFile string) (bool, error) {
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

	// download to nFile
	err = fi.Download(locFile)
	if err != nil {
		return false, err
	}

	return true, nil
}

// GetExecutableFile #
func GetExecutableFile(url string, dir string, xFile string) (bool, error) {
	ext := path.Ext(xFile)
	oFile := xFile[0:len(xFile)-len(ext)] + ".old"
	nFile := xFile[0:len(xFile)-len(ext)] + ".new"

	RemoveFile(oFile)
	RemoveFile(nFile)

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

	// download to nFile
	err = fi.Download(nFile)
	if err != nil {
		return false, err
	}

	// if OK: xFile to oFile
	if FileExists(xFile) {
		err = os.Rename(xFile, oFile)
		if err != nil {
			return false, err
		}
	}

	// if OK: newFile to xFile
	if FileExists(nFile) {
		err = os.Rename(nFile, xFile)
		if err != nil {
			// back to xFile
			os.Rename(oFile, xFile)
			return false, err
		}
	}

	if FileExists(xFile) {
		Chmod(xFile, 0755)
		RemoveFile(oFile)
	}

	return true, nil
}

// GetExecutable #
func GetExecutable(url string, dir string) (bool, error) {
	xFile, _ := os.Executable()

	return GetExecutableFile(url, dir, xFile)
}
