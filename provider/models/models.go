package models

import "cs-tf-provider/client"

type ProviderMeta struct {
	CSClient *client.CSClient
	Token    string
}

type AuthResponse struct {
	Token string
}
