package falcon

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type OnCallGroInfo struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
xxxxxxxx
}

// GetOnCallGroInfo 获取 OnCall 组信息
func (f *FALCON) GetOnCallGroInfo (onCallGroName string) (*OnCallGroInfo, error) {
	path := fmt.Sprintf("xxxxxxxxxxoncallGroup=%s&mine=false&limit=2&offset=0", onCallGroName)
	if res, err := f.GetInfos(path); err == nil {
		var onCallGroInfoList []OnCallGroInfo
		err = json.Unmarshal(*res, &onCallGroInfoList)
		if err != nil {
			fmt.Println("json解析失败：", err)
			return nil, err
		}

		// 查询的结果有多条，这里只返回准确匹配的
		for _, onCallGroInfo := range onCallGroInfoList {
			if onCallGroInfo.Name == onCallGroName {
				return &onCallGroInfo, nil
			}
		}

		return nil, fmt.Errorf("未查询到告警组信息")
	} else {
		fmt.Println("获取结果失败：", err)
		return nil, err
	}
}

type Expression struct {
	ID                 int         `json:"id"`
	MetricName         string      `json:"metric_name"`
xxxxxxx
	Func               string      `json:"func"`
xxxxx
	Thresholds         []struct {
		xxxxxxxxx
		RunEnd    string `json:"run_end"`
	} `json:"thresholds"`
	OncallGroupName string      `json:"oncall_group_name"`
	OncallGroupID   int         `json:"oncall_group_id"`
	OncallGroups    interface{} `json:"oncall_groups"`
}

type ExpressionsByOncallGroStruct struct {
	Expressions []Expression `json:"expressions"`
	Count int `json:"count"`
}

