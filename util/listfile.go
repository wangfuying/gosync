package util

import (
	"io/ioutil"
	"os"
	"time"
)

type FileInfo struct {
	Name     string
	Size     int64
	Mode     os.FileMode
	Time     time.Time
	Dir      bool
	FullPath string
	Flag     string
}

func ListFiles(dir string) (ret []*FileInfo, err error) {
	files, err := ioutil.ReadDir(dir)
	if nil != err {
		return nil, err
	}
	datas := make([]*FileInfo, 0)
	for _, file := range files {
		fileinfo := &FileInfo{
			file.Name(),
			file.Size(),
			file.Mode(),
			file.ModTime(),
			file.IsDir(),
			"",
			"",
		}
		datas = append(datas, fileinfo)
	}
	return datas, err
}

func ScanFile(path string, suffix string, handler interface{}) int {
	files, err := ioutil.ReadDir(path)
	if nil != err {
		Error(err.Error())
		return 0
	}
	if len(suffix) == 0 {
		suffix = ""
	}
	for _, file := range files {
		if file.IsDir() {
			// go in
			newSuffix := suffix + file.Name() + string(os.PathSeparator)
			ScanFile(path+string(os.PathSeparator)+file.Name(), newSuffix, handler)
		} else {
			handler.(func(os.FileInfo, string))(file, suffix)
		}
	}
	return 1
}
