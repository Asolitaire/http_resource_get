package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

var (
	a, f, d string
	h       bool
)

func init() {
	flag.BoolVar(&h, "h", false, "这是帮助")
	flag.StringVar(&d, "d", "", "下载到指定目录")
	flag.StringVar(&a, "a", "", "url")
	flag.StringVar(&f, "f", "", "资源url文件(多行)")

	flag.Usage = usage
}

func main() {
	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(3)
	}
	if d != "" {
		fmt.Println("资源将下载到你指定的目录")
		_, err := dirExists(d)
		if err != nil {
			err := os.MkdirAll(d, os.ModePerm)
			if err != nil {
				fmt.Println("创建目录失败", err)
				os.Exit(3)
			}
		}
	} else {
		fmt.Println("资源将下载到当前目录下的getResource内")
	}
	if a != "" {
		fmt.Printf("下载:%v -> ", a)
		getToFile(&a)
	}
	if f != "" {
		_, err := fileToGet(&f)
		if err != nil {
			panic(err)
		}
	}
	os.Exit(3)
}

func usage() {
	fmt.Fprintf(os.Stderr, `get:download http resource file(Can automatically generate multiple levels of directories based on the PATH structure of the URL, not suitable for downloading large files). version: get/0.0.1
Usage: get [-h] [-d <路径>] [-a <单个url>] [-f <包含多行url的文件名>]

Options:
`)
	flag.PrintDefaults()
}

func dirExists(dir string) (bool, error) {
	_, err := os.Stat(dir)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func fileToGet(urlfile *string) (bool, error) {
	file, err := os.OpenFile(*urlfile, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("打开文件出错", err)
		return false, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, err
	}

	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		if err != nil {
			if err == io.EOF && line != "" {
				fmt.Printf("\nread line:%v->", line)
				getToFile(&line)
				return true, nil
			} else if err == io.EOF && line == "" {
				fmt.Println("File read ok!")
				return false, nil
			} else {
				fmt.Println("Read file error!", err)
				return false, err
			}
		}
		fmt.Printf("\nread line:%v->", line)
		if line != "" {
			getToFile(&line)
		}
	}
	return true, err
}

func getToFile(urlTxt *string) bool {
	u, err := url.Parse(*urlTxt)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	pwd, _ := os.Getwd()
	allPath := pwd + "/getResource" + u.Path
	if d != "" {
		allPath = d + u.Path
	}
	fileDir := path.Dir(allPath)
	isDir, err := dirExists(fileDir)
	if isDir != true {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", *urlTxt, nil)
	if err != nil {
		//panic(err)
		fmt.Printf(" -> url连接失败: %v\n", err)
		return false
	}
	//req.Header.Add("Accept", "*/*")
	//req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.5;q=0.4")
	//req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.15(0x17000f31) NetType/4G Language/zh_CN")
	req.Header.Add("Referer", "http://er.frrjs.cn:443/home/cn?i=1157620_&code=071aVg0w3p62oV27fq0w3Q4SIL0aVg05")
	response, err := client.Do(req)

	if err != nil {
		//panic(err)
		fmt.Printf(" 读取下载流失败: %v\n", err.Error())
		return false
	}

	defer response.Body.Close()
	b, err := ioutil.ReadAll(response.Body)

	_, err = respToFile(allPath, b)
	if err != nil {
		fmt.Printf("%v -> 保存文件失败: %v\n", allPath, err)
		return false
	}
	fmt.Printf("%v -> ok\n", allPath)
	return true
}

func respToFile(fileAllPath string, contents []byte) (bool, error) {
	createFile, err := os.Create(fileAllPath)
	if err != nil {
		fmt.Printf("\n创建文件%v失败，原因：%v\n", fileAllPath, err)
		return false, err
	}
	//defer createFile.close()

	_, err = createFile.Write(contents)
	if err != nil {
		fmt.Printf("\ndownload file->%v faild:%v\n", fileAllPath, err)
		return false, err
	}
	createFile.Sync()
	return true, err
}
