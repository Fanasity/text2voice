package model

type ServerCfg string

const CONFIG ServerCfg = "config"

type System struct {
	NodeId string `yaml:"nodeID"`
	Type   int    `yaml:"type"`
}

type KafkaConfig struct {
	Gateway                  []string `yaml:"gateway"`
	AutoCommitInterval       string   `yaml:"autoCommitInterval"`
	BatchSize                int64    `yaml:"batchSize"`
	BufferMemory             int64    `yaml:"bufferMemory"`
	MaxPollRecords           int64    `yaml:"maxPollRecords"`
	PropertiesSessionTimeout int64    `yaml:"propertiesSessionTimeout"`
	TaskTopic                string   `yaml:"taskTopic"`
	ResultTopic              string   `yaml:"resultTopic"`
}

type MinIOConfig struct {
	AccessKey string `yaml:"accessKey"`
	Bucket    string `yaml:"bucket"`
	Endpoint  string `yaml:"endpoint"`
	UseSecure bool   `yaml:"useSecure"`
	Public    string `yaml:"public"`
	Secret    string `yaml:"secret"`
}

type Config struct {
	System System      `yaml:"system"`
	Kafka  KafkaConfig `yaml:"kafka"`
	MinIO  MinIOConfig `yaml:"minio"`
}
