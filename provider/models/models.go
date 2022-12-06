package models

import "cs-tf-provider/client"

type ProviderMeta struct {
	CSClient *client.CSClient
	Token    string
}

type AuthResponse struct {
	Token   *string `json:"Token,omitempty"`
	Code    *string `json:"Code,omitempty"`
	Message *string `json:"Message,omitempty"`
}
