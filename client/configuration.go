package client

// Configuration stores the configuration for a Client
type Configuration struct {
	URL             string
	AccessKeyID     string
	SecretAccessKey string
	AWSServiceName  string
	Region          string
	Login           *Login
}

// NewConfiguration creates a default Configuration struct
func NewConfiguration() *Configuration {
	cfg := &Configuration{
		AWSServiceName: "s3",
	}

	return cfg
}

type Login struct {
	Username     string
	Password     string
	ParentUserID string `json:"ParentUserID,omitempty"`
}
