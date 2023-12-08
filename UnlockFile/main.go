// Package main
// @Description: 解锁单个文件
package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("参数长度有误")
		return
	}
	filePath := os.Args[1]
	dstFilePath := filePath + ".temp"
	copyFile(filePath, dstFilePath)
	err := os.Remove(filePath)
	if err != nil {
		log.Printf("文件%v未执行成功", filePath)
	}
}

func copyFile(sourcePath, dstFilePath string) (err error) {
	source, _ := os.Open(sourcePath)
	destination, _ := os.OpenFile(dstFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	defer source.Close()
	defer destination.Close()
	buf := make([]byte, 4096)
	for {
		n, err := source.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}
		if _, err := destination.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
