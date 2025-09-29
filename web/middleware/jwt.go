package middleware

import (
	"github.com/gin-gonic/gin"
	base_utils "github.com/mark0725/go-app-base/utils"
	base_web "github.com/mark0725/go-app-base/web"
)

// JwtConfig mirrors the JWT plugin configuration parameters.
type JwtConfig struct {
	Anonymous         string   `json:"anonymous,omitempty"`          // Optional “anonymous” consumer identifier used when authentication fails.
	ClaimsToVerify    []string `json:"claims_to_verify,omitempty"`   // Registered claims to verify (accepted: "exp", "nbf").
	CookieNames       []string `json:"cookie_names,omitempty"`       // Cookie names inspected to retrieve JWTs (default: []).
	HeaderNames       []string `json:"header_names,omitempty"`       // HTTP header names inspected to retrieve JWTs (default: ["authorization"]).
	KeyClaimName      string   `json:"key_claim_name,omitempty"`     // Claim containing the key that identifies the secret (default: "iss").
	MaximumExpiration int      `json:"maximum_expiration,omitempty"` // Maximum allowed JWT lifetime in seconds, 0–31536000 (default: 0).
	Realm             string   `json:"realm,omitempty"`              // Value returned in the WWW-Authenticate realm attribute.
	RunOnPreflight    bool     `json:"run_on_preflight,omitempty"`   // Whether to run authentication on OPTIONS preflight requests (default: true).
	SecretIsBase64    bool     `json:"secret_is_base64,omitempty"`   // Treat credential secret as base64-encoded (default: false).
	UriParamNames     []string `json:"uri_param_names,omitempty"`    // Query parameters inspected to retrieve JWTs (default: ["jwt"]).
}

func newJwtConfig(params map[string]any) *JwtConfig {
	conf := JwtConfig{
		Anonymous: "",
	}

	base_utils.MapToStruct(params, &conf)

	return &conf
}
func (m *WebMiddlewareSecurity) Jwt(params map[string]any, r *gin.RouterGroup) gin.HandlerFunc {
	// conf := newKeyAuthConfig(params)

	return func(c *gin.Context) {
		if _, ok := c.Get(base_web.CtxKeyAuthenticatedConsumer); !ok {
			authedConsumer := base_web.AuthenticatedConsumer{
				Id:       "xilin.gao",
				Username: "高西林",
			}
			c.Set(base_web.CtxKeyAuthenticatedConsumer, &authedConsumer)
		}
		c.Next()

	}
}

// func GenJwtToken(issue string, audience string, userId string, exp int64) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"iss": issue,             // 签发者
// 		"sub": userId,            // 主题
// 		"aud": audience,          // 受众
// 		"iat": time.Now().Unix(), // 签发时间
// 		// "nbf": time.Now().Unix(),   // 生效时间
// 		"exp": exp,
// 	})

// 	tokenString, err := token.SignedString([]byte("secret"))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

// func VerifyJwtToken(tokenString string, secretKey string) (*jwt.Token, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		// 确保使用的是期望的签名方法
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		aud, ok := token.Claims.(jwt.MapClaims)["aud"].(string)
// 		if !ok {
// 			return nil, fmt.Errorf("aud claim is missing or not a string")
// 		}

// 		secretKey, err := getAppSecret(aud)
// 		if err != nil {
// 			return nil, fmt.Errorf("aud secret not found")
// 		}

// 		return secretKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 		fmt.Println(claims["sub"])
// 	} else {
// 		return nil, errors.New("invalid token")
// 	}

// 	return token, nil
// }

// func checkAppToken(token string) (string, error) {
// 	tokenInfo, err := getTokenInfo(token)
// 	if err != nil {
// 		return "", errors.New("token not found")
// 	}

// 	now := time.Now().Unix()
// 	if tokenInfo.ExpireTime != 0 && now > tokenInfo.ExpireTime {
// 		return "", errors.New("token expired")
// 	}

// 	return tokenInfo.AppId, nil
// }
