package test

import (
	"encoding/json"
	"falconService/falcon"
	"fmt"
	"strings"
	"sync"
	"testing"
)

//testPoint 22504
func TestGetOnCallGroInfo (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxx",
		SecretAccessKeyId: "xxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "testPoint"

	if onCallGroInfo, err := myFalcon.GetOnCallGroInfo(onCallGroupName); err == nil {
		fmt.Println((*onCallGroInfo).Name, (*onCallGroInfo).ID)
	} else {
		fmt.Println(err)
	}
}

func TestGetExpressionsByOncallGroTmp (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxx",
		SecretAccessKeyId: "xxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "testPoint"
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errMsg []error

	limit := 100

	var expressionsByOncallGroList falcon.ExpressionsByOncallGroStruct

	// 定义并发部分
	emit := func(offset int) {
		defer wg.Done()
		path := fmt.Sprintf("/xxxxx?oncall=%s&limit=%d&offset=%v", onCallGroupName, limit, offset)
		if ans, err := myFalcon.GetInfos(path); err == nil{
			var expressionsByOncallGro falcon.ExpressionsByOncallGroStruct
			err = json.Unmarshal(*ans, &expressionsByOncallGro)
			if err == nil && len(expressionsByOncallGro.Expressions) > 0 {
				mu.Lock()
				expressionsByOncallGroList.Expressions = append(expressionsByOncallGroList.Expressions, expressionsByOncallGro.Expressions...)
				mu.Unlock()
			} else {
				errMsg = append(errMsg, err)
			}
		} else {
			errMsg = append(errMsg, err)
		}
	}

	// 默认 limit 为 100
	path := fmt.Sprintf("xxxxxxxxxxx?oncall=%s&limit=%d", onCallGroupName, limit)

	if body, err := myFalcon.GetInfos(path); err == nil {
		err = json.Unmarshal(*body, &expressionsByOncallGroList)
		if err != nil {
			errMsg = append(errMsg, err)
		}

		// 如果数目大于 limit，则需要并发处理
		if expressionsByOncallGroList.Count > limit {
			page := (expressionsByOncallGroList.Count - 1) / limit + 1
			wg.Add(page - 1)
			for i := 2; i <= page; i ++ {
				go emit(i)
			}
			wg.Wait()
		}

	} else {
		errMsg = append(errMsg, err)
	}


}

func TestGetExpressionsByOncallGro (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxxxxxx",
		SecretAccessKeyId: "xxxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "【xxxxxl"
	//onCallGroupName := "testPoint"
	if expressionsByOncallGroList, errMsg := myFalcon.GetExpressionsByOncallGro(onCallGroupName); errMsg == nil {
		fmt.Println("Falcon count：", (*expressionsByOncallGroList).Count)
		fmt.Println("输出count：", len((*expressionsByOncallGroList).Expressions))
	} else {
		fmt.Println(*errMsg)
	}
}

func TestModLabels (t *testing.T) {

	addList := []string{"cluster", "c4", "ak"}

	includeLabels := []string{
		"cluster=c3,c4",
		"domain=posth5.g.mi.com",
		"job=mife-nginx-com",
		"owt=mfe",
	}

	labelMap := make(map[string][]string)
	for _, label := range includeLabels {
		tmpLabel := strings.Split(label, "=")
		key := tmpLabel[0]
		vals := strings.Split(tmpLabel[1], ",")

		// 如果有标签 owt=mfe，则不进行修改
		if key == "owt" {
			for _, val := range vals {
				if val == "mfe" {
					fmt.Println("label含有mfe")
					return
				}
			}
		}

		// 在这里对 labels 进行处理
		if key == addList[0] {
			var tmpLabels = []string{}
			keyExitFlag := false
			for _, val := range vals {
				// 如果目标值已经存在，则不添加
				if val == addList[2] {
					fmt.Println("目标值已经存在")
					return
				}

				// 如果 匹配值 c4 存在，则添加目标值 ak
				if val == addList[1] {
					val = val + "," + addList[2]
					keyExitFlag = true
				}
				tmpLabels = append(tmpLabels, val)
			}
			if !keyExitFlag {
				fmt.Println("目标值不存在")
				return
			}
		}

		labelMap[key] = vals
	}

	// 如果到这里，说明可以修改了
	fmt.Println(labelMap)

}

