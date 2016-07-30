package oss2bdc

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func toGbk(str string) string {
	bytes, _ := ioutil.ReadAll(transform.NewReader(bytes.NewReader([]byte(str)), simplifiedchinese.GBK.NewEncoder()))
	return string(bytes)
}

// Compress 压缩文件夹
func Compress(dir string) {
	config := GetConfig()
	parts := strings.Split(dir, "/")
	zipPath := config.RawPath + parts[1] + ".zip"
	zipFile, err := os.Create(zipPath)
	defer func() {
		zipFile.Close()
		newZipPath := config.ZipPath + parts[1] + ".zip"
		os.MkdirAll(config.ZipPath, 777)
		if err := os.Rename(zipPath, newZipPath); err != nil {
			log.Fatalln("os.Rename:", err)
		}
	}()
	if err != nil {
		log.Fatalln("os.Create:", err)
	}
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err := filepath.Walk(
		config.RawPath+dir,
		func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				log.Fatalln("os.Open:", err)
			}
			defer file.Close()
			parts := strings.Split(path, string(os.PathSeparator))
			writer, err := zipWriter.Create(toGbk(parts[3] + "/" + parts[4]))
			if err != nil {
				log.Fatalln("zipWriter.Create:", err)
			}

			bytes, err := ioutil.ReadAll(file)
			if _, err = writer.Write(bytes); err != nil {
				log.Fatal(err)
			}

			return nil
		}); err != nil {
		log.Fatalln("listDir:", err)
	}
}
