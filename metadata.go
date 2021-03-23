package rocketmq

type Metadata struct {
	// sdk proto (http or tcp),default tcp
	AccessProto string `json:"accessProto,omitempty"`

	// rocketmq Credentials
	AccessKey string `json:"accessKey,omitempty"`

	// rocketmq Credentials
	AccessSecret string `json:"accessSecret,omitempty"`

	// rocketmq's endpoint, optional, just for http proto
	Endpoint string `json:"endpoint,omitempty"`

	// rocketmq's name server, optional
	NameServer string `json:"nameServer,omitempty"`

	// consumer group for rocketmq's subscribers, suggested to provide
	ConsumerGroup string `json:"consumerGroup,omitempty"`

	// consumer group for rocketmq's subscribers, suggested to provide, just for http proto
	ConsumerBatchSize int `json:"consumerBatchSize,string,omitempty"`

	// consumer group for rocketmq's subscribers, suggested to provide, just for cgo proto
	ConsumerThreadNums int `json:"consumerThreadNums,string,omitempty"`

	// rocketmq's namespace, optional
	Namespace string `json:"namespace,omitempty"`

	// rocketmq's name server domain, optional
	NameServerDomain string `json:"nameServerDomain,omitempty"`

	// retry times to connect rocketmq's broker, optional
	Retries int `json:"retries,string,omitempty"`

	// topics to subscribe, use delimiter ',' to separate if more than one topics are configured, optional
	Topics string `json:"topics,omitempty"`
}
