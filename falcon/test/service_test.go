package test

import (
	"falconService/falcon"
	"fmt"
	"testing"
)

func TestUpdateExprLabelsTmp (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxx",
		SecretAccessKeyId: "xxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "testPoint"
	modList := []string{"cluster", "c4", "ak"}

	if expressionList, errMsg := myFalcon.GetExpressionsByOncallGro(onCallGroupName); errMsg == nil {
		for _, expression := range (*expressionList).Expressions {

			// 在这里增加处理的规则即可
			includeLabels := expression.IncludeLabels
			modLabels, mfeFlag := myFalcon.ModExprLabels(&includeLabels, &modList);
			if mfeFlag {
				fmt.Println("含有mfe，不进行修改")
				continue
			}

			if modLabels != nil {
				expression.IncludeLabels = *modLabels
			} else {
				fmt.Println("不符合修改条件")
				continue
			}

			if sucess := myFalcon.ChangeExprLabels(&expression); sucess {
				fmt.Println("修改告警规则成功")
			} else {
				fmt.Println("修改告警规则失败")
			}
		}
	} else {
		fmt.Println(*errMsg)
	}
}


func TestUpdateExprLabels (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxxx",
		SecretAccessKeyId: "xxxxxxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "testPoint"
	modList := []string{"cluster", "c4", "ak"}

	if ok := myFalcon.ModExprLabelByOnCallGroService(onCallGroupName, &modList); ok {
		fmt.Println("修改规则成功")
	} else {
		fmt.Println("修改规则失败")
	}

}


func TestUpdateTemandStra (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "AKWS3GA5E6XCDMWEU4",
		SecretAccessKeyId: "KHFTmVutRHamO5Tu2/S8EQqlSfgXfDgjQC/wwSvc",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	onCallGroupName := "testPoint"
	modList := []string{"cluster", "c4", "ak"}

	if templateList, errMsg := myFalcon.GetTemplatesByOnCallGro(onCallGroupName); errMsg == nil {
		for _, template := range *templateList {
			if templateInfo, err := myFalcon.GetTemplatesByID(template.ID); err == nil {
				tag := (*templateInfo).Template.CommonLabels
				tags, mfeFlag := myFalcon.ChangeTemandStraTags(&tag, &modList)

				if mfeFlag {
					fmt.Println("含有mfe，不进行修改")
					continue
				}

				if tags != "" {
					// 开始修改告警模板
					(*templateInfo).Template.CommonLabels = tags
					if ok := myFalcon.ModTemplate(templateInfo); ok {
						fmt.Println("修改告警模板成功")
						continue
					} else {
						fmt.Println("修改告警模板失败")
						return
					}
				} else {
					// 修改告警策略
					fmt.Println("请修改告警策略")
				}
			} else {
				fmt.Println("获取模板信息失败：", err)
			}
		}
	} else {
		fmt.Println("errMsg：", errMsg)
	}
}


func TestModTempansStraTagByTempIDService (t *testing.T) {
	myFalcon := falcon.FALCON{
		AccessKeyId: "xxxx",
		SecretAccessKeyId: "xxxxx",
	}

	status := myFalcon.GetApiKey()
	if !status {
		return
	}

	templateID := 16039
	modList := []string{"cluster", "c4", "ak"}

	if ok := myFalcon.ModTempansStraTagByTempIDService(templateID, &modList); !ok {
		fmt.Println("修改失败。。。。。。")
	} else {
		fmt.Println("修改成功。。。。。。。。。。")
	}
}






