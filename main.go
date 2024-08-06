package main

import (
	"bufio"
	"falconService/falcon"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func LabeLs () {
	pwd, _ := os.Getwd()
	logsDir := pwd + "/logs"
	files, _ := ioutil.ReadDir(logsDir)
	fmt.Println("可回退的文件：")
	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() && (strings.HasSuffix(fileName, ".txt") || strings.HasSuffix(fileName, ".txt.back")) {
			fmt.Println(fileName)
		}
	}
}

func Help () {
	fmt.Println(
		"ls  查看可回退的标签\n" +
			"back <filename> 修改标签\n" +
			"mod 修改标签\n" +
			"help 查看帮助命令")
}

// falconLabel
func main() {
	pwd, _ := os.Getwd()
	logsDir := pwd + "/logs"
	_, err := os.Stat(logsDir)
	if err != nil {
		err := os.Mkdir(logsDir, os.ModePerm)
		if err != nil {
			fmt.Println("err: ", err)
			fmt.Println("在当前目录下创建文件夹 logs 失败，请手动创建")
			return
		}
	}

	var AKSK string
	var myFalcon falcon.FALCON
	fmt.Print("请输入ak++sk> ")
	fmt.Scanf("%s", &AKSK)
	AKSK = strings.Trim(AKSK, " ")
	tmpAKSK := strings.Split(AKSK, "++")
	if len(tmpAKSK) == 2 {
		myFalcon.AccessKeyId = tmpAKSK[0]
		myFalcon.SecretAccessKeyId = tmpAKSK[1]

		status := myFalcon.GetApiKey()
		if !status {
			fmt.Println("初始化 myFalcon 失败")
			return
		}
		fmt.Println("初始化falcon成功")
	} else {
		fmt.Println("输入格式 ak++sk")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("falconTools> ")
		scanner.Scan()
		input := scanner.Text()

		if err := scanner.Err(); err != nil {
			fmt.Println("读取输入失败：", err)
			return
		}
		commands := strings.Fields(input)

		if len(commands) == 0 {
			continue
		}

		switch commands[0] {
		case "ls":
			LabeLs()
		case "back":
			for _, fileName := range commands[1:] {
				pwd, _ := os.Getwd()
				logsPath := pwd + "/logs/" + fileName
				_, err := os.Stat(logsPath)
				if err != nil {
					fmt.Println("文件不存在", fileName)
				} else {
					if (strings.HasSuffix(fileName, ".txt") || strings.HasSuffix(fileName, ".txt.back")) {
						fmt.Println("开始回退 ", fileName)
						if ok := myFalcon.GoBackLabels(fileName); !ok {
							fmt.Println("回退失败", fileName)
						}
						fmt.Println("回退成功")
					} else {
						fmt.Println("该文件的格式可能不对", fileName)
					}
				}
			}
		case "mod":
			// mod testPoint cluster c4 ak
			if len(commands) != 5 {
				fmt.Println("输入格式可能存在问题")
				continue
			}

			onCallGroupName := commands[1]
			falcon.LogPath =  "./logs/" + onCallGroupName + "_" + time.Now().Format("20060102150405") + ".txt"

			modList := commands[2:]

			if ok := myFalcon.ModExprLabelByOnCallGroService(onCallGroupName, &modList); !ok {
				fmt.Println("修改表达式失败")
			}
			if errMsg := myFalcon.ModTempansStraTagByOnCallGroService(onCallGroupName, &modList); errMsg != nil{
				fmt.Println("模板和策略修改失败：", *errMsg)
			}
			fmt.Println("修改成功")

		case "exit":
			return
		default:
			Help()
		}
		
	}

}






