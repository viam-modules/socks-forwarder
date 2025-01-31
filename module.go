// Package main defines a module that serves a generic service that interfaces with the
// socks-forwarder systemd service.
package main

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
	"go.viam.com/rdk/services/generic"
)

var socksForwarderControllerModel = resource.NewModel("viam", "ble-socks", "controller")

// socksForwarderController accepts `DoCommand` calls to start, stop and restart the socks
// forwarder systemd process.
type socksForwarderController struct {
	resource.Named
	resource.TriviallyReconfigurable
	resource.TriviallyCloseable
	resource.TriviallyValidateConfig

	logger logging.Logger
}

func newSocksForwarderController(ctx context.Context,
	_ resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (resource.Resource, error) {
	return &socksForwarderController{Named: conf.ResourceName().AsNamed(), logger: logger}, nil
}

// DoCommand can accept a "command" key with the values ["start", "stop" or "restart"].
// All other commands will result in an error. "start" idempotently starts the socks
// forwarder service. "stop" idempotently stops the socks forwarder service. "restart"
// forcibly restarts the socks forwarder service (stopping it if it is currently running.)
func (sfc *socksForwarderController) DoCommand(
	ctx context.Context,
	req map[string]interface{},
) (map[string]interface{}, error) {
	cmd, ok := req["command"]
	if !ok {
		return nil, errors.New("missing 'command' string")
	}
	switch cmd {
	case "start":
		sfc.logger.Info("Received request to start the socks-forwarder systemd service")
		if err := controlService(ctx, "start"); err != nil {
			sfc.logger.Errorw("Error starting the socks-forwarder systemd service", "error", err)
		}
		sfc.logger.Info("Started the socks-forwarder systemd service")
	case "stop":
		sfc.logger.Info("Stopping the socks-forwarder systemd service")
		if err := controlService(ctx, "stop"); err != nil {
			sfc.logger.Errorw("Error stopping the socks-forwarder systemd service", "error", err)
		}
		sfc.logger.Info("Stopped the socks-forwarder systemd service")
	case "restart":
		sfc.logger.Info("Restarting the socks-forwarder systemd service")
		if err := controlService(ctx, "restart"); err != nil {
			sfc.logger.Errorw("Error restarting the socks-forwarder systemd service", "error", err)
		}
		sfc.logger.Info("Restarted the socks-forwarder systemd service")
	default:
		return nil, fmt.Errorf("unknown 'command' %q", cmd)
	}

	return nil, nil
}

// Uses `systemctl` to put the `socks-forwarder` service into `state`.
func controlService(ctx context.Context, state string) error {
	cmd := exec.CommandContext(ctx, "systemctl", state, "socks-forwarder")
	return cmd.Run()
}

func main() {
	resource.RegisterService(generic.API, socksForwarderControllerModel,
		resource.Registration[resource.Resource, resource.NoNativeConfig]{
			Constructor: newSocksForwarderController,
		})
	module.ModularMain("socks-forwarder", resource.APIModel{generic.API, socksForwarderControllerModel})
}
