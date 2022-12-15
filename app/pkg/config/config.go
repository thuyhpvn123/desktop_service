package config

type IConfig interface {
	GetVersion() string
	GetPrivateKey() []byte
}