func TestGetTemplatesByOnCallGro (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxxxxx",
		SecretAccessKeyId: "xxxxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "testPoint"

	if templatesList, errMsg := myFalcon.GetTemplatesByOnCallGro(onCallGroupName); errMsg == nil {
		fmt.Println(templatesList)
	} else {
		fmt.Println("errMsg：", errMsg)
	}
}

func TestGetTemplatesByID (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxxxx",
		SecretAccessKeyId: "xxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}
	//16039
	templateID := 16039
	if templateInfo, err := myFalcon.GetTemplatesByID(templateID); err == nil {
		fmt.Println((*templateInfo).Template.CommonLabels)
	} else {
		fmt.Println("获取模板信息失败：", err)
	}
}

func TestModTagsTmp (t *testing.T) {
	tags := "xxx=xxx|ngin,xx=xxx,xx=xx|xxx,xx=xxx,xxx=xx|xxx"
	modList := []string{"cluster", "c4", "ak"}

	modTags := []string{}
	keyExitFlag := false
	for _, tag := range strings.Split(tags, ",") {
		tmpKey := strings.Split(tag, "=")
		key := tmpKey[0]
		vals := strings.Split(tmpKey[1], "|")

		// 如果有标签 owt=mfe，则不进行修改
		if key == "owt" {
			for _, val := range vals {
				if val == "mfe" {
					fmt.Println("tag 含有 mfe")
					//return nil, true
					return
				}
			}
		}

		tmpTags := []string{}
		// 在这里对 目标字段进行处理
		if key == modList[0] {
			keyExitFlag = true
			keyExitFlag := false
			for _, val := range vals {
				// 如果目标值已经存在，则不添加
				if val == modList[2] {
					fmt.Println("目标值已经存在")
					//return nil, false
					return
				}

				// 如果 匹配值 c4 存在，则添加目标值 ak
				if val == modList[1] {
					val = val + "|" + modList[2]
					keyExitFlag = true
				}
				tmpTags = append(tmpTags, val)
			}

			if !keyExitFlag {
				fmt.Println("目标值不存在")
				//return nil, false
				return
			}
		} else {
			tmpTags = vals
		}
		// 存储处理后的 label
		tags := fmt.Sprintf("%s=%s", key, strings.Join(tmpTags, "|"))
		//modTags += tags
		modTags = append(modTags, tags)
	}
	if !keyExitFlag {
		fmt.Println("key 不存在：", modList[0])
		return
	}
	fmt.Println(strings.Join(modTags, ","))

}




func TestGetStrategyByTempID (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxx",
		SecretAccessKeyId: "xxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	templateID := 16039

	if strategyList, errMsg := myFalcon.GetStrategyByTempID(templateID); errMsg == nil {
		//fmt.Println(*strategyList)
		for _, strategies := range *strategyList {
			fmt.Println(strategies)
		}
	} else {
		fmt.Println(errMsg)
	}

}


func TestModStrategyWay (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxxx",
		SecretAccessKeyId: "xxxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	templateID := 16039
	modList := []string{"cluster", "c4", "ak"}

	if strategyList, errMsg := myFalcon.GetStrategyByTempID(templateID); errMsg == nil {
		//fmt.Println(*strategyList)
		changeFlag := false
		for _, strategies := range *strategyList {
			for idx, stragy := range strategies.Strategies {
				//fmt.Println(stragy)
				tag := stragy.Tags
				tags, mfeFlag := myFalcon.ChangeTemandStraTags(&tag, &modList)
				if mfeFlag {
					fmt.Println("含有mfe，不进行修改")
					continue
				}

				// 该规则需要修改
				if tags != "" {
					changeFlag = true
					strategies.Strategies[idx].Tags = tags
				} else {
					fmt.Println("未匹配到目标字段")
				}
			}
			if changeFlag {
				if ok := myFalcon.ModStrategy(&strategies); ok {
					fmt.Println("修改策略成功")
				} else {
					fmt.Println("修改策略失败")
				}
			}
		}



	} else {
		fmt.Println(errMsg)
	}

}

func TestModStrategyByTempID (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxxxxxx",
		SecretAccessKeyId: "xxxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	templateID := 16039
	modList := []string{"cluster", "c4", "ak"}

	if ok := myFalcon.ModStrategyByTempID(templateID, &modList); ok {
		fmt.Println("修改模板下的策略成功")
	} else {
		fmt.Println("修改模板下的策略失败")
	}
}
