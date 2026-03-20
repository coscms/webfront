package machinery

import (
	"net/url"

	"github.com/RichardKnop/machinery/v2/config"

	backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
	brokersiface "github.com/RichardKnop/machinery/v2/brokers/iface"
	locksiface "github.com/RichardKnop/machinery/v2/locks/iface"
)

var (
	brokerEngines  = map[string]func(*url.URL, *config.Config) (brokersiface.Broker, error){}
	backendEngines = map[string]func(*url.URL, *config.Config) (backendsiface.Backend, error){}
	lockEngines    = map[string]func(*url.URL, *config.Config) (locksiface.Lock, error){}
)

func RegisterBrokerEngine(name string, f func(*url.URL, *config.Config) (brokersiface.Broker, error)) {
	brokerEngines[name] = f
}

func RegisterBackendEngine(name string, f func(*url.URL, *config.Config) (backendsiface.Backend, error)) {
	backendEngines[name] = f
}

func RegisterLockEngine(name string, f func(*url.URL, *config.Config) (locksiface.Lock, error)) {
	lockEngines[name] = f
}
