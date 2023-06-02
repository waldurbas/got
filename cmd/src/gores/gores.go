package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const resHead = `package resx

import "github.com/waldurbas/got/res"

// Res Data
var Res *res.Resdata

func init() {
	Res = res.New()`

const resFoot = `}`

var isDebug bool

func main() {

	for _, v := range os.Args[1:] {
		if strings.ToLower(v) == "-debug" {
			isDebug = true
			break
		}
	}

	headPrinted := false
	for _, v := range os.Args[1:] {
		ls := strings.ToLower(v)
		if ls == "-debug" {
			isDebug = true
			continue
		}

		if ls[:2] == "-d" {
			ix := strings.Index(v, "=")
			if ix > 0 {
				fname := v[ix+1:]
				fmt.Printf("\nFile: [%s]\n", fname)
				bData, err := ioutil.ReadFile(fname)
				if err != nil {
					log.Fatal(err)
				}

				//Add("gstock.eql", `H4sIAAA.... `)

				ix = bytes.IndexAny(bData, "`")
				if ix > -1 {
					bData = bData[ix+1:]
				}

				ix = bytes.IndexAny(bData, "`")
				if ix > -1 {
					bData = bData[:ix]
				}

				cData := string(bData)
				fmt.Printf("\ncoded=\"%v\"\n", cData)

				cData = strings.Trim(cData, "\n")

				dData, _ := base64.StdEncoding.DecodeString(cData)
				rdata := bytes.NewReader(dData)
				r, _ := gzip.NewReader(rdata)
				s, _ := ioutil.ReadAll(r)

				fmt.Printf("\ndecode=[%v]\n", string(s))

				fmt.Printf("\n\nlen GZIP=%d,len(s)=%d\n", len(cData), len(s))
			}

			return
		}

		if !headPrinted {
			fmt.Printf(resHead)
			headPrinted = true
		}

		cData := encodeFile(v)
		f := filepath.Base(v)

		fmt.Printf("\n    Res.Add(")

		a := 128
		le := len(cData)
		i := a - 12 - len(f)
		fmt.Printf("\"%s\", `%s", f, cData[:i])
		for i < le {
			n := i + a + 1
			if n >= le {
				n = le
			}
			fmt.Printf("\n%s", cData[i:n])
			i = n
		}
		fmt.Printf("`)\n")
	}

	if headPrinted {
		fmt.Printf(resFoot)
	}
}

func fileToTrimString(filename string) (dBytes []byte, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	for {
		bb, _, err := reader.ReadLine()

		if err == io.EOF {
			break
		}

		bb = bytes.TrimRight(bb, " ")
		bb = bytes.TrimLeft(bb, " \t")

		if len(bb) > 0 {
			if bytes.IndexAny(bb, "--") == 0 {
				continue
			}

			if bytes.IndexAny(bb, "#$&") == 0 {
				continue
			}

			bb = append(bb, 10)

			dBytes = append(dBytes, bb...)
		}
	}

	return
}

func encodeFile(filename string) (cData string) {

	rawbytes, err := fileToTrimString(filename)
	if err != nil {
		return
	}

	if isDebug {
		fmt.Printf("\ndebug[\n%v]\n", string(rawbytes))
	}

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	gz.Write(rawbytes)
	gz.Flush()
	gz.Close()

	cData = base64.StdEncoding.EncodeToString(b.Bytes())
	return
}
