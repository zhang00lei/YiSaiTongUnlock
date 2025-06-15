// Package main
// @Description: 根据右键菜单解锁目录文件，或单个文件
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/bytedance/gopkg/util/gopool"
	"github.com/schollz/progressbar/v3"
)

const EDITOR_PATH = "C:\\Users\\pc\\AppData\\Local\\Temp\\GoLand"

var exe_path string
var wg sync.WaitGroup

//加密文件标志
var lockedByte []byte

func init_info() bool {
	pathTemp := filepath.Dir(os.Args[0])
	if pathTemp == EDITOR_PATH {
		exe_path, _ = os.Getwd()
	} else {
		exe_path = pathTemp
	}
	lockedByte = append(lockedByte, 20)
	lockedByte = append(lockedByte, 35)
	lockedByte = append(lockedByte, 101)
	return true
}

func ReadBlock(filePth string, bufSize int) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := make([]byte, bufSize)
	bfRd := bufio.NewReader(f)
	_, err = bfRd.Read(buf)
	return buf, err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("参数长度有误")
		fmt.Scanln()
		return
	}
	if !init_info() {
		return
	}

	pathTemp := os.Args[1]
	info, err := os.Stat(pathTemp)
	if err != nil {
		fmt.Println("无法获取文件或目录信息：", err)
		fmt.Scanln()
		return
	}
	if info.IsDir() {
		lock := sync.Mutex{}
		s := spinner.New(spinner.CharSets[59], 500*time.Millisecond)
		s.Prefix = "搜索加密文件中 "
		s.Start()

		allFile, _ := getAllFileIncludeSubFolder(pathTemp)
		needFile := getNeedUnlockFile(allFile)
		s.Stop()

		bar := progressbar.Default(int64(len(needFile)))
		unlockCount := 0
		poolTemp := gopool.NewPool("Unlock", 5, gopool.NewConfig())
		for _, filePath := range needFile {
			wg.Add(1)
			temp := filePath
			poolTemp.Go(func() {
				UnlockFile(temp)
				lock.Lock()
				unlockCount++
				bar.Add(1)
				lock.Unlock()
				wg.Done()
			})
		}
		wg.Wait()
		fmt.Println("操作完成")
		fmt.Scanln()
	} else if info.Mode().IsRegular() {
		//解密当前文件
		data, err := ReadBlock(pathTemp, 4)
		if err != nil {
			return
		}
		if !bytes.Equal(data[1:], lockedByte) {
			//log.Println("文件未加密，跳过解锁：", pathTemp)
			return
		}
		UnlockFile(pathTemp)
	} else {
		fmt.Println("文件类型不支持")
	}
}

//
// UnlockFile
//  @Description: 解密文件
//  @param pathTemp 需要解密的文件路径
//  @param unlockCfg
//
func UnlockFile(pathTemp string) {
	// 1. 保存原始时间
	fileInfo, err := os.Stat(pathTemp)
	if err != nil {
		log.Println("无法获取文件信息:", err)
		return
	}
	originalModTime := fileInfo.ModTime()

	// 2. 原有解密流程
	docPath := pathTemp + ".docx"
	if err := os.Rename(pathTemp, docPath); err != nil {
		log.Println("重命名失败:", err)
		return
	}

	unlockPath := filepath.Join(exe_path, "wps.exe")
	cmd := exec.Command(unlockPath, docPath) // 直接传参，避免字符串解析问题
	if err := cmd.Run(); err != nil {
		log.Println("解密失败:", err)
		return
	}

	// 3. 恢复文件并重置时间
	dstFilePath := docPath + ".temp"
	if err := os.Rename(dstFilePath, pathTemp); err != nil {
		log.Println("恢复文件名失败:", err)
		return
	}
	if err := os.Chtimes(pathTemp, originalModTime, originalModTime); err != nil {
		log.Println("恢复时间失败:", err)
	}
}

//
// getAllFileIncludeSubFolder
//  @Description: 获取目录下所有文件（包含子目录）
//  @param folder
//  @return []string
//  @return error
//
func getAllFileIncludeSubFolder(folder string) ([]string, error) {
	var result []string
	filepath.Walk(folder, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if !info.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	return result, nil
}

func getNeedUnlockFile(allFiles []string) []string {
	var result []string
	var lock sync.Mutex
	poolTemp := gopool.NewPool("Unlock", 50, gopool.NewConfig())
	for _, pathTemp := range allFiles {
		filePath := pathTemp
		wg.Add(1)
		poolTemp.Go(func() {
			defer wg.Done()
			isLocked := fileIsLocked(filePath)
			if isLocked {
				lock.Lock()
				result = append(result, filePath)
				lock.Unlock()
			}
		})
	}
	wg.Wait()
	return result
}

func fileIsLocked(filePath string) bool {
	data, err := ReadBlock(filePath, 4)
	if err != nil {
		return false
	}
	if !bytes.Equal(data[1:], lockedByte) {
		//log.Println("文件未加密，跳过解锁：", pathTemp)
		return false
	}
	return true
}
