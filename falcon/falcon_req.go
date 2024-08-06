package falcon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// https://wiki.n.miui.com/pages/viewpage.action?pageId=536078551

func (f *FALCON) GetApiKey () (bool){
	// curl -i -X GET 'http://api.falcon.srv/api/v2/auth/key' -H 'ak:AKWS3GA5E6XCDMWEU4' -H 'sk:KHFTmVutRHamO5Tu2/S8EQqlSfgXfDgjQC/wwSvc'
	// X-Api-Key: 667571696e67666569/7995a5879c8aa06041f1ad947206a31ff9cb651a770a78cbe5d62351ab4a3d2a6a07d8bfd7c1e3a951e18c7a275e9f0b
	req, _ := http.NewRequest("GET", "http://api.falcon.srv/api/v2/auth/key", nil)
	// 添加 ak、sk
	req.Header.Set("ak", f.AccessKeyId)
	req.Header.Set("sk",f.SecretAccessKeyId)

	resp, err := (&http.Client{}).Do(req)
	defer resp.Body.Close()

	if err != nil {
		return false
	}

	if resp.Status == "204 No Content" {
		for k, v := range resp.Header {
			if k == "X-Api-Key" {
				f.ApiKey = v[0]
				return true
			}
		}
	}
	return false
}

func (f *FALCON) GetInfos (path string) (*[]byte, error) {
	backEnd := "http://api.falcon.srv/api/v2"
	url := backEnd + path
	req, _ := http.NewRequest("GET", url, nil)
	// 设置 api key
	req.Header.Set("x-api-key", f.ApiKey)

	resp, err := (&http.Client{}).Do(req)
	defer resp.Body.Close()

	if err != nil {
		return nil, err
	}

	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Printf("Get++%s++%s\n", url, body)
	return &body, nil
}

func (f *FALCON) PostInfos (path string, params interface{}) (*[]byte, error){
	backEnd := "http://api.falcon.srv/api/v2"
	url := backEnd + path
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println("json err:", err)
		return nil, err
	}
	postBody := bytes.NewBuffer(b)

	//fmt.Printf("Post++%s++%s\n", url, postBody)

	req, _ := http.NewRequest("POST", url, postBody)
	req.Header.Set("x-api-key", f.ApiKey)
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", url, err)
		return nil, err
	}

	if resp.Status == "200 OK" {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("读取数据失败：", err)
			return nil, err
		}
		return &body, nil
	} else {
		return nil, fmt.Errorf("异常状态码：%v", resp.Status)
	}
}

func (f *FALCON) PutInfos (path string, params interface{}) (*[]byte, error){
	backEnd := "http://api.falcon.srv/api/v2"
	url := backEnd + path
	b, err := json.Marshal(params)
	if err != nil {
		fmt.Println("json err:", err)
		return nil, err
	}

	putBody := bytes.NewBuffer(b)

	//fmt.Println("putBody: ", putBody)
	//fmt.Printf("Put++%s++%s\n", url, putBody)

	req, _ := http.NewRequest(http.MethodPut, url, putBody)
	req.Header.Set("x-api-key", f.ApiKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("请求失败", url, err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println("body: ", string(body))
	if err != nil {
		fmt.Println("读取数据失败：", err)
		return nil, err
	}

	if resp.Status == "200 OK" {
		return &body, nil
	} else {
		return nil, fmt.Errorf("异常状态码：%v; %s", resp.Status, string(body))
	}
}
