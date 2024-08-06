package falcon

import (
	"fmt"
	"sync"
)

func (f *FALCON) ModExprLabelByOnCallGroService (onCallGroupName string, modList *[]string) bool {
	if expressionList, errMsg := f.GetExpressionsByOncallGro(onCallGroupName); errMsg == nil {
		for _, expression := range (*expressionList).Expressions {

			// 由于 falcon 侧原因，当前拿不到 action 部分的数据，这里需要通过 ID 进行二次查询
			expressionDetail, err := f.GetExpressionByExprID(expression.ID)
			if err != nil {
				//fmt.Println("获取表达式信息失败", err)
				continue
			}

			includeLabels := (*expressionDetail).IncludeLabels
			modLabels, mfeFlag := f.ModExprLabels(&includeLabels, modList);
			if mfeFlag {
				//fmt.Println("含有mfe，不进行修改")
				continue
			}

			if modLabels != nil {
				(*expressionDetail).IncludeLabels = *modLabels
			} else {
				//fmt.Println("不符合修改条件")
				continue
			}

			if sucess := f.ChangeExprLabels(expressionDetail); sucess {
				//fmt.Println("修改告警规则成功")
				//return true
			} else {
				fmt.Println("修改告警规则失败")
				return false
			}
		}
		return true
	} else {
		//fmt.Println(*errMsg)
		return false
	}
}

func (f *FALCON) ModTempansStraTagByOnCallGroService (onCallGroupName string, modList *[]string) *[]error {
	var wg sync.WaitGroup
	if templateList, errMsg := f.GetTemplatesByOnCallGro(onCallGroupName); errMsg == nil {
		wg.Add(len(*templateList))
		for _, template := range *templateList {
			go func(templateID int) {
				defer wg.Done()
				f.ModTempansStraTagByTempIDService(templateID, modList)
			}(template.ID)
		}
		wg.Wait()
		return nil
	} else {
		return errMsg
	}
}

func (f *FALCON) ModTempansStraTagByTempIDService (templateID int, modList *[]string) bool {
	if templateInfo, err := f.GetTemplatesByID(templateID); err == nil {
		tag := (*templateInfo).Template.CommonLabels

		//fmt.Println("tag to handle: ", tag)
		//if tag == "" {
		//	fmt.Println("该策略没有 label：", templateID,templateInfo)
		//	return true
		//}

		tags, mfeFlag := f.ChangeTemandStraTags(&tag, modList)

		if mfeFlag {
			//fmt.Println("含有mfe，不进行修改")
			//continue
			return true
		}

		if tags != "" {
			// 开始修改告警模板
			(*templateInfo).Template.CommonLabels = tags
			if ok := f.ModTemplate(templateInfo); ok {
				//fmt.Println("修改告警模板成功")
				return true
			} else {
				//fmt.Println("修改告警模板失败")
				return false
			}
		} else {
			// 修改告警策略
			//fmt.Println("请修改告警策略")
			if ok := f.ModStrategyByTempID(templateID, modList); !ok {
				return false
			} else {
				return true
			}
		}
	} else {
		//fmt.Println("获取模板信息失败：", err)
		return false
	}
}
