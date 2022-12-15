package test_structs

type TestConfig struct {
	Version string
}

func (t *TestConfig) GetVersion() string {
	return t.Version
}

func (t *TestConfig) GetPrivateKey() []byte {
	return nil
}
