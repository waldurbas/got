package htf

// ----------------------------------------------------------------------------------
// util.go (https://github.com/waldurbas/got/htf)
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
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"
	"time"
)

// DownloadFile #
func DownloadFile(url string, toFile string) error {
	// Create the file with .tmp extension, so that we won't overwrite a
	// file until it's downloaded fully
	tmpFile := toFile + ".tmp"
	os.Remove(tmpFile)

	out, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	// Create our bytes counter and pass it to be used alongside our writer
	counter := &WriteCounter{}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	// The progress use the same line so print a new line once it's finished downloading
	fmt.Println()
	out.Close()

	// Rename the tmp file back to the original file
	time.Sleep(1 * time.Second)
	err = os.Rename(tmpFile, toFile)
	if err != nil {
		return err
	}

	return nil
}

var wg sync.WaitGroup

// ParalellDownloadFile #
func ParalellDownloadFile(url string, toFile string) error {
	flen, err := URLfileSize(url)
	if err != nil {
		return err
	}

	// in 256 kB Bloecken
	maxParts := flen / (1024 * 256)
	if maxParts > 10 {
		maxParts = 10
	}

	lenPart := flen / maxParts         // Bytes for each Go-routine
	diff := flen % maxParts            // Get the remaining for the last request
	body := make([]string, maxParts+1) // Make up a temporary array to hold the data to be written to the file

	fmt.Print(" ")
	for i := 0; i < maxParts; i++ {
		wg.Add(1)

		min := lenPart * i       // Min range
		max := lenPart * (i + 1) // Max range

		if i == maxParts-1 {
			max += diff // Add the remaining bytes in the last request
		}

		go func(min int, max int, i int) {
			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)

			// Add the data for the Range header of the form "bytes=0-100"
			rangeHeader := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
			req.Header.Add("Range", rangeHeader)

			r, _ := client.Do(req)
			defer func() {
				r.Body.Close()
			}()

			b, _ := ioutil.ReadAll(r.Body)
			body[i] = string(b)
			wg.Done()
			fmt.Print(".")
		}(min, max, i)
	}
	wg.Wait()

	outFile := toFile + ".tmp"
	os.Remove(outFile)

	f, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	for i := 0; i < maxParts; i++ {
		f.Write([]byte(body[i]))
	}
	f.Close()

	os.Remove(toFile)
	err = os.Rename(outFile, toFile)
	if err != nil {
		return err
	}

	return nil
}

// URLfileSize #
func URLfileSize(url string) (int, error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}

	// Is our request ok?
	if resp.StatusCode != http.StatusOK {
		err := errors.New(resp.Status)
		return 0, err
	}

	// the Header "Content-Length" will let us know
	// the total file size to download
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))

	return size, nil
}

func readableBytes(n uint64) string {
	sizes := []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

	b := float64(1024)
	e := math.Floor(math.Log(float64(n)) / math.Log(b))
	sfx := sizes[int(e)]
	v := float64(n) / math.Pow(b, math.Floor(e))
	f := "%.0f"
	if v < 10 {
		f = "%.1f"
	}

	return fmt.Sprintf(f+" %s", v, sfx)
}

func urlDownloadListFile(url string) (string, error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Chmod #
func Chmod(name string, m uint32) error {
	return os.Chmod(name, fs.FileMode(m))
}

// RemoveOldFile #
func RemoveOldFile(xFile string) string {
	ext := path.Ext(xFile)
	oFile := xFile[0:len(xFile)-len(ext)] + ".old"
	if FileExists(oFile) {
		os.Remove(oFile)
	}

	return oFile
}

// FileExists #
func FileExists(filename string) bool {
	f, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !f.IsDir()
}
