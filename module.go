// Package main defines a module that serves a sensor that interfaces with the
// socks-forwarder systemd service.
package main

import (
	"context"
	"errors"
	"fmt"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

var socksForwarderSensorModel = resource.NewModel("viam", "sensor", "socks-forwarder-sensor")

// socksForwarderSensor reports the current tx/rx of hci0 (bluetooth adapter)
// through its `Readings`. It also accepts `DoCommand` calls to start, stop and
// restart the socks forwarder systemd process.
type socksForwarderSensor struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	resource.TriviallyValidateConfig

	logger logging.Logger
}

func newSocksForwarderSensor(_ context.Context,
	_ resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (resource.Resource, error) {
	return &socksForwarderSensor{Named: conf.ResourceName().AsNamed(), logger: logger}, nil
}

// DoCommand can accept a "command" key with the values ["start", "stop" or
// "restart"]. All other commands will result in error. "start" idempotently
// starts the socks forwarder service. "stop" idempotently stops the socks
// forwarder service. "restart" forcibly restarts the socks forwarder service
// (stopping it if it is currently running).
func (sfs *socksForwarderSensor) DoCommand(
	ctx context.Context,
	req map[string]interface{},
) (map[string]interface{}, error) {
	cmd, ok := req["command"]
	if !ok {
		return nil, errors.New("missing 'command' string")
	}
	switch cmd {
	case "start":
		sfs.logger.Info("Starting the socks-forwarder systemd service...")
		// TODO
		sfs.logger.Info("Started the socks-forwarder systemd service")
	case "stop":
		sfs.logger.Info("Stopping the socks-forwarder systemd service...")
		// TODO
		sfs.logger.Info("Stopped the socks-forwarder systemd service")
	case "restart":
		sfs.logger.Info("Restarting the socks-forwarder systemd service...")
		// TODO
		sfs.logger.Info("Restarted the socks-forwarder systemd service")
	default:
		return nil, fmt.Errorf("unknown 'command' %q", cmd)
	}

	return nil, nil
}

func (sfs *socksForwarderSensor) Readings(
	ctx context.Context,
	_ map[string]interface{},
) (map[string]interface{}, error) {
	readings := make(map[string]interface{})
	// TODO gather readings from syscall to `ifconfig`.
	return readings, nil
}

func main() {
	resource.RegisterComponent(sensor.API, socksForwarderSensorModel,
		resource.Registration[resource.Resource, resource.NoNativeConfig]{
			Constructor: newSocksForwarderSensor,
		})
	module.ModularMain("socks-forwarder-module", resource.APIModel{sensor.API, socksForwarderSensorModel})
}
