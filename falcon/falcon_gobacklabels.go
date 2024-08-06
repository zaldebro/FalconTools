package falcon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func (f *FALCON) GoBackLabels (fileName string) bool {
	path := "./logs/"  + fileName
	LogPath = path + ".back"
	file, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file open fail", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		recordType := strings.Split(scanner.Text(), "++")
		if recordType[0] == "expression" {
			var expressionByExprID ExpressionsByOncallGroStruct
			err = json.Unmarshal([]byte(recordType[1]), &expressionByExprID)
			if err != nil {
				fmt.Println("解析表达式失败：", err)
				return false
			}
			expression := expressionByExprID.Expressions[0]
			if sucess := f.ChangeExprLabels(&expression); !sucess {
				fmt.Println("回退告警规则失败")
				return false
			}
		} else if recordType[0] == "template" {
			var templateByID TemplateByIDStruct
			err = json.Unmarshal([]byte(recordType[1]), &templateByID)
			if err != nil {
				fmt.Println("解析模板失败：", err)
				return false
			}
			if ok := f.ModTemplate(&templateByID); !ok {
				fmt.Println("回退模板失败")
			}
		} else if recordType[0] == "strategy" {
			var strategyByTempID StrategyByTempIDStruct
			err = json.Unmarshal([]byte(recordType[1]), &strategyByTempID)
			if err != nil {
				fmt.Println("解析策略失败：", err)
				return false
			}
			if ok := f.ModStrategy(&strategyByTempID); !ok {
				fmt.Println("修改策略失败")
				return false
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("回退完成")
		}
	}

	return true
}











