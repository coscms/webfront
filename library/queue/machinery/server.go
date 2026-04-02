package machinery

import (
	"fmt"
	"net/url"

	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"

	backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
	brokersiface "github.com/RichardKnop/machinery/v2/brokers/iface"
	locksiface "github.com/RichardKnop/machinery/v2/locks/iface"

	backendeager "github.com/RichardKnop/machinery/v2/backends/eager"
	brokereager "github.com/RichardKnop/machinery/v2/brokers/eager"
	lockeager "github.com/RichardKnop/machinery/v2/locks/eager"
)

func NewServer(tasks map[string]interface{}, configPaths ...interface{}) (*Server, error) {
	return NewServerWithCustom(nil, nil, nil, tasks, configPaths...)
}

func NewServerWithCustom(broker brokersiface.Broker, backend backendsiface.Backend, lock locksiface.Lock,
	tasks map[string]interface{}, configPaths ...interface{}) (*Server, error) {
	s := &Server{
		broker:  broker,
		backend: backend,
		lock:    lock,
	}
	err := s.ParseConfig(configPaths...)
	if err != nil {
		return s, err
	}
	err = s.initBackendEngine()
	if err != nil {
		return s, err
	}
	err = s.newServer()
	if err != nil {
		return s, err
	}
	if tasks != nil {
		err = s.RegisterTasks(tasks)
	}
	return s, err
}

type Server struct {
	config  *config.Config
	broker  brokersiface.Broker
	backend backendsiface.Backend
	lock    locksiface.Lock
	worker  *machinery.Worker
	*machinery.Server
}

func (s *Server) initBackendEngine() error {
	if s.broker == nil && len(s.config.Broker) > 0 {
		info, err := url.Parse(s.config.Broker)
		if err != nil {
			return err
		}
		f, y := brokerEngines[info.Scheme]
		if !y {
			return fmt.Errorf(`[queue.machinery] unsupported broker: %s`, info.Scheme)
		}
		s.broker, err = f(info, s.config)
		if err != nil {
			return fmt.Errorf(`[queue.machinery] failed to init broker %s: %w`, info.Scheme, err)
		}
	}
	if s.backend == nil && len(s.config.ResultBackend) > 0 {
		info, err := url.Parse(s.config.ResultBackend)
		if err != nil {
			return err
		}
		f, y := backendEngines[info.Scheme]
		if !y {
			return fmt.Errorf(`[queue.machinery] unsupported backend: %s`, info.Scheme)
		}
		s.backend, err = f(info, s.config)
		if err != nil {
			return fmt.Errorf(`[queue.machinery] failed to init backend %s: %w`, info.Scheme, err)
		}
	}
	if s.lock == nil && len(s.config.Lock) > 0 {
		info, err := url.Parse(s.config.Lock)
		if err != nil {
			return err
		}
		f, y := lockEngines[info.Scheme]
		if !y {
			return fmt.Errorf(`[queue.machinery] unsupported lock: %s`, info.Scheme)
		}
		s.lock, err = f(info, s.config)
		if err != nil {
			return fmt.Errorf(`[queue.machinery] failed to init lock %s: %w`, info.Scheme, err)
		}
	}
	return nil
}

func (s *Server) ParseConfig(configPaths ...interface{}) (err error) {
	var configPath string
	if len(configPaths) > 0 {
		switch c := configPaths[0].(type) {
		case string:
			configPath = c
		case config.Config:
			s.config = &c
			return
		case *config.Config:
			s.config = c
			return
		}
	}
	if len(configPath) > 0 {
		s.config, err = config.NewFromYaml(configPath, true)
		return
	}

	s.config, err = config.NewFromEnvironment()
	return
}

// newServer Create server instance
func (s *Server) newServer() (err error) {
	broker := s.broker
	backend := s.backend
	lock := s.lock
	if broker == nil {
		broker = brokereager.New()
	}
	if backend == nil {
		backend = backendeager.New()
	}
	if lock == nil {
		lock = lockeager.New()
	}
	s.Server = machinery.NewServer(s.config, broker, backend, lock)
	return
}

func (s *Server) Close() error {
	if s.worker != nil {
		s.worker.Quit()
	}
	return nil
}
