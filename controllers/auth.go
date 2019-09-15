package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/alimoeeny/gooauth2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"lyanna/models"
	"lyanna/utils"
	"net/http"
)

type GithubUserInfo struct {
	AvatarURL         string      `json:"avatar_url"`
	Email             string	  `json:"email"`
	GID        		  int64          `json:"id"`
	Name              string      `json:"name"`
	Login 			  string 	  `json:"login"`
	HtmlUrl 			string    `json:"html_url"`
}

func AuthGet(c *gin.Context) {
	uuid := utils.UUID()
	session := sessions.Default(c)
	session.Delete(models.SESSION_GITHUB_STATE)
	session.Set(models.SESSION_GITHUB_STATE, uuid)
	session.Save()
	authurl := fmt.Sprintf(models.Conf.GitHub.AuthUrl,models.Conf.GitHub.ClientID,uuid)
	c.Redirect(http.StatusFound, authurl)
}

func Oauth2Callback(c *gin.Context) {
	//var (
	//	userInfo *GithubUserInfo
	//	user     *models.User
	//)
	code := c.Query("code")
	state := c.Query("state")

	session := sessions.Default(c)
	fmt.Println(state)
	fmt.Println(session.Get(models.SESSION_GITHUB_STATE))
	if len(state) == 0 || state != session.Get(models.SESSION_GITHUB_STATE) {
		c.Abort()
		return
	}
	session.Delete(models.SESSION_GITHUB_STATE)
	session.Save()

	token, err := GetTokenByCode(code)
	fmt.Printf("token:%v",token)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	gitHubUser, err := getGitHubUserInfoByToken(token)
	fmt.Println(gitHubUser)
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}

	gituser := &models.GitHubUser{
		GID:gitHubUser.GID,
		Email:gitHubUser.Email,
		Picture: gitHubUser.AvatarURL,
		UserName:gitHubUser.Name,
		NickName:gitHubUser.Login,
		Url:gitHubUser.HtmlUrl,
	}

	gituser,err = gituser.FirstOrCreate()
	if err == nil {
		s := sessions.Default(c)
		s.Clear()
		s.Set(models.SESSION_KEY, gituser.GID)
		s.Save()
		c.Redirect(http.StatusMovedPermanently,"/post/2")
	}


}

func GetTokenByCode(code string)(accessToken string, err error) {
	var (
		transport *oauth.Transport
		token     *oauth.Token
	)
	transport = &oauth.Transport{
		Config:&oauth.Config{
			ClientId: models.Conf.GitHub.ClientID,
			ClientSecret:models.Conf.GitHub.ClientSecret,
			TokenURL:models.Conf.GitHub.TokenUrl,
			RedirectURL:models.Conf.GitHub.RedirectUrl,
			Scope:"email profile",
		},
	}
	token, err = transport.Exchange(code)
	if err != nil {
		fmt.Println(err)
		return
	}
	accessToken = token.AccessToken
	tokenCache := oauth.CacheFile("./request.token")
	if err = tokenCache.PutToken(token); err != nil {
		log.Println(err)
		return
	}
	return
}

func getGitHubUserInfoByToken(token string)(*GithubUserInfo, error){
	var (
		resp *http.Response
		body []byte
		err  error
	)
	resp, err = http.Get(fmt.Sprintf("https://api.github.com/user?access_token=%s", token))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var userInfo GithubUserInfo
	err = json.Unmarshal(body, &userInfo)
	fmt.Println(userInfo)
	return &userInfo, err
}
