package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"lyanna/models"
	"time"
)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

func GenList(n int) []int {
	ret := make([]int, n)
	for i := 0; i < n; i++ {
		ret[i] = i
	}
	return ret
}

func Add(a1, a2 int) int {
	return a1 + a2
}

func RenderSingleCommnet(comment *models.Comment)(string,error){
	var (
		commentHTML string
		err error
	)
	funcMap := template.FuncMap{
		"dateFormat": DateFormat,
	}
	tpl := template.New("singleComment.html")
	tpl.Funcs(funcMap)
	tpl, err= tpl.ParseFiles("./views/front/singleComment.html")
	if err != nil {
		fmt.Println(err)
		return commentHTML,err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf,&comment)
	if err != nil {
		return commentHTML,err
	}
	commentHTML = buf.String()
	return commentHTML,nil
}
