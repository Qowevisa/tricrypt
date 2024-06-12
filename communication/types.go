package communication

type ClientForServer struct {
	ID       uint16
	Nickname string
}

type RegisteredUser struct {
	IsRegistered bool
	ID           uint16
	Name         string
}

const (
	LINK_STATUS_CREATED = 1 + iota
	LINK_STATUS_EXPIRED
)

const (
	LINK_LEN_V1 = 8
)

// NOTE: Data should be 32 or 64 bytes
type Link struct {
	Status   uint8
	Data     []byte
	UseCount uint16
}
