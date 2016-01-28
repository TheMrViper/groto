package groto

type (
	ClientRequest interface {
		// packet name
		Name() string
	}

	ClientResponse interface{}

	ServerRequest interface {
		// packet name
		Name() string

		// new pointer creator
		Create() ServerRequest

		// packet handler
		Handler() ServerResponse
	}

	ServerResponse interface{}

	RequestHeader struct {
		Action string
	}
	ResponseHeader struct {
		ErrorCode int
	}
)
