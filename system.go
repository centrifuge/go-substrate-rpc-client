package substrate

type System struct {
	client Client
}

func NewSystemRPC(client Client) *System {
	return &System{client:client}
}
