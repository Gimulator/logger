package concluder

import client "github.com/Gimulator/client-go"

type Concluder interface {
	Send(client.Object) error
}
