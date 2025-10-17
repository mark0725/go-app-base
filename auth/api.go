package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	base_db "github.com/mark0725/go-app-base/db"
	"github.com/mark0725/go-app-base/entities"
	base_web "github.com/mark0725/go-app-base/web"
)

type AuthApi struct{}

var g_AuthApi = &AuthApi{}

func (obj *AuthApi) Login(c *gin.Context) {
	logger.Trace("reqData:", c.Request.Body)

	reqParams := map[string]any{}
	if err := c.ShouldBindJSON(&reqParams); err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, base_web.ApiReponse{Code: "BadRequest", Message: "Message parameter error"})
		return
	}

	if _, ok := reqParams["username"]; !ok {
		logger.Error("username is required")
		c.JSON(http.StatusBadRequest, base_web.ApiReponse{Code: "BadRequest", Message: "username is required"})
		return
	}

	if _, ok := reqParams["password"]; !ok {
		logger.Error("password is required")
		c.JSON(http.StatusBadRequest, base_web.ApiReponse{Code: "BadRequest", Message: "password is required"})
		return
	}

	username := reqParams["username"].(string)
	sqlParams := map[string]any{
		"USER_ID": username,
		"ORG_ID":  g_appConfig.Org.OrgId,
	}

	recs, err := base_db.DBQueryEnt[entities.IdmUser](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_IDM_USER, "ORG_ID={ORG_ID} AND USER_ID={USER_ID}", sqlParams)
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		c.JSON(http.StatusInternalServerError, base_web.ApiReponse{Code: "ERROR", Message: "DBQueryEnt fail"})
		return
	}

	if len(recs) == 0 {
		c.JSON(http.StatusBadRequest, base_web.ApiReponse{Code: "DATA_EXIST", Message: "not found user " + username})
		return
	}

	userInfo := recs[0]

	password := reqParams["password"].(string)
	password_hash := sha256.Sum256([]byte(password))
	if userInfo.Passwd != hex.EncodeToString(password_hash[:]) {
		c.JSON(http.StatusBadRequest, base_web.ApiReponse{Code: "USER_PASSWORD_FAULT", Message: "password error"})
		return
	}

	authedConsumer := base_web.AuthenticatedConsumer{
		Id:           userInfo.UserId,
		Username:     userInfo.UserName,
		OrgId:        userInfo.OrgId,
		ConsumerType: "user",
	}
	c.Set(base_web.CtxKeyAuthenticatedConsumer, &authedConsumer)
	session := sessions.Default(c)
	sessionData, _ := json.Marshal(authedConsumer)
	session.Set(base_web.CtxKeyAuthenticatedConsumer, string(sessionData))
	if err := session.Save(); err != nil {
		logger.Error("session save", "error", err)
	}

	c.JSON(http.StatusOK, base_web.ApiReponse{Code: "OK", Message: "OK"})
}

func (obj *AuthApi) Loginout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(base_web.CtxKeyAuthenticatedConsumer)
	if err := session.Save(); err != nil {
		logger.Error("session save", "error", err)
	}
	c.Redirect(http.StatusFound, "/")
}
