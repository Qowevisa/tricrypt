package communication

type ClientForServer struct {
	ID       uint16
	Nickname string
}

type RegisteredUser struct {
	ID uint16
}

const (
	LINK_STATUS_CREATED = 1 + iota
	LINK_STATUS_EXPIRED
)

const (
	LINK_LEN_V1 = 32
)

// NOTE: Data should be 32 or 64 bytes
type Link struct {
	Status   uint8
	Data     []byte
	UseCount uint32
}
