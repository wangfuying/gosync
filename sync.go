package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"gosync/util"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mode string

var syncDir string = "." + string(os.PathSeparator) + "Files" + string(os.PathSeparator)

// server
var serverPort int = 666

// client
var scanPeriod time.Duration = 1000 * 1000 * 1000 // 1 second
var ipAddr string
var datas []*util.FileInfo

var workResultLock sync.WaitGroup

func main() {
	flag.StringVar(&mode, "m", "client", "How a client or a server you want to create ?")
	flag.StringVar(&syncDir, "d", syncDir, "What dirctory would you want to sync ?")
	flag.StringVar(&ipAddr, "a", "127.0.0.1:"+strconv.Itoa(serverPort), "What dirctory would you want to sync ?")
	flag.Parse()
	fmt.Println("Go File SyncEr v1.0")

	if !strings.Contains(syncDir, string(os.PathSeparator)) {
		syncDir = syncDir + string(os.PathSeparator)
	}
	if strings.ToLower(mode) == "client" {
		util.Info("Run as CLIENT ")
		workResultLock.Add(1)
		go Client()
	} else {
		util.Info("Run as SERVER with port %d", serverPort)
		workResultLock.Add(1)
		go Server()
	}

	// wait
	workResultLock.Wait()

}

func Client() {
	// init
	flag, err := util.PathExists(syncDir)
	if nil != err {
		return
	}
	if !flag {
		os.Mkdir(syncDir, os.ModePerm)
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ipAddr) //获取一个TCP地址信息,TCPAddr
	if nil != err {
		util.Error("Can not parse Server IP Address !!!")
		goto END
	}

	for {
		time.Sleep(scanPeriod)
		util.Info("Scanning Files ...")
		// datas, err := util.ListFiles(syncDir)
		datas = make([]*util.FileInfo, 0)
		util.ScanFile(syncDir, "."+string(os.PathSeparator), EnumFileHandler)
		jsonStr, err := json.Marshal(datas)
		fmt.Println(string(jsonStr))
		conn, err := net.DialTCP("tcp", nil, tcpAddr) //创建一个TCP连接:TCPConn
		if nil != err {
			util.Error(err.Error())
			continue
		}
		defer conn.Close()
		len, err := conn.Write(jsonStr)
		if nil != err {
			util.Error(err.Error())
		}
		util.Trace("Successed Write to Server %d bytes", len)
		conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		result, err := util.ReadData(conn) //获得返回数据
		if nil != err {
			util.Error(err.Error())
		}
		fmt.Println(string(result))
		/**
		**/
		time.Sleep(scanPeriod * 10)
	}
END:
	workResultLock.Done()
}

func Server() {
	// init
	flag, err := util.PathExists(syncDir)
	if nil != err {
		return
	}
	if !flag {
		os.Mkdir(syncDir, os.ModePerm)
	}
	// listen
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:"+strconv.Itoa(serverPort))
	if nil != err {
		util.Error("Listen failed 1 ! %s", err.Error())
	}
	listener, err := net.ListenTCP("tcp", tcpAddr) //监听一个端口
	defer listener.Close()
	if nil != err {
		util.Error("Listen failed 2 ! %s", err.Error())
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go ServerHandler(conn)
	}
}

func ServerHandler(conn net.Conn) {
	data, _ := util.ReadData(conn)
	// result, err := ioutil.ReadAll(conn)
	var files []*util.FileInfo
	err2 := json.Unmarshal(data, &files)
	if nil != err2 {
		util.Error(err2.Error())
	}
	// fmt.Println(len(files))
	for _, file := range files {
		CompareFile(file)
		fmt.Println(file)
	}
	data, err := json.Marshal(files)
	if nil != err {
		util.Error(err.Error())
	} else {
		conn.Write(data)
	}
	defer conn.Close()
}

func EnumFileHandler(file os.FileInfo, suffix string) {
	// util.Info(suffix + file.Name())
	fileinfo := &util.FileInfo{
		file.Name(),
		file.Size(),
		file.Mode(),
		file.ModTime(),
		file.IsDir(),
		suffix + file.Name(),
		"",
	}
	datas = append(datas, fileinfo)
}

func CompareFile(file *util.FileInfo) (int, error) {
	var localPath string
	if strings.HasPrefix(file.FullPath, "."+string(os.PathSeparator)) {
		nameRune := []rune(file.FullPath)
		localPath = syncDir + string(nameRune[2:])
		// fmt.Println(localPath)
	}
	if !util.FileExist(localPath) {
		file.Flag = "not found in server"
		return -1, nil
	} else {
		// compare time
		thisFile, err := os.Stat(localPath)
		if nil != err {
			fmt.Println(err.Error())
			return 9, nil
		}
		if thisFile.ModTime().Before(file.Time) {
			file.Flag = "newer in server"
		} else if thisFile.ModTime().After(file.Time) {
			file.Flag = "oldder in server"
		} else {
			file.Flag = "same"
		}
	}
	return 0, nil
}
