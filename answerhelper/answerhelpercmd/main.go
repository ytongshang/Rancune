package main

import (
	"strconv"
	"sync"
	"time"


	"strings"

	"github.com/ytongshang/rancune/answerhelper/chushou"
	"github.com/ytongshang/rancune/answerhelper/ocr"
	"github.com/ytongshang/rancune/answerhelper"
	"github.com/ytongshang/rancune/answerhelper/search"
	"github.com/ytongshang/rancune/answerhelper/util"
	"github.com/fatih/color"
 	"github.com/nsf/termbox-go"
)

func main() {
	answerhelper.InitConfig()
	chushou.InitQuestion()

	go func() {
		// 先去请求token
		ocr.GetToken()
	}()

	go func() {
		// 创建截图文件夹
		util.MkDirIfNotExist(answerhelper.ImagePath)
	}()

	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	color.Yellow("请按空格键开始搜索答案...\n")

Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeySpace:
				getQuestion()
				color.Yellow("\n\n请按空格键开始搜索答案...\n")
			default:
				break Loop
			}
		}
	}
}

func getQuestion() {
	start := time.Now()

	var accessToken string
	var imageBase64 string

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		token, err := ocr.GetToken()
		if err != nil {
			color.Red("获取百度Ocr token失败", err)
			return
		}
		color.Green("step1.获取百度Ocr token成功")
		accessToken = token
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 截图
		timeStamp := getUnixStamp()
		image, err := util.AndroidImageCapture(answerhelper.ImagePath + timeStamp + ".png")
		if err != nil {
			color.Red("截图失败", err)
			return
		}
		// 裁剪
		cutimage, err := util.ImageCut(image, answerhelper.AppConf.Left, answerhelper.AppConf.Top,
			answerhelper.AppConf.Right, answerhelper.AppConf.Bottom)
		if err != nil {
			color.Red("裁剪失败", err)
			return
		}
		// 保存
		cropfile := answerhelper.ImagePath + timeStamp + "_crop.png"
		err = util.SavePNGFile(cropfile, cutimage)
		if err != nil {
			color.Red("保存失败", err)
			return
		}
		imagestr, err := util.ImageToBase64(cropfile)
		if err != nil {
			color.Red("image base64失败", err)
			return
		}
		color.Green("step2.截图base64成功")
		imageBase64 = imagestr
	}()

	wg.Wait()

	if accessToken == "" || imageBase64 == "" {
		return
	}

	question, answer, err := ocr.Ocr(accessToken, imageBase64)
	if err != nil {
		color.Red("图像识别出错")
		return
	}
	color.Green("step3.识别题目成功\n\n")
	color.Green("题目：%s\n", question)

	color.Green("从题库中搜索答案：\n")
	ans := chushou.GetAnswer(question)
	if ans != "" {
		color.Green("题目存在于题库中，推荐答案：%s\n", ans)
	} else {
		color.Cyan("%s \n\n", question)
		opposite := strings.Contains(question, "不")
		if opposite {
			color.Red("注意题目中的否定\n\n")
		}
		result := search.GetSearchResult(question, answer)

		total := make(map[string]*search.Result)
		for _, answer := range answer {
			total[answer] = &search.Result{}
		}
		for engine, answerResult := range result {
			color.Red("================%s搜索==============", engine)

			for key, value := range answerResult {
				color.Green("%s : 结果总数 %d ， 答案出现频率： %d", key, value.Sum, value.Freq)
				total[key].Freq += value.Freq
				total[key].Sum += value.Sum
			}
			color.Red("======================================")
		}
		color.Red("================总共==============")
		for key, value := range total {
			color.Green("%s : 结果总数 %d ， 答案出现频率： %d", key, value.Sum, value.Freq)
			total[key].Freq += value.Freq
			total[key].Sum += value.Sum
		}
		color.Red("======================================")
	}

	color.Cyan("\n耗时：%v", time.Now().Sub(start))

}

func getUnixStamp() string {
	now := time.Now().Unix()
	return strconv.Itoa(int(now))
}
