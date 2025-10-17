package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	base_db "github.com/mark0725/go-app-base/db"
	"github.com/mark0725/go-app-base/entities"
	base_utils "github.com/mark0725/go-app-base/utils"
	base_web "github.com/mark0725/go-app-base/web"
)

// ConfigObject represents the configuration schema described above.
type KeyAuthConfig struct {
	Anonymous                 string          `json:"anonymous,omitempty"`                   // optional consumer UUID or username used when authentication fails
	HideCredentials           bool            `json:"hide_credentials,omitempty"`            // if true, strips the credential from the request (default: false)
	IdentityRealms            []IdentityRealm `json:"identity_realms,omitempty"`             // list of Konnect Identity Realms used to source a consumer
	KeyInHeader               bool            `json:"key_in_header,omitempty"`               // enables reading the credential from request header (default: false)
	KeyInQuery                bool            `json:"key_in_query,omitempty"`                // enables reading the credential from query parameter (default: false)
	KeyInBody                 bool            `json:"key_in_body,omitempty"`                 // enables reading the credential from request body (default: false)
	KeyInAuthorizationHeader  bool            `json:"key_in_authorization_header,omitempty"` // enables reading the credential from Authorization header (default: true)
	AuthorizationHeaderPrefix string          `json:"authorization_header_prefix,omitempty"` // authorization header prefix (default: "Bearer")
	KeyNames                  []string        `json:"key_names,omitempty"`                   // header/query/body parameter names to look for (default: ["apikey"])
	Realm                     string          `json:"realm,omitempty"`                       // realm value returned in WWW-Authenticate header on auth failure
	RunOnPreflight            bool            `json:"run_on_preflight,omitempty"`            // if true, authenticates OPTIONS preflight requests (default: true)
}

// IdentityRealm defines an individual Konnect Identity Realm entry.
type IdentityRealm struct {
	ID     string `json:"id,omitempty"`     // UUID of the identity realm
	Region string `json:"region,omitempty"` // region identifier
	Scope  string `json:"scope,omitempty"`  // allowed values: "cp" or "realm"
}

func newKeyAuthConfig(params map[string]any) *KeyAuthConfig {
	conf := KeyAuthConfig{
		Anonymous:                 "",
		HideCredentials:           false,
		IdentityRealms:            []IdentityRealm{},
		KeyInHeader:               false,
		KeyInQuery:                false,
		KeyInBody:                 false,
		KeyInAuthorizationHeader:  true,
		AuthorizationHeaderPrefix: "Bearer",
		KeyNames:                  []string{"apikey"},
		Realm:                     "",
		RunOnPreflight:            false,
	}

	base_utils.MapToStruct(params, &conf)

	return &conf
}
func (m *WebMiddlewareSecurity) KeyAuth(params map[string]any, r gin.IRoutes) gin.HandlerFunc {
	conf := newKeyAuthConfig(params)
	return func(c *gin.Context) {
		if _, ok := c.Get(base_web.CtxKeyAuthenticatedConsumer); ok {
			c.Next()
			return
		}

		apiKey, err := extractAPIKey(c, conf)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Credential is required"})
			return
		}

		credentialInfo, err := getCredential(apiKey, conf)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Credential"})
			return
		}
		// if credentialInfo, err := repo.From[entities.AuthCredential]().Get(apiKey); err != nil {
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "Invalid API key"})
		// 	return
		// }
		c.Set(base_web.CtxKeyAuthenticatedConsumer, &base_web.AuthenticatedConsumer{
			Id:           credentialInfo.ObjectId,
			Username:     "",
			ConsumerType: credentialInfo.ObjectType,
			OrgId:        credentialInfo.OrgId,
		})

		credential := base_web.AuthenticatedCredential{
			Id:         credentialInfo.CredentialId,
			Identifier: credentialInfo.CredentialId,
			Name:       credentialInfo.CredentialName,
			Type:       credentialInfo.CredentialType,
			ObjectId:   credentialInfo.ObjectId,
			ObjectType: credentialInfo.ObjectType,
		}
		c.Set(base_web.CtxKeyAuthenticatedCredential, &credential)

		logger.Debugf("AuthenticatedCredential: %+v", credential)
		c.Next()
	}
}

func extractAPIKey(c *gin.Context, conf *KeyAuthConfig) (string, error) {
	if conf.KeyInAuthorizationHeader {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			return "", errors.New("missing Authorization header")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], conf.AuthorizationHeaderPrefix) {
			return "", errors.New("invalid Authorization header format")
		}

		return strings.TrimSpace(parts[1]), nil
	}

	if conf.KeyInHeader {
		authHeader := c.GetHeader(conf.KeyNames[0])
		if authHeader != "" {
			return strings.TrimSpace(authHeader), nil
		}
	}

	if conf.KeyInQuery {
		authHeader := c.Query(conf.KeyNames[0])
		if authHeader != "" {
			return strings.TrimSpace(authHeader), nil
		}

	}

	if conf.KeyInBody {
		var body map[string]any
		if err := c.BindJSON(&body); err != nil {
			return "", err
		}
		if k, ok := body[conf.KeyNames[0]]; ok {
			return k.(string), nil
		}
	}

	return "", fmt.Errorf("missing credential key")
}

func getCredential(apiKey string, conf *KeyAuthConfig) (*entities.AuthCredential, error) {
	sqlParams := map[string]any{
		"KEY":    apiKey,
		"ORG_ID": g_appConfig.Org.OrgId,
	}
	recs, err := base_db.DBQueryEnt[entities.AuthCredential](base_db.DB_CONN_NAME_DEFAULT, entities.DB_TABLE_AUTH_CREDENTIAL, "ORG_ID={ORG_ID} AND CREDENTIAL_TYPE='key' AND CREDENTIAL_VALUE={KEY}", sqlParams)
	if err != nil {
		logger.Error("DBQueryEnt fail: ", err)
		return nil, err
	}
	if len(recs) == 0 {
		return nil, errors.New("invalid api key")
	}

	return recs[0], nil
}
