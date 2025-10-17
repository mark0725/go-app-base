package middleware

type OIDCConfig struct {
	Anonymous                                    string       `json:"anonymous,omitempty"`                                         // 失败鉴权时使用的“匿名”consumer
	Audience                                     []string     `json:"audience,omitempty"`                                          // 授权端点 audience
	AudienceClaim                                []string     `json:"audience_claim,omitempty"`                                    // audience 所在 claim
	AudienceRequired                             []string     `json:"audience_required,omitempty"`                                 // 授权成功所需 audience
	AuthMethods                                  []string     `json:"auth_methods,omitempty"`                                      // 启用的 grant / credential 类型
	AuthenticatedGroupsClaim                     []string     `json:"authenticated_groups_claim,omitempty"`                        // 已认证组 claim
	AuthorizationCookieDomain                    string       `json:"authorization_cookie_domain,omitempty"`                       // 授权 Cookie Domain
	AuthorizationCookieHTTPOnly                  *bool        `json:"authorization_cookie_http_only,omitempty"`                    // 授权 Cookie HttpOnly
	AuthorizationCookieName                      string       `json:"authorization_cookie_name,omitempty"`                         // 授权 Cookie 名
	AuthorizationCookiePath                      string       `json:"authorization_cookie_path,omitempty"`                         // 授权 Cookie Path
	AuthorizationCookieSameSite                  string       `json:"authorization_cookie_same_site,omitempty"`                    // 授权 Cookie SameSite
	AuthorizationCookieSecure                    *bool        `json:"authorization_cookie_secure,omitempty"`                       // 授权 Cookie Secure
	AuthorizationEndpoint                        string       `json:"authorization_endpoint,omitempty"`                            // 授权端点
	AuthorizationQueryArgsClient                 []string     `json:"authorization_query_args_client,omitempty"`                   // 透传到授权端点的客户端查询参数
	AuthorizationQueryArgsNames                  []string     `json:"authorization_query_args_names,omitempty"`                    // 透传查询参数名
	AuthorizationQueryArgsValues                 []string     `json:"authorization_query_args_values,omitempty"`                   // 透传查询参数值
	AuthorizationRollingTimeout                  *int         `json:"authorization_rolling_timeout,omitempty"`                     // 授权码流 Session 滚动超时
	BearerTokenCookieName                        string       `json:"bearer_token_cookie_name,omitempty"`                          // Bearer-Token Cookie 名
	BearerTokenParamType                         []string     `json:"bearer_token_param_type,omitempty"`                           // Bearer-Token 查找位置
	ByUsernameIgnoreCase                         *bool        `json:"by_username_ignore_case,omitempty"`                           // 用户名匹配是否忽略大小写
	CacheIntrospection                           *bool        `json:"cache_introspection,omitempty"`                               // 缓存 introspection
	CacheTokenExchange                           *bool        `json:"cache_token_exchange,omitempty"`                              // 缓存 token-exchange
	CacheTokens                                  *bool        `json:"cache_tokens,omitempty"`                                      // 缓存 token
	CacheTokensSalt                              string       `json:"cache_tokens_salt,omitempty"`                                 // token 缓存盐
	CacheTTL                                     *int         `json:"cache_ttl,omitempty"`                                         // 默认缓存 ttl
	CacheTTLMax                                  *int         `json:"cache_ttl_max,omitempty"`                                     // 最大缓存 ttl
	CacheTTLMin                                  *int         `json:"cache_ttl_min,omitempty"`                                     // 最小缓存 ttl
	CacheTTLNeg                                  *int         `json:"cache_ttl_neg,omitempty"`                                     // 负缓存 ttl
	CacheTTLResurrect                            *int         `json:"cache_ttl_resurrect,omitempty"`                               // 复活 ttl
	CacheUserInfo                                *bool        `json:"cache_user_info,omitempty"`                                   // 缓存 userinfo
	ClaimsForbidden                              []string     `json:"claims_forbidden,omitempty"`                                  // 禁用的 claim
	ClientAlg                                    []string     `json:"client_alg,omitempty"`                                        // 客户端 JWT 算法
	ClientArg                                    string       `json:"client_arg,omitempty"`                                        // 选取 client 的参数名
	ClientAuth                                   []string     `json:"client_auth,omitempty"`                                       // token/introspection 端点认证方式
	ClientCredentialsParamType                   []string     `json:"client_credentials_param_type,omitempty"`                     // client_credential 查找位置
	ClientID                                     []string     `json:"client_id,omitempty"`                                         // client_id
	ClientJWK                                    []JWK        `json:"client_jwk,omitempty"`                                        // private_key_jwt 使用的 JWK
	ClientSecret                                 []string     `json:"client_secret,omitempty"`                                     // client_secret
	ClusterCacheRedis                            *RedisConfig `json:"cluster_cache_redis,omitempty"`                               // 共享缓存 Redis
	ClusterCacheStrategy                         string       `json:"cluster_cache_strategy,omitempty"`                            // 集群缓存策略
	ConsumerBy                                   []string     `json:"consumer_by,omitempty"`                                       // consumer 映射字段
	ConsumerClaim                                []string     `json:"consumer_claim,omitempty"`                                    // consumer 映射 claim
	ConsumerOptional                             *bool        `json:"consumer_optional,omitempty"`                                 // consumer 映射失败是否继续
	CredentialClaim                              []string     `json:"credential_claim,omitempty"`                                  // 虚拟凭据 claim
	DisableSession                               []string     `json:"disable_session,omitempty"`                                   // 禁用颁发 session 的 grants
	DiscoveryHeadersNames                        []string     `json:"discovery_headers_names,omitempty"`                           // discovery 额外 header 名
	DiscoveryHeadersValues                       []string     `json:"discovery_headers_values,omitempty"`                          // discovery 额外 header 值
	DisplayErrors                                *bool        `json:"display_errors,omitempty"`                                    // 失败时返回详细错误
	Domains                                      []string     `json:"domains,omitempty"`                                           // 允许的 hd
	DownstreamAccessTokenHeader                  string       `json:"downstream_access_token_header,omitempty"`                    // 下游 access-token Header
	DownstreamAccessTokenJWKHeader               string       `json:"downstream_access_token_jwk_header,omitempty"`                // 下游 access-token JWK Header
	DownstreamHeadersClaims                      []string     `json:"downstream_headers_claims,omitempty"`                         // 下游 Header claim 路径
	DownstreamHeadersNames                       []string     `json:"downstream_headers_names,omitempty"`                          // 下游 Header 名
	DownstreamIDTokenHeader                      string       `json:"downstream_id_token_header,omitempty"`                        // 下游 id-token Header
	DownstreamIDTokenJWKHeader                   string       `json:"downstream_id_token_jwk_header,omitempty"`                    // 下游 id-token JWK Header
	DownstreamIntrospectionHeader                string       `json:"downstream_introspection_header,omitempty"`                   // 下游 introspection Header
	DownstreamIntrospectionJWTHeader             string       `json:"downstream_introspection_jwt_header,omitempty"`               // 下游 introspection JWT Header
	DownstreamRefreshTokenHeader                 string       `json:"downstream_refresh_token_header,omitempty"`                   // 下游 refresh-token Header
	DownstreamSessionIDHeader                    string       `json:"downstream_session_id_header,omitempty"`                      // 下游 session-id Header
	DownstreamUserInfoHeader                     string       `json:"downstream_user_info_header,omitempty"`                       // 下游 userinfo Header
	DownstreamUserInfoJWTHeader                  string       `json:"downstream_user_info_jwt_header,omitempty"`                   // 下游 userinfo JWT Header
	DPoPProofLifetime                            *int         `json:"dpop_proof_lifetime,omitempty"`                               // DPoP proof 生命周期
	DPoPUseNonce                                 *bool        `json:"dpop_use_nonce,omitempty"`                                    // DPoP 是否使用 nonce
	EnableHSSignatures                           *bool        `json:"enable_hs_signatures,omitempty"`                              // 允许 HS*
	EndSessionEndpoint                           string       `json:"end_session_endpoint,omitempty"`                              // 结束会话端点
	ExposeErrorCode                              *bool        `json:"expose_error_code,omitempty"`                                 // 返回错误码 header
	ExtraJWKSURIs                                []string     `json:"extra_jwks_uris,omitempty"`                                   // 额外 jwks_uri
	ForbiddenDestroySession                      *bool        `json:"forbidden_destroy_session,omitempty"`                         // 403 时销毁 session
	ForbiddenErrorMessage                        string       `json:"forbidden_error_message,omitempty"`                           // 403 错误信息
	ForbiddenRedirectURI                         []string     `json:"forbidden_redirect_uri,omitempty"`                            // 403 跳转 URI
	GroupsClaim                                  []string     `json:"groups_claim,omitempty"`                                      // groups claim
	GroupsRequired                               []string     `json:"groups_required,omitempty"`                                   // 所需 groups
	HideCredentials                              *bool        `json:"hide_credentials,omitempty"`                                  // 移除凭据
	HTTPProxy                                    string       `json:"http_proxy,omitempty"`                                        // HTTP 代理
	HTTPProxyAuthorization                       string       `json:"http_proxy_authorization,omitempty"`                          // HTTP 代理认证
	HTTPVersion                                  *float64     `json:"http_version,omitempty"`                                      // 请求使用的 HTTP 版本
	HTTPSProxy                                   string       `json:"https_proxy,omitempty"`                                       // HTTPS 代理
	HTTPSProxyAuthorization                      string       `json:"https_proxy_authorization,omitempty"`                         // HTTPS 代理认证
	IDTokenParamName                             string       `json:"id_token_param_name,omitempty"`                               // id-token 参数名
	IDTokenParamType                             []string     `json:"id_token_param_type,omitempty"`                               // id-token 查找位置
	IgnoreSignature                              []string     `json:"ignore_signature,omitempty"`                                  // 跳过签名验证的 grant
	IntrospectJWTTokens                          *bool        `json:"introspect_jwt_tokens,omitempty"`                             // introspect JWT access-token
	IntrospectionAccept                          string       `json:"introspection_accept,omitempty"`                              // introspection Accept
	IntrospectionCheckActive                     *bool        `json:"introspection_check_active,omitempty"`                        // 检查 active
	IntrospectionEndpoint                        string       `json:"introspection_endpoint,omitempty"`                            // introspection 端点
	IntrospectionEndpointAuthMethod              string       `json:"introspection_endpoint_auth_method,omitempty"`                // introspection 认证方式
	IntrospectionHeadersClient                   []string     `json:"introspection_headers_client,omitempty"`                      // 客户端透传 header
	IntrospectionHeadersNames                    []string     `json:"introspection_headers_names,omitempty"`                       // 额外 header 名
	IntrospectionHeadersValues                   []string     `json:"introspection_headers_values,omitempty"`                      // 额外 header 值
	IntrospectionHint                            string       `json:"introspection_hint,omitempty"`                                // introspection hint
	IntrospectionPostArgsClient                  []string     `json:"introspection_post_args_client,omitempty"`                    // 客户端透传 post 参数
	IntrospectionPostArgsClientHeaders           []string     `json:"introspection_post_args_client_headers,omitempty"`            // 从 Header 取的 post 参数
	IntrospectionPostArgsNames                   []string     `json:"introspection_post_args_names,omitempty"`                     // post 参数名
	IntrospectionPostArgsValues                  []string     `json:"introspection_post_args_values,omitempty"`                    // post 参数值
	IntrospectionTokenParamName                  string       `json:"introspection_token_param_name,omitempty"`                    // introspection token 参数名
	Issuer                                       string       `json:"issuer"`                                                      // discovery 地址 / Issuer
	IssuersAllowed                               []string     `json:"issuers_allowed,omitempty"`                                   // token 中允许的 iss
	JWTSessClaim                                 string       `json:"jwt_session_claim,omitempty"`                                 // JWT session claim
	JWTSessCookie                                string       `json:"jwt_session_cookie,omitempty"`                                // JWT session Cookie 名
	Keepalive                                    *bool        `json:"keepalive,omitempty"`                                         // HTTP Keep-Alive
	Leeway                                       *int         `json:"leeway,omitempty"`                                            // 时间字段容差
	LoginAction                                  string       `json:"login_action,omitempty"`                                      // 登录后行为
	LoginMethods                                 []string     `json:"login_methods,omitempty"`                                     // 启用登录的 grants
	LoginRedirectMode                            string       `json:"login_redirect_mode,omitempty"`                               // redirect 模式
	LoginRedirectURI                             []string     `json:"login_redirect_uri,omitempty"`                                // 登录成功跳转 URI
	LoginTokens                                  []string     `json:"login_tokens,omitempty"`                                      // 登录时返回 Token 列表
	LogoutMethods                                []string     `json:"logout_methods,omitempty"`                                    // 激活登出的 HTTP 方法
	LogoutPostArg                                string       `json:"logout_post_arg,omitempty"`                                   // body 参数触发登出
	LogoutQueryArg                               string       `json:"logout_query_arg,omitempty"`                                  // query 参数触发登出
	LogoutRedirectURI                            []string     `json:"logout_redirect_uri,omitempty"`                               // 登出跳转 URI
	LogoutRevoke                                 *bool        `json:"logout_revoke,omitempty"`                                     // 登出时撤销 token
	LogoutRevokeAccessToken                      *bool        `json:"logout_revoke_access_token,omitempty"`                        // 撤销 access_token
	LogoutRevokeRefreshToken                     *bool        `json:"logout_revoke_refresh_token,omitempty"`                       // 撤销 refresh_token
	LogoutURISuffix                              string       `json:"logout_uri_suffix,omitempty"`                                 // 登出请求 URI 后缀
	MaxAge                                       *int         `json:"max_age,omitempty"`                                           // 与 auth_time 比较的最大秒数
	MTLSIntrospectionEndpoint                    string       `json:"mtls_introspection_endpoint,omitempty"`                       // mTLS introspection 端点别名
	MTLSRevocationEndpoint                       string       `json:"mtls_revocation_endpoint,omitempty"`                          // mTLS revocation 端点别名
	MTLSTokenEndpoint                            string       `json:"mtls_token_endpoint,omitempty"`                               // mTLS token 端点别名
	NoProxy                                      string       `json:"no_proxy,omitempty"`                                          // 不经代理的主机
	PasswordParamType                            []string     `json:"password_param_type,omitempty"`                               // 用户名/密码查找位置
	PreserveQueryArgs                            *bool        `json:"preserve_query_args,omitempty"`                               // AuthCode 流保留查询串
	ProofOfPossessionAuthMethodsValidation       *bool        `json:"proof_of_possession_auth_methods_validation,omitempty"`       // PoP 时校验 auth_methods
	ProofOfPossessionDPoP                        string       `json:"proof_of_possession_dpop,omitempty"`                          // DPoP 验证级别
	ProofOfPossessionMTLS                        string       `json:"proof_of_possession_mtls,omitempty"`                          // mTLS PoP 验证级别
	PushedAuthorizationRequestEndpoint           string       `json:"pushed_authorization_request_endpoint,omitempty"`             // PAR 端点
	PushedAuthorizationRequestEndpointAuthMethod string       `json:"pushed_authorization_request_endpoint_auth_method,omitempty"` // PAR 认证方式
	RedirectURI                                  []string     `json:"redirect_uri,omitempty"`                                      // redirect_uri 白名单
	Redis                                        *RedisConfig `json:"redis,omitempty"`                                             // Session Redis
	RediscoveryLifetime                          *int         `json:"rediscovery_lifetime,omitempty"`                              // 发现重试间隔
	RefreshTokenParamName                        string       `json:"refresh_token_param_name,omitempty"`                          // refresh_token 参数名
	RefreshTokenParamType                        []string     `json:"refresh_token_param_type,omitempty"`                          // refresh_token 查找位置
	RefreshTokens                                *bool        `json:"refresh_tokens,omitempty"`                                    // 自动刷新 token
	RequireProofKeyForCodeExchange               *bool        `json:"require_proof_key_for_code_exchange,omitempty"`               // 强制 PKCE
	RequirePushedAuthorizationRequests           *bool        `json:"require_pushed_authorization_requests,omitempty"`             // 强制 PAR
	RequireSignedRequestObject                   *bool        `json:"require_signed_request_object,omitempty"`                     // 强制 JAR
	ResolveDistributedClaims                     *bool        `json:"resolve_distributed_claims,omitempty"`                        // 解析分布式 claim
	ResponseMode                                 string       `json:"response_mode,omitempty"`                                     // authorization response_mode
	ResponseType                                 []string     `json:"response_type,omitempty"`                                     // authorization response_type
	Reverify                                     *bool        `json:"reverify,omitempty"`                                          // 每次使用 session 时重新验证 token
	RevocationEndpoint                           string       `json:"revocation_endpoint,omitempty"`                               // 撤销端点
	RevocationEndpointAuthMethod                 string       `json:"revocation_endpoint_auth_method,omitempty"`                   // 撤销端点认证方式
	RevocationTokenParamName                     string       `json:"revocation_token_param_name,omitempty"`                       // revocation token 参数
	RolesClaim                                   []string     `json:"roles_claim,omitempty"`                                       // roles claim
	RolesRequired                                []string     `json:"roles_required,omitempty"`                                    // 所需 roles
	RunOnPreflight                               *bool        `json:"run_on_preflight,omitempty"`                                  // 是否处理 OPTIONS
	Scopes                                       []string     `json:"scopes,omitempty"`                                            // 请求 scope
	ScopesClaim                                  []string     `json:"scopes_claim,omitempty"`                                      // scopes claim
	ScopesRequired                               []string     `json:"scopes_required,omitempty"`                                   // 所需 scopes
	SearchUserInfo                               *bool        `json:"search_user_info,omitempty"`                                  // 使用 userinfo 查询补充
	SessionAbsoluteTimeout                       *int         `json:"session_absolute_timeout,omitempty"`                          // Session 绝对过期
	SessionAudience                              string       `json:"session_audience,omitempty"`                                  // Session audience
	SessionCookieDomain                          string       `json:"session_cookie_domain,omitempty"`                             // Session Cookie Domain
	SessionCookieHTTPOnly                        *bool        `json:"session_cookie_http_only,omitempty"`                          // Session Cookie HttpOnly
	SessionCookieName                            string       `json:"session_cookie_name,omitempty"`                               // Session Cookie 名
	SessionCookiePath                            string       `json:"session_cookie_path,omitempty"`                               // Session Cookie Path
	SessionCookieSameSite                        string       `json:"session_cookie_same_site,omitempty"`                          // Session Cookie SameSite
	SessionCookieSecure                          *bool        `json:"session_cookie_secure,omitempty"`                             // Session Cookie Secure
	SessionEnforceSameSubject                    *bool        `json:"session_enforce_same_subject,omitempty"`                      // 同一 subject 限制
	SessionHashStorageKey                        *bool        `json:"session_hash_storage_key,omitempty"`                          // 存储 key 哈希
	SessionHashSubject                           *bool        `json:"session_hash_subject,omitempty"`                              // 存储 subject 哈希
	SessionIdlingTimeout                         *int         `json:"session_idling_timeout,omitempty"`                            // Session 空闲超时
	SessionMemcachedHost                         string       `json:"session_memcached_host,omitempty"`                            // Memcached host
	SessionMemcachedPort                         *int         `json:"session_memcached_port,omitempty"`                            // Memcached port
	SessionMemcachedPrefix                       string       `json:"session_memcached_prefix,omitempty"`                          // Memcached 前缀
	SessionMemcachedSocket                       string       `json:"session_memcached_socket,omitempty"`                          // Memcached unix socket
	SessionRemember                              *bool        `json:"session_remember,omitempty"`                                  // Persistent Session
	SessionRememberAbsoluteTimeout               *int         `json:"session_remember_absolute_timeout,omitempty"`                 // Persistent 绝对超时
	SessionRememberCookieName                    string       `json:"session_remember_cookie_name,omitempty"`                      // Persistent Cookie 名
	SessionRememberRollingTimeout                *int         `json:"session_remember_rolling_timeout,omitempty"`                  // Persistent 滚动超时
	SessionRequestHeaders                        []string     `json:"session_request_headers,omitempty"`                           // 向上游发送的 Session Header
	SessionResponseHeaders                       []string     `json:"session_response_headers,omitempty"`                          // 向下游发送的 Session Header
	SessionRollingTimeout                        *int         `json:"session_rolling_timeout,omitempty"`                           // Session 滚动超时
	SessionSecret                                string       `json:"session_secret,omitempty"`                                    // Session 加密密钥
	SessionStorage                               string       `json:"session_storage,omitempty"`                                   // Session 存储方式
	SessionStoreMetadata                         *bool        `json:"session_store_metadata,omitempty"`                            // 存储 Session 元数据
	SSLVerify                                    *bool        `json:"ssl_verify,omitempty"`                                        // 验证 IdP 证书
	Timeout                                      *int         `json:"timeout,omitempty"`                                           // 网络超时 (ms)
	TLSClientAuthCertID                          string       `json:"tls_client_auth_cert_id,omitempty"`                           // 用于 mTLS 的客户端证书 ID
	TLSClientAuthSSLVerify                       *bool        `json:"tls_client_auth_ssl_verify,omitempty"`                        // mTLS 校验证书
	TokenCacheKeyIncludeScope                    *bool        `json:"token_cache_key_include_scope,omitempty"`                     // 缓存 key 包含 scope
	TokenEndpoint                                string       `json:"token_endpoint,omitempty"`                                    // token 端点
	TokenEndpointAuthMethod                      string       `json:"token_endpoint_auth_method,omitempty"`                        // token 端点认证方式
	TokenExchangeEndpoint                        string       `json:"token_exchange_endpoint,omitempty"`                           // token-exchange 端点
	TokenHeadersClient                           []string     `json:"token_headers_client,omitempty"`                              // 客户端透传 header
	TokenHeadersGrants                           []string     `json:"token_headers_grants,omitempty"`                              // 仅在指定 grants 返回 header
	TokenHeadersNames                            []string     `json:"token_headers_names,omitempty"`                               // 额外 header 名
	TokenHeadersPrefix                           string       `json:"token_headers_prefix,omitempty"`                              // 下游 header 前缀
	TokenHeadersReplay                           []string     `json:"token_headers_replay,omitempty"`                              // 透传给下游的 header
	TokenHeadersValues                           []string     `json:"token_headers_values,omitempty"`                              // 额外 header 值
	TokenPostArgsClient                          []string     `json:"token_post_args_client,omitempty"`                            // 客户端透传 post 参数
	TokenPostArgsNames                           []string     `json:"token_post_args_names,omitempty"`                             // 额外 post 参数名
	TokenPostArgsValues                          []string     `json:"token_post_args_values,omitempty"`                            // 额外 post 参数值
	UnauthorizedDestroySession                   *bool        `json:"unauthorized_destroy_session,omitempty"`                      // 401 时销毁 session
	UnauthorizedErrorMessage                     string       `json:"unauthorized_error_message,omitempty"`                        // 401 错误信息
	UnauthorizedRedirectURI                      []string     `json:"unauthorized_redirect_uri,omitempty"`                         // 401 跳转 URI
	UnexpectedRedirectURI                        []string     `json:"unexpected_redirect_uri,omitempty"`                           // 意外状态跳转 URI
	UpstreamAccessTokenHeader                    string       `json:"upstream_access_token_header,omitempty"`                      // 上游 access-token Header
	UpstreamAccessTokenJWKHeader                 string       `json:"upstream_access_token_jwk_header,omitempty"`                  // 上游 access-token JWK Header
	UpstreamHeadersClaims                        []string     `json:"upstream_headers_claims,omitempty"`                           // 上游 Header claim (顶层)
	UpstreamHeadersNames                         []string     `json:"upstream_headers_names,omitempty"`                            // 上游 Header 名
	UpstreamIDTokenHeader                        string       `json:"upstream_id_token_header,omitempty"`                          // 上游 id-token Header
	UpstreamIDTokenJWKHeader                     string       `json:"upstream_id_token_jwk_header,omitempty"`                      // 上游 id-token JWK Header
	UpstreamIntrospectionHeader                  string       `json:"upstream_introspection_header,omitempty"`                     // 上游 introspection Header
	UpstreamIntrospectionJWTHeader               string       `json:"upstream_introspection_jwt_header,omitempty"`                 // 上游 introspection JWT Header
	UpstreamRefreshTokenHeader                   string       `json:"upstream_refresh_token_header,omitempty"`                     // 上游 refresh-token Header
	UpstreamSessionIDHeader                      string       `json:"upstream_session_id_header,omitempty"`                        // 上游 session-id Header
	UpstreamUserInfoHeader                       string       `json:"upstream_user_info_header,omitempty"`                         // 上游 userinfo Header
	UpstreamUserInfoJWTHeader                    string       `json:"upstream_user_info_jwt_header,omitempty"`                     // 上游 userinfo JWT Header
	UserInfoAccept                               string       `json:"userinfo_accept,omitempty"`                                   // userinfo Accept
	UserInfoEndpoint                             string       `json:"userinfo_endpoint,omitempty"`                                 // userinfo 端点
	UserInfoHeadersClient                        []string     `json:"userinfo_headers_client,omitempty"`                           // 客户端透传 header
	UserInfoHeadersNames                         []string     `json:"userinfo_headers_names,omitempty"`                            // 额外 header 名
	UserInfoHeadersValues                        []string     `json:"userinfo_headers_values,omitempty"`                           // 额外 header 值
	UserInfoQueryArgsClient                      []string     `json:"userinfo_query_args_client,omitempty"`                        // 客户端透传 query
	UserInfoQueryArgsNames                       []string     `json:"userinfo_query_args_names,omitempty"`                         // 额外 query 名
	UserInfoQueryArgsValues                      []string     `json:"userinfo_query_args_values,omitempty"`                        // 额外 query 值
	UsingPseudoIssuer                            *bool        `json:"using_pseudo_issuer,omitempty"`                               // 伪 Issuer
	VerifyClaims                                 *bool        `json:"verify_claims,omitempty"`                                     // 校验标准 claim
	VerifyNonce                                  *bool        `json:"verify_nonce,omitempty"`                                      // 校验 nonce
	VerifyParameters                             *bool        `json:"verify_parameters,omitempty"`                                 // 校验配置与 discovery
	VerifySignature                              *bool        `json:"verify_signature,omitempty"`                                  // 校验签名
}

// JWK 定义（仅列出常用字段，若需全部字段可再扩展）
type JWK struct {
	Alg     string   `json:"alg,omitempty"`      // 算法
	Crv     string   `json:"crv,omitempty"`      // Curve
	D       string   `json:"d,omitempty"`        // 私钥 d
	DP      string   `json:"dp,omitempty"`       // dp
	DQ      string   `json:"dq,omitempty"`       // dq
	E       string   `json:"e,omitempty"`        // 公钥 e
	K       string   `json:"k,omitempty"`        // 对称 key
	Kid     string   `json:"kid,omitempty"`      // key id
	Kty     string   `json:"kty,omitempty"`      // key type
	N       string   `json:"n,omitempty"`        // 公钥 n
	P       string   `json:"p,omitempty"`        // p
	Q       string   `json:"q,omitempty"`        // q
	Qi      string   `json:"qi,omitempty"`       // qi
	R       string   `json:"r,omitempty"`        // r
	X       string   `json:"x,omitempty"`        // x
	Y       string   `json:"y,omitempty"`        // y
	X5c     []string `json:"x5c,omitempty"`      // x5c
	X5t     string   `json:"x5t,omitempty"`      // x5t
	X5tS256 string   `json:"x5t#S256,omitempty"` // x5t#S256
	X5u     string   `json:"x5u,omitempty"`      // x5u
	Use     string   `json:"use,omitempty"`      // 用途
	KeyOps  []string `json:"key_ops,omitempty"`  // key 操作
	Issuer  string   `json:"issuer,omitempty"`   // 发行者
}

// Redis 结点（用于 cluster_nodes / sentinel_nodes）
type RedisNode struct {
	Host string `json:"host,omitempty"` // 主机
	Port int    `json:"port,omitempty"` // 端口
}

// 通用 Redis / Redis-Cluster / Sentinel 配置
type RedisConfig struct {
	ClusterMaxRedirections int         `json:"cluster_max_redirections,omitempty"` // Cluster 重定向重试
	ClusterNodes           []RedisNode `json:"cluster_nodes,omitempty"`            // Cluster 节点
	ConnectTimeout         int         `json:"connect_timeout,omitempty"`          // 连接超时 ms
	ConnectionIsProxied    *bool       `json:"connection_is_proxied,omitempty"`    // 经过代理
	Database               int         `json:"database,omitempty"`                 // DB
	Host                   string      `json:"host,omitempty"`                     // Host
	KeepaliveBacklog       int         `json:"keepalive_backlog,omitempty"`        // keepalive backlog
	KeepalivePoolSize      int         `json:"keepalive_pool_size,omitempty"`      // keepalive pool
	Password               string      `json:"password,omitempty"`                 // AUTH 密码
	Port                   int         `json:"port,omitempty"`                     // 端口
	Prefix                 string      `json:"prefix,omitempty"`                   // key 前缀
	ReadTimeout            int         `json:"read_timeout,omitempty"`             // 读超时
	SendTimeout            int         `json:"send_timeout,omitempty"`             // 写超时
	SentinelMaster         string      `json:"sentinel_master,omitempty"`          // Sentinel master
	SentinelNodes          []RedisNode `json:"sentinel_nodes,omitempty"`           // Sentinel 节点
	SentinelPassword       string      `json:"sentinel_password,omitempty"`        // Sentinel 密码
	SentinelRole           string      `json:"sentinel_role,omitempty"`            // Sentinel 角色
	SentinelUsername       string      `json:"sentinel_username,omitempty"`        // Sentinel 用户名
	ServerName             string      `json:"server_name,omitempty"`              // SNI
	Socket                 string      `json:"socket,omitempty"`                   // Unix Socket
	SSL                    *bool       `json:"ssl,omitempty"`                      // SSL
	SSLVerify              *bool       `json:"ssl_verify,omitempty"`               // 校验证书
	Username               string      `json:"username,omitempty"`                 // ACL 用户名
}
