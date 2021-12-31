package packer

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"io/fs"
	"io/ioutil"
)

func PackBytes(json []byte) []byte {
	var err error
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	var zipFile io.Writer
	zipFile, err = zipWriter.Create("data")
	if err == nil {
		_, err = zipFile.Write(json)
	}
	err = zipWriter.Close()
	result := []byte(base64.StdEncoding.EncodeToString(buf.Bytes()))
	//logger.Println("Packer stats: ", len(json), len(result), float64(len(json))/float64(len(result)))
	return result
}

func UnpackString(b64 string) string {
	var err error
	var data []byte
	data, err = base64.StdEncoding.DecodeString(b64)
	if err == nil {
		buf := bytes.NewReader(data)
		var zipFile *zip.Reader
		zipFile, err = zip.NewReader(buf, buf.Size())
		var file fs.File
		file, err = zipFile.Open("data")
		if err == nil {
			var bs []byte
			bs, err = ioutil.ReadAll(file)
			if err == nil {
				return string(bs)
			}
			_ = file.Close()
		}
	}
	return ""
}
