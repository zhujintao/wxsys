package middleware

import (
	"net/http"
 	"github.com/gin-gonic/gin"
	"time"
	"github.com/chanxuehong/wechat/mp/user/oauth2"
	"github.com/chanxuehong/sid"
	"github.com/chanxuehong/session"
	"github.com/chanxuehong/rand"

)

const (
	oriid     = ""
	appid     = ""
	appsecret = ""
	token     = ""
	enaeskey  = ""
)

var (
	sessionStorage = session.New(20*60, 60*60)
	backcallurl = ""
	oauth2Config   = oauth2.NewOAuth2Config(
		appid,                 
		appsecret,             
		backcallurl,		
		//"snsapi_userinfo",    
		"snsapi_base",     
	)


	oauth2Client = oauth2.Client{Config: oauth2Config,}
	state string
)


func WxOauth(c *gin.Context)  {
	path := c.Request.URL.Path
	if path != "/logout" {
		if code := c.Query("code"); code == "" {
			islogin := 1
			if cookie, err := c.Request.Cookie("sid");err == nil {
				if _, err := sessionStorage.Get(cookie.Value); err == nil {
					islogin = 0
				}
			}
			if islogin == 1 {
				state = string(rand.NewHex())
				oauth2Config.RedirectURI = backcallurl + path
				AuthCodeURL := oauth2Config.AuthCodeURL(state, nil)
				c.Redirect(302, AuthCodeURL)
			}

		}else {

				sid := sid.New()
				if err := sessionStorage.Add(sid, state); err != nil {
					c.String(401, "error")
				}
				cookie := http.Cookie{
					Name:     "sid",
					Value:    sid,
					HttpOnly: true,
				}
				http.SetCookie(c.Writer, &cookie)
				oauth2Client.Exchange(code)
		}

	}

}

