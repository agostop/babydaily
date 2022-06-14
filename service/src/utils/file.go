/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"archive/zip"
	"bufio"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	loggers "management_backend/src/logger"
	"os"
	"path/filepath"
	"strings"
)

var log = loggers.GetLogger(loggers.ModuleWeb)

func Utf8ToGBK(text string) (string, error) {
	dst := make([]byte, len(text)*2)
	tr := simplifiedchinese.GB18030.NewEncoder()
	nDst, _, err := tr.Transform(dst, []byte(text), true)
	if err != nil {
		return text, err
	}
	return string(dst[:nDst]), nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer func() {
		err = src.Close()
		if err != nil {
			log.Error("src file close err :", err.Error())
		}
	}()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer func() {
		err = dst.Close()
		if err != nil {
			log.Error("src file close err :", err.Error())
		}
	}()
	return io.Copy(dst, src)
}

func RePlace(fileName, oldStr, newStr string) error {
	in, err := os.Open(fileName)
	if err != nil {
		log.Error("open file fail:", err)
		return err
	}
	defer func() {
		err = in.Close()
		if err != nil {
			log.Error("Close file fail:", err)
		}
		err = os.Remove(fileName)
		if err != nil {
			log.Error("Remove file fail:", err)
		}
	}()

	out, err := os.OpenFile(fileName+".sh", os.O_RDWR|os.O_CREATE, 0766)
	if err != nil {
		log.Error("Open write file fail:", err)
		return err
	}
	defer func() {
		err = out.Close()
		if err != nil {
			log.Error("out file close err :", err.Error())
		}
	}()

	br := bufio.NewReader(in)
	index := 1
	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error("read err:", err)
			return err
		}
		newLine := strings.Replace(string(line), oldStr, newStr, -1)
		_, err = out.WriteString(newLine + "\n")
		if err != nil {
			log.Error("write to file fail:", err)
			return err
		}
		index++
	}
	return nil
}

func Zip(srcDir string, zipFileName string) error {

	// 预防：旧文件无法覆盖
	err := os.RemoveAll(zipFileName)
	if err != nil {
		log.Error("Remove zipFile err :", err.Error())
		return err
	}

	// 创建：zip文件
	zipfile, _ := os.Create(zipFileName)
	defer func() {
		err = zipfile.Close()
		if err != nil {
			log.Error("zip file close err :", err.Error())
		}
	}()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer func() {
		err = archive.Close()
		if err != nil {
			log.Error("zip file close err :", err.Error())
		}
	}()

	// 遍历路径信息
	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == srcDir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcDir+`\`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer func() {
				err = file.Close()
				if err != nil {
					log.Error("zip file close err :", err.Error())
				}
			}()
			_, err = io.Copy(writer, file)
			if err != nil {
				log.Error("Copy file err :", err.Error())
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Error("filepath Walk err :", err.Error())
		return err
	}
	return nil
}
