package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PC struct {
	client *http.Client
	biz    string
	uin    string
}

type Article struct {
	AppMsgExtInfo struct {
		ContentURL string `json:"content_url"`
		Title      string `json:"title"`
	} `json:"app_msg_ext_info"`
	CommMsgInfo struct {
		Datetime int64 `json:"datetime"`
	} `json:"comm_msg_info"`
}

type Response struct {
	GeneralMsgList string `json:"general_msg_list"`
}

func NewPC(biz, uin string) (*PC, error) {
	client := &http.Client{}

	return &PC{
		client: client,
		biz:    biz,
		uin:    uin,
	}, nil
}

func (p *PC) getURLs(key string, offset int) ([]Article, error) {
	reqURL := "https://mp.weixin.qq.com/mp/profile_ext"
	params := url.Values{}
	params.Set("action", "getmsg")
	params.Set("__biz", p.biz)
	params.Set("f", "json")
	params.Set("offset", strconv.Itoa(offset))
	params.Set("count", "10")
	params.Set("uin", p.uin)
	params.Set("key", key)

	req, err := http.NewRequest("GET", reqURL+"?"+params.Encode(), nil)
	req.Header.Add("Cookies", "")
	if err != nil {
		return nil, err
	}

	resp, err := p.client.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res Response
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return nil, err
	}

	if res.GeneralMsgList == "" {
		return nil, errors.New("获取文章连接失败，检查参数")
	}

	var msgList struct {
		List []Article `json:"list"`
	}
	err = json.Unmarshal([]byte(res.GeneralMsgList), &msgList)
	if err != nil {
		return nil, err
	}

	return msgList.List, nil
}

/*
startTimestamp:文章开始时间
startCount:开始采集的条数
endCount:结束采集的条数
比如要采集第10-20条的文章
startCount:10
endCount:20
*/
func getHistoryURLs(biz, uin, key string, startTimestamp int64, startCount, endCount int) ([]Article, []string, error) {
	pc, err := NewPC(biz, uin)
	if err != nil {
		return nil, nil, err
	}

	var allArticles []Article
	var justURLs []string
	for {
		articles, err := pc.getURLs(key, startCount)
		if err != nil {
			return nil, nil, err
		}
		if len(articles) == 0 {
			break
		}
		startCount += 10
		allArticles = append(allArticles, articles...)
		for _, article := range articles {
			justURLs = append(justURLs, article.AppMsgExtInfo.ContentURL)
		}
		lastTimestamp := articles[len(articles)-1].CommMsgInfo.Datetime
		fmt.Println(startCount, timestampToDate(lastTimestamp))
		if lastTimestamp <= startTimestamp || startCount >= endCount {
			break
		}
		time.Sleep(5 * time.Second)
	}
	return allArticles, justURLs, nil
}

func timestampToDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	return t.Format("2006-01-02")
}

// 待完善
// 采集到文章的url后，传给DownloadHtml()
func main() {
	// 通过抓包工具获取，点击该公众号的文章即可从抓包工具里查到相关参数的值
	biz := ""
	uin := ""
	key := ""

	_, justURLs, err := getHistoryURLs(biz, uin, key, 0, 0, 1)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("抓取到的文章链接")
	for _, url := range justURLs {
		fmt.Println(url)
	}
	fmt.Println("Total:", len(justURLs))
}
