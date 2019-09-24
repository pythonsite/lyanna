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

type HH struct{
	Comments []*models.Comment
	Githubuser interface{}
	Post *models.Post
	Pages int
	CommentNum int
}

func RenderAllComment(hh HH) (string,error) {
	var (
		AllCommentHTML string
		err error
	)
	funcMap := template.FuncMap{
		"dateFormat": DateFormat,
		"genList":GenList,
		"add": Add,
	}
	tpl := template.New("comment.html")
	tpl.Funcs(funcMap)
	tpl, err = tpl.ParseFiles("./views/front/comment.html")
	if err != nil {
		fmt.Println(err)
		return AllCommentHTML,err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf,hh)
	if err != nil {
		fmt.Println(err)
		return AllCommentHTML,err
	}
	AllCommentHTML = buf.String()
	return AllCommentHTML,nil
}

func RenderSingleComment(comment *models.Comment)(string,error){
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

func GetMapValue(key string, origin map[string][]*models.Post)[]*models.Post{
	return origin[key]
}