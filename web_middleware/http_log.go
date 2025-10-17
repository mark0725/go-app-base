package middleware

// Config 定义了插件的全部可配置项
type HttpLogConfig struct {
	ContentType       string                 `json:"content_type"`                   // 发送的数据类型，默认 application/json
	CustomFieldsByLua map[string]interface{} `json:"custom_fields_by_lua,omitempty"` // 以 Lua 代码生成的自定义字段
	FlushTimeout      float64                `json:"flush_timeout,omitempty"`        // 当 queue_size > 1 时，最大发送间隔（秒）
	Headers           map[string]string      `json:"headers,omitempty"`              // 附加到上游请求的自定义 Header
	HTTPEndpoint      string                 `json:"http_endpoint"`                  // 上游 HTTP 端点，必填
	Keepalive         int64                  `json:"keepalive,omitempty"`            // 空闲连接保活时长（毫秒），默认 60000
	Method            string                 `json:"method,omitempty"`               // HTTP 方法，支持 POST/PUT/PATCH，默认 POST
	Queue             HttpLogQueueConfig     `json:"queue,omitempty"`                // 内部异步队列配置
	QueueSize         int                    `json:"queue_size,omitempty"`           // 每次发送的最大日志条目数
	RetryCount        int                    `json:"retry_count,omitempty"`          // 发送失败时的最大重试次数
	Timeout           int64                  `json:"timeout,omitempty"`              // 发送请求的超时时间（毫秒），默认 10000
}

// QueueConfig 定义了异步队列相关的详细行为
type HttpLogQueueConfig struct {
	ConcurrencyLimit   int     `json:"concurrency_limit,omitempty"`    // 并发定时器数量，-1 表示无限制
	InitialRetryDelay  float64 `json:"initial_retry_delay,omitempty"`  // 首次重试前的等待时间（秒）
	MaxBatchSize       int     `json:"max_batch_size,omitempty"`       // 单批处理的最大条目数
	MaxBytes           int     `json:"max_bytes,omitempty"`            // 队列中允许等待的最大字节数
	MaxCoalescingDelay float64 `json:"max_coalescing_delay,omitempty"` // 首条进入队列后允许的最大等待时间（秒）
	MaxEntries         int     `json:"max_entries,omitempty"`          // 队列中等待的最大条目数
	MaxRetryDelay      float64 `json:"max_retry_delay,omitempty"`      // 指数退避的最大间隔（秒）
	MaxRetryTime       float64 `json:"max_retry_time,omitempty"`       // 在放弃之前的最大重试时长（秒）
}
