package falcon

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestGoBackExpr (t *testing.T) {

	myFalcon := FALCON{
		AccessKeyId: "xxxxx",
		SecretAccessKeyId: "xxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	fileName := "testPoint_20240713175207.txt"
	path := "../logs/"  + fileName

	file, err := os.OpenFile(path, os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file open fail", err)
		return
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
				return
			}
			expression := expressionByExprID.Expressions[0]
			if sucess := myFalcon.ChangeExprLabels(&expression); !sucess {
				fmt.Println("回退告警规则失败")
				return
			} else {
				continue
			}
		}

		if recordType[0] == "template" {
			var templateByID TemplateByIDStruct
			err = json.Unmarshal([]byte(recordType[1]), &templateByID)
			if err != nil {
				fmt.Println("解析模板失败：", err)
				return
			}

			if ok := myFalcon.ModTemplate(&templateByID); !ok {
				fmt.Println("回退模板失败")
			} else {
				continue
			}
		}

		if recordType[0] == "strategy" {
			var strategyByTempID StrategyByTempIDStruct
			err = json.Unmarshal([]byte(recordType[1]), &strategyByTempID)
			if err != nil {
				fmt.Println("解析策略失败：", err)
				return
			}

			if ok := myFalcon.ModStrategy(&strategyByTempID); !ok {
				fmt.Println("修改策略失败")
				return
				//return false
			} else {
				//fmt.Println("修改策略成功")
			}
		}
	}

}