// GetExpressionsByOncallGro 获取 onCall组 下面的告警表达式
func (f *FALCON) GetExpressionsByOncallGro (onCallGroupName string) (*ExpressionsByOncallGroStruct, *[]error) {

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errMsg []error

	// 每次返回数据的数目
	limit := 100

	var expressionsByOncallGroList ExpressionsByOncallGroStruct

	// 定义并发部分
	emit := func(offset int) {
		defer wg.Done()
		path := fmt.Sprintf("xxxxxxxxxx?oncall=%s&limit=%d&offset=%v", onCallGroupName, limit, offset)
		if ans, err := f.GetInfos(path); err == nil{
			var expressionsByOncallGro ExpressionsByOncallGroStruct
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
	path := fmt.Sprintf("xxxxxxxxxxxxxxxh?oncall=%s&limit=%d", onCallGroupName, limit)

	if body, err := f.GetInfos(path); err == nil {
		err = json.Unmarshal(*body, &expressionsByOncallGroList)
		if err != nil {
			errMsg = append(errMsg, err)
			return nil, &errMsg
		}

		// 如果数目大于 limit，则需要并发处理
		if expressionsByOncallGroList.Count > limit {
			page := (expressionsByOncallGroList.Count - 1) / limit + 1
			wg.Add(page - 1)
			for i := 1; i < page; i ++ {
				go emit(i * 100)
			}
			wg.Wait()
		}

		return &expressionsByOncallGroList, nil
	} else {
		errMsg = append(errMsg, err)
		return nil, &errMsg
	}
}

// GetExpressionByExprID 通过 ID 获取告警表达式
func (f *FALCON) GetExpressionByExprID (exprID int) (*Expression, error) {
	path := fmt.Sprintf("xxxxxx?id=%d",exprID)
	if body , err := f.GetInfos(path); err == nil {
		var expressionByExprID ExpressionsByOncallGroStruct
		err = json.Unmarshal(*body, &expressionByExprID)
		if err != nil {
			return nil, err
		}
		// 这里使用 id 查询，是唯一的
		expression := expressionByExprID.Expressions[0]
		return &expression, nil
	} else {
		return nil, err
	}
}

type PutUpdateInfo struct {
	Action struct {
	xxxxxxxxxx_groups"`
	} `json:"action"`
	Expression struct {
		Metric        string   `json:"metric"`
xxxxxxxxlback_mail"`
		OrgID              int         `json:"org_id"`
		AddThresholds []struct {
			xxxxxx
		} `json:"add_thresholds"`
		UpdateThresholds []UpdateThreshold `json:"update_thresholds"`
	} `json:"expression"`
}

type UpdateThreshold struct {
	Id        int    `json:"id"`
	xxxxxxxxx
}

type OncallGroups struct {
	ID int `json:"id"`
	Name string `json:"name"`
}

// ChangeExprLabels 修改表达式的 labels， 请注意，有的数据如果查不到，则会默认为空！！
func (f *FALCON) ChangeExprLabels (alertExpressionInfo *Expression) bool {

	var putUpdateInfo PutUpdateInfo

	// 动作部分 Action，
	xxxxxxxxxxxxxxx

	// 修改之前存储一下日志
	recordPath := fmt.Sprintf("/expression/search?id=%d",(*alertExpressionInfo).ID)
	if body, err := f.GetInfos(recordPath); err == nil {
		if ok := f.Logging("expression", body); !ok {
			fmt.Println("记录表达式日志失败")
			return false
		}
	} else {
		fmt.Println("日志-获取表达式信息失败：", err)
		return false
	}

	// 开始修改报警表达式
	path := fmt.Sprintf("/expression/update/%v", alertExpressionInfo.ID)
	if res, err := f.PutInfos(path, putUpdateInfo); err == nil {
		fmt.Println("修改表达式成功：", string(*res))
		return true
	} else {
		fmt.Println("修改表达式失败：", err)
		return false
	}
}

// ModExprLabels  bool 表示是否是 xxxx
// 当返回 map 为 nil 或者 bool 为 真时，均不需要修改
func (f *FALCON) ModExprLabels (includeLabels, modList *[]string) (*[]string, bool) {
	if *includeLabels == nil {
		return nil, false
	}

	modLabels := []string{}
	fieldExitFlag := false
	for _, label := range (*includeLabels) {
		tmpLabel := strings.Split(label, "=")
		key := tmpLabel[0]
		vals := strings.Split(tmpLabel[1], ",")

		// 如果有标签 xxx=xxxx，则不进行修改
		if key == "xxxx" {
			for _, val := range vals {
				if val == "xxxxxx" {
					//fmt.Println("label含有xxxx")
					return nil, true
				}
			}
		}

		tmpLabels := []string{}
		// 在这里对 目标字段进行处理
		if key == (*modList)[0] {
			fieldExitFlag = true
			keyExitFlag := false
			for _, val := range vals {
				// 如果目标值已经存在，则不添加
				if val == (*modList)[2] {
					//fmt.Println("目标值已经存在")
					return nil, false
				}

				// 如果 匹配值 c4 存在，则添加目标值 ak
				if val == (*modList)[1] {
					val = val + "," + (*modList)[2]
					keyExitFlag = true
				}
				tmpLabels = append(tmpLabels, val)
			}

			if !keyExitFlag {
				//fmt.Println("目标值不存在")
				return nil, false
			}
		} else {
			tmpLabels = vals
		}
		// 存储处理后的 label
		labels := fmt.Sprintf("%s=%s", key, strings.Join(tmpLabels, ","))
		modLabels = append(modLabels, labels)
	}
	if !fieldExitFlag {
		//fmt.Println("字段不存在", (*modList)[0])
		return nil, false
	}
	//fmt.Println("modLabels: ", modLabels)
	return &modLabels, false
}

type TemplatesByOnCallGroStruct struct {
	xxxxxxxxxxxx
}

// GetTemplatesByOnCallGro 根据 onCall组 获取告警模板
func (f *FALCON) GetTemplatesByOnCallGro (onCallGroupName string) (*[]TemplatesByOnCallGroStruct, *[]error) {

	var templatesList []TemplatesByOnCallGroStruct

	var errMsg []error
	var wg sync.WaitGroup
	var stopOnce sync.Once
	var mu sync.Mutex
	stop := make(chan struct{})
	maxGoroutines := 10 // 限制同时运行的最大goroutine数量
	guard := make(chan struct{}, maxGoroutines)

	emit := func(offset int) {
		defer func() {
			<- guard
			wg.Done()
		}()
		path := fmt.Sprintf("xxxxxxxxxxname=%s&limit=10&offset=%v", onCallGroupName, offset)
		if templatesBody, err := f.GetInfos(path); err == nil {
			var templates []TemplatesByOnCallGroStruct
			err = json.Unmarshal(*templatesBody, &templates)
			if err != nil {
				//fmt.Println("json 解析失败：", err)
				errMsg = append(errMsg, err)
				return
			}

			// 此时没数据查询完毕，停止并发
			if len(templates) < 10 {
				mu.Lock()
				templatesList = append(templatesList, templates...)
				mu.Unlock()
				stopOnce.Do(func() {
					close(stop)
				})
			} else {
				mu.Lock()
				templatesList = append(templatesList, templates...)
				mu.Unlock()
			}
		} else {
			errMsg = append(errMsg, err)
		}
	}

	offset := 0
	for {
		select {
		case <- stop:
			goto wait
		default:
			wg.Add(1)
			guard <- struct{}{}
			go emit(offset)
			offset = offset + 10
		}
	}
wait:
	wg.Wait()
	if len(errMsg) > 0 {
		return nil, &errMsg
	}
	//fmt.Println(templatesList, errMsg)
	return &templatesList, nil
}

type TemplateByIDStruct struct {
	Template struct {
		xxxxxxxx
	} `json:"template"`
	Action struct {
		xxxxxxxx
	} `json:"action"`
}

func (f *FALCON) GetTemplatesByID (templateID int) (*TemplateByIDStruct, error) {
	path := fmt.Sprintf("/xxxxxx?id=%v", templateID)
	if body, err := f.GetInfos(path); err == nil {
		var templateByID TemplateByIDStruct
		err = json.Unmarshal(*body, &templateByID)
		if err != nil {
			return nil, err
		}
		return &templateByID, nil
	} else {
		return nil, err
	}
}

func (f *FALCON) ChangeTemandStraTags (tags *string, modList *[]string) (string, bool) {

	if *tags == "" {
		return "", false
	}

	modTags := []string{}
	keyExitFlag := false
	fmt.Println("Tags: ", *tags)
	for _, tag := range strings.Split(*tags, ",") {
		tmpKey := strings.Split(tag, "=")
		key := tmpKey[0]
		vals := strings.Split(tmpKey[1], "|")

		// 如果有标签 owt=mfe，则不进行修改
		if key == "owt" {
			for _, val := range vals {
				if val == "xxxx" {
					//fmt.Println("tag 含有 xxx")
					return "", true
				}
			}
		}

		tmpTags := []string{}
		// 在这里对 目标字段进行处理
		if key == (*modList)[0] {
			for _, val := range vals {
				// 如果目标值已经存在，则不添加
				if val == (*modList)[2] {
					//fmt.Println("目标值已经存在")
					// 如果模板存在该标签，则不需要修改模板；
					// 由于不跳过改模板会导致继续修改模板下的策略，这里返回 true
					return "", true
					//return "", false
				}

				// 如果 匹配值 xxxx 存在，则添加目标值 xxxx
				if val == (*modList)[1] {
					val = val + "|" + (*modList)[2]
					keyExitFlag = true
				}
				tmpTags = append(tmpTags, val)
			}

		} else {
			tmpTags = vals
		}
		// 存储处理后的 tags
		tags := fmt.Sprintf("%s=%s", key, strings.Join(tmpTags, "|"))
		modTags = append(modTags, tags)
	}
	if !keyExitFlag {
		//fmt.Println("key 不存在：", (*modList)[1])
		return "", false
	}
	//fmt.Println(strings.Join(modTags, ","))
	return strings.Join(modTags, ","), false
}

// ModTemplate 修改告警模板
func (f *FALCON) ModTemplate (templateInfo *TemplateByIDStruct) bool {

	// 修改之前存储一下日志
	recordPath := fmt.Sprintf("xxxxe?id=%v",  (*templateInfo).Template.ID)
	if body, err := f.GetInfos(recordPath); err == nil {
		if ok := f.Logging("template", body); !ok {
			fmt.Println("记录模板日志失败")
			return false
		}
	} else {
		fmt.Println("日志-获取模板信息失败：", err)
		return false
	}

	path := fmt.Sprintf("/template/%v", (*templateInfo).Template.ID)
	if res, err := f.PutInfos(path, *templateInfo); err == nil {
		fmt.Println("更新成功的模板: ", string(*res))
		return true
	} else {
		fmt.Println("记录日志获取模板信息失败：", err)
		return false
	}
}

type StrategyByTempIDStruct struct {
	xxxx
	Strategies []struct {
		xxxxx
	} `json:"strategies"`
	Pause int `json:"pause"`
}

func (f *FALCON) GetStrategyByTempID (templateID int) (*[]StrategyByTempIDStruct, *[]error) {
	var strategyByTempList []StrategyByTempIDStruct

	var errMsg []error
	var mu sync.Mutex
	var stopOnce sync.Once
	var wg sync.WaitGroup
	maxGoroutines := 10 // 限制同时运行的最大goroutine数量
	guard := make(chan struct{}, maxGoroutines)

	stop := make(chan struct{})

	emit := func(offset int) {
		defer func() {
			<- guard
			wg.Done()
		}()
		path := fmt.Sprintf("/strategy/search?tid=%d&limit=10&offset=%d", templateID, offset)
		if strategyByTempBody, err := f.GetInfos(path); err == nil {
			var strategyByTemp []StrategyByTempIDStruct

			err = json.Unmarshal(*strategyByTempBody, &strategyByTemp)
			if err != nil {
				errMsg = append(errMsg, err)
				return
			}

			if len(strategyByTemp) < 10 {
				mu.Lock()
				strategyByTempList = append(strategyByTempList, strategyByTemp...)
				mu.Unlock()
				stopOnce.Do(func() {
					close(stop)
				})
			} else {
				mu.Lock()
				strategyByTempList = append(strategyByTempList, strategyByTemp...)
				mu.Unlock()
			}

		} else {
			errMsg = append(errMsg, err)
		}
	}

	offset := 0
	for {
		select {
		case <- stop:
			goto wait
		default:
			wg.Add(1)
			guard <- struct{}{}
			go emit(offset)
			offset += 10
		}
	}

	wait:
	wg.Wait()
	if len(errMsg) > 0 {
		return nil, &errMsg
	}
	return &strategyByTempList, nil

}

func (f *FALCON) ModStrategy (strategies *StrategyByTempIDStruct) bool {
	var putUpdateStrategies PutUpdateStrategies
xxxxxxxxxxxxxx

	// 修改之前存储一下日志
	recordPath := fmt.Sprintf("/xxxxx?id=%v", strategies.ID)
	if body, err := f.GetInfos(recordPath); err == nil {
		if ok := f.Logging("strategy", body); !ok {
			fmt.Println("记录策略日志失败：")
			return false
		}
	} else {
		fmt.Println("日志-获取策略信息失败：", err)
		return false
	}

	path := fmt.Sprintf("/xxxx/%d", strategies.ID)
	if res, err := f.PutInfos(path, putUpdateStrategies); err == nil {
		fmt.Println("更新成功的策略： ", string(*res))
		return true
	} else {
		fmt.Println("修改报警策略失败：", err)
		return false
	}
}

type PutUpdateStrategies struct {
xxxxxx
	AddStrategies []struct {
		Metric      string `json:"metric"`
	xxxxx
	} `json:"add_strategies"`
	xxxxxxxxe_strategies"`
}

type UpdateStrategy struct {
	ID          int    `json:"id"`
xxxxxxxxxx
}

func (f *FALCON) ModStrategyByTempID (templateID int, modList *[]string) bool {
	if strategyList, errMsg := f.GetStrategyByTempID(templateID); errMsg == nil {
		//changeFlag := false
		for _, strategies := range *strategyList {
			for idx, stragy := range strategies.Strategies {
				tag := stragy.Tags
				tags, mfeFlag := f.ChangeTemandStraTags(&tag, modList)
				if mfeFlag {
					//fmt.Println("含有mfe，不进行修改")
					continue
				}

				// 该规则需要修改
				if tags != "" {
					//changeFlag = true
					strategies.Strategies[idx].Tags = tags
					if ok := f.ModStrategy(&strategies); !ok {
						//fmt.Println("修改策略失败")
						return false
					} else {
						//fmt.Println("修改策略成功")
					}
				} else {
					//fmt.Println("未匹配到目标字段")
				}
			}
			//if changeFlag {
			//	if ok := f.ModStrategy(&strategies); !ok {
			//		//fmt.Println("修改策略失败")
			//		return false
			//	} else {
			//		//fmt.Println("修改策略成功")
			//	}
			//	changeFlag = false
			//}
		}
	} else {
		fmt.Println(errMsg)
		return false
	}
	return true
}
