package types

type ECDHConfig struct {
	UseDefaultBasepoint bool
	CustomBasepoint     []byte
}
