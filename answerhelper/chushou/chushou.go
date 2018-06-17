package chushou

import (
	"encoding/xml"
	"io/ioutil"

	"strings"

	"bytes"
	"unicode"

	"github.com/fatih/color"
)

type Candidate struct {
	XMlName   xml.Name   `xml:"questions"`
	Questions []Question `xml:"question"`
}

type Question struct {
	XMlName     xml.Name `xml:"question"`
	ID          int      `xml:"id,attr"`
	Content     string   `xml:"content,attr"`
	RightOption int      `xml:"rightOption,attr"`
	Options     []Option `xml:"option"`
}

type Option struct {
	XMlName xml.Name `xml:"option"`
	ID      int      `xml:"id,attr"`
	Content string   `xml:"content,attr"`
}

var C *Candidate

func InitQuestion() {
	path := "./answerhelper/question.xml"
	content, err := ioutil.ReadFile(path)
	if err != nil {
		color.Red("open xml failed, %v", err)
		return
	}
	C = &Candidate{}
	err = xml.Unmarshal(content, C)
	if err != nil {
		color.Red("xml unmarshal failed,%v", err)
		return
	}
}

func GetAnswer(question string) string {
	var answer string
	question = SimplyQuestion(question)
	color.Green("简化后的题目为：%s\n", question)
Exit:
	for _, ques := range C.Questions {
		if strings.Contains(ques.Content, question) {
			for _, option := range ques.Options {
				if option.ID == ques.RightOption {
					answer = option.Content
					break Exit
				}
			}
		}
	}
	return answer
}

func SimplyQuestion(question string) string {
	strings.Replace(question, " ", "", -1)
	index := strings.Index(question, ".")
	if index != -1 {
		question = question[index:]
	}
	var buf bytes.Buffer
	var size int
	for _, c := range question {
		if unicode.IsDigit(c) && size < 3 {
			size++
			continue
		}
		if c == '.' {
			continue
		}
		if c == '?' || c == '？' {
			continue
		}
		buf.WriteRune(c)
	}
	question = buf.String()
	return question
}
