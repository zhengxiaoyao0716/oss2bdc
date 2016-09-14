package oss2bdc

import (
	"log"
	"os"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var bucket *oss.Bucket

func init() {
	config := GetConfig()

	client, err := oss.New(config.Oss.Endpoint, config.Oss.Key, config.Oss.Secret)
	if err != nil {
		log.Fatalln("oss.New: ", err)
	}
	ossBucket, err := client.Bucket(config.Oss.BucketName)
	if err != nil {
		log.Fatalln("client.Bucket: ", err)
	}
	bucket = ossBucket
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func replaceSpecialChar(rawStr string) string {
	var str string
	str = strings.Replace(rawStr, ":", "：", -1)
	str = strings.Replace(str, "*", "　", -1)
	str = strings.Replace(str, "?", "？", -1)
	str = strings.Replace(str, "\"", "”", -1)
	str = strings.Replace(str, "<", "《", -1)
	str = strings.Replace(str, ">", "》", -1)
	str = strings.Replace(str, "|", "　", -1)
	return str
}

func downloadObject(name string) {
	var dir, objPath string
	if engine == nil {
		objPath = GetConfig().RawPath + replaceSpecialChar(name)
		dir = objPath[0:strings.LastIndex(objPath, "/")]
		os.MkdirAll(dir, 777)
	} else {
		parts := strings.Split(toGbk(name), "/")
		team := replaceSpecialChar(Phone2Team(parts[2]))
		dir = GetConfig().RawPath + parts[0] + "/" + parts[1] + "/" + team + "/"
		objPath = dir + parts[3] + ".jpg"
		if exist(objPath) {
			return
		}
	}
	os.MkdirAll(dir, 777)
	if err := bucket.DownloadFile(name, objPath, 1024*1024, oss.Routines(3), oss.Checkpoint(true, "")); err != nil {
		switch err.(type) {
		case oss.ServiceError:
			log.Println("bucket.DownloadFile: ", err.(oss.ServiceError))
		case error:
			log.Println("bucket.DownloadFile: ", err)
		}
	}
}

// Download 下载
func Download(dir string) {
	prefix := oss.Prefix(dir)
	marker := oss.Marker("")
	for {
		lsRes, err := bucket.ListObjects(oss.MaxKeys(1000), marker, prefix)
		if err != nil {
			log.Fatalln("bucket.ListObjects:", err)
		}

		prefix = oss.Prefix(lsRes.Prefix)
		marker = oss.Marker(lsRes.NextMarker)

		for _, properties := range lsRes.Objects {
			// log.Println("Downloading: ", properties.Key, properties.Size)
			downloadObject(properties.Key)
		}

		if !lsRes.IsTruncated {
			break
		}
	}
}
