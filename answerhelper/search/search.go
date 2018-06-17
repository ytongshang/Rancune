package search

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/ytongshang/rancune/httpclient"
)

type Result struct {
	Sum  int32
	Freq int32
}

func GetSearchResult(question string, answers []string) map[string]map[string]*Result {
	if question == "" {
		return nil
	}

	var wg sync.WaitGroup

	res := make(map[string]map[string]*Result)

	wg.Add(1)
	go func() {
		defer wg.Done()
		res["百度"] = generalSearch("http://www.baidu.com/s?wd=%s", `百度为您找到相关结果约([\d\,]+)`, question, answers)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		res["Sougo"] = generalSearch("https://www.sogou.com/web?query=%s", `搜狗已为您找到约([\d\,]+)`, question, answers)
	}()
	wg.Wait()

	return res
}

func generalSearch(uriFormat string, regxp string, question string, answers []string) map[string]*Result {
	resultMap := make(map[string]*Result, len(answers))

	searchURL := fmt.Sprintf(uriFormat, url.QueryEscape(question))
	questionBody, err := httpclient.Get(searchURL).Bytes()
	if err != nil {
		return nil
	}
	var wg sync.WaitGroup

	for _, answer := range answers {
		searchResult := new(Result)
		resultMap[answer] = searchResult

		wg.Add(1)
		go func(answer string, result *Result) {
			defer wg.Done()
			//搜question结果中answer出现的次数
			result.Freq = int32(strings.Count(string(questionBody), answer))
		}(answer, searchResult)

		wg.Add(1)
		go func(answer string, result *Result) {
			defer wg.Done()
			//搜question + answer 总共出现的次数
			keyword := fmt.Sprintf("%s %s", question, answer)
			searchURL := fmt.Sprintf(uriFormat, url.QueryEscape(keyword))
			body, err := httpclient.Get(searchURL).Bytes()
			if err != nil {
				color.Red("search %s error", answer)
			} else {
				reg, _ := regexp.Compile(regxp)
				find := reg.FindAllStringSubmatch(string(body), -1)
				if len(find) > 0 {
					sum := find[0][1]
					sum = strings.Replace(sum, ",", "", -1)
					result.Sum = MustInt32(sum)
				}
			}
		}(answer, searchResult)
	}
	wg.Wait()
	return resultMap
}

func MustInt32(str string) int32 {
	i, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return int32(0)
	}
	return int32(i)
}
