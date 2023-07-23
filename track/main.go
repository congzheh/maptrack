package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	dataOut = make(map[string]int)
	url     = "https://n.virtualairline.cn/p.php"
)

func main() {
	//无限循环
	for {
		Data := ProcessData()
		//判断dataOut的键是否有不是Data中的键，如果有，则对应的值+1
		for key, _ := range dataOut {
			if _, ok := Data[key]; !ok {
				if dataOut[key] < 3 {
					dataOut[key]++
				} else {
					//删除dataOut中的键
					delete(dataOut, key)
					//删除文件
					os.Remove("f/" + key + ".txt")
				}
			}
		}
		//处理Data中的数据，将经纬度数据（第四第五个），输出到key为文件名，对应的txt中
		for filename, values := range Data {
			//给dataOut中的filename设置值为0
			dataOut[filename] = 0
			// 检查文件是否存在
			_, err := os.Stat("f/" + filename + ".txt")
			if os.IsNotExist(err) {
				// 文件不存在，创建并写入数据
				fileContent := values[4] + "," + values[5] + "\n"
				err = ioutil.WriteFile("f/"+filename+".txt", []byte(fileContent), 0644)
				if err != nil {
					fmt.Println("无法写入文件:", err)
					return
				}
			} else {
				// 文件存在，读取原有数据并追加新数据
				fileContent, err := ioutil.ReadFile("f/" + filename + ".txt")
				if err != nil {
					fmt.Println("无法读取文件:", err)
					return
				}
				fileContent = append(fileContent, []byte(values[4]+","+values[5]+"\n")...)

				err = ioutil.WriteFile("f/"+filename+".txt", fileContent, 0644)
				if err != nil {
					fmt.Println("无法写入文件:", err)
					return
				}
			}
		}
		// 读取文件中原有的数据
		//10s一次
		time.Sleep(10 * time.Second)
	}
}

func ProcessData() map[string][]string {
	dataMap := make(map[string][]string)
	lines := make([]string, 0)
	resp, err := http.Get(url)
	/*
		测试数据：
			PILOT:1730:CSN881::29.515390:-98.280319:770:0:330::::37040:XNATC:1200::::
			PILOT:1730:CSN882::22.321335:113.881544:37:3:71:ZLXY:RJTT:WJC G212 VYK W34 VAPGU W100 ORAVA W201 UNSEK A326 DONVO G597 AGAVO Y644 GONAV G597 LANAT Y51 SAMON Y517 SYOEN Y31 GENJI Y10 GODIN:37040:XNATC:1200:B77L:35100:RJAA:
	*/
	scanner := bufio.NewScanner(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		goto doExit
	}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		goto doExit
	}
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) >= 3 {
			key := parts[2]
			value := parts

			/*
				value中的数据如下：
				   0 PILOT:
				   1 1730:
				   2 CSN881:
				   3 :
				   4 22.321335:
				   5 113.881544:
				   6 37:
				   7 3:
				   8 71:
				   9 ZLXY:
				   10 RJTT:
				   11 WJC G212 VYK W34 VAPGU W100 ORAVA W201 UNSEK A326 DONVO G597 AGAVO Y644 GONAV G597 LANAT Y51 SAMON Y517 SYOEN Y31 GENJI Y10 GODIN:
				   12 37040:
				   13 XNATC:
				   14 1200:
				   15 B77L:
				   16 35100:
				   17 RJAA:
			*/
			dataMap[key] = append(dataMap[key], value...)
		}
	}
	// 打印结果
	/*
		for key, values := range dataMap {
			fmt.Printf("键：%s\n", key)
			for _, value := range values {
				//i := 0
				//fmt.Printf(value + "\n")
			}
			//fmt.Println("---------")
		}
	*/
	return dataMap
doExit:
	fmt.Println("退出")
	return dataMap

}
