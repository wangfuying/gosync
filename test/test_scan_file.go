package main

import (
	"fmt"
	"gosync/util"
	"os"
)

func main() {
	fmt.Println("Start Test !")
	util.ScanFile("./", "."+string(os.PathSeparator), Handler)
}

func Handler(file os.FileInfo, suffix string) {
	fmt.Println(suffix + file.Name())
}
