// Package main defines a module that serves a sensor that interfaces with the
// socks-forwarder systemd service.
package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"
	"go.viam.com/rdk/resource"
)

var socksForwarderSensorModel = resource.NewModel("viam", "socks", "forwarder")

// Keeps track of whether hciconfig is usable on the current machine. Will stop `Readings`
// from calling hciconfig if set to false.
var hciconfigUsable bool

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

func newSocksForwarderSensor(ctx context.Context,
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
		if err := controlService(ctx, "start"); err != nil {
			sfs.logger.Errorw("Error starting the socks-forwarder systemd service", "error", err)
		}
		sfs.logger.Info("Started the socks-forwarder systemd service")
	case "stop":
		sfs.logger.Info("Stopping the socks-forwarder systemd service...")
		if err := controlService(ctx, "stop"); err != nil {
			sfs.logger.Errorw("Error stopping the socks-forwarder systemd service", "error", err)
		}
		sfs.logger.Info("Stopped the socks-forwarder systemd service")
	case "restart":
		sfs.logger.Info("Restarting the socks-forwarder systemd service...")
		if err := controlService(ctx, "restart"); err != nil {
			sfs.logger.Errorw("Error restarting the socks-forwarder systemd service", "error", err)
		}
		sfs.logger.Info("Restarted the socks-forwarder systemd service")
	default:
		return nil, fmt.Errorf("unknown 'command' %q", cmd)
	}

	return nil, nil
}

// Uses `systemctl` to put the "socks-forwarder" service into `state`.
func controlService(ctx context.Context, state string) error {
	cmd := exec.CommandContext(ctx, "systemctl", state, "socks-forwarder")
	return cmd.Run()
}

/* Readings from hciconfig look like
hci0:	Type: Primary  Bus: UART
	BD Address: 14:D4:24:7F:75:A4  ACL MTU: 1021:8  SCO MTU: 64:1
	UP RUNNING PSCAN
	RX bytes:3602 acl:0 sco:0 events:375 errors:0
	TX bytes:62782 acl:0 sco:0 commands:375 errors:0
*/

// Readings returns the RX and TX byte totals of the hci0 (bluetooth)
// adapter.
func (sfs *socksForwarderSensor) Readings(
	ctx context.Context,
	_ map[string]interface{},
) (map[string]interface{}, error) {
	cmd := exec.CommandContext(ctx, "hciconfig")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		sfs.logger.Errorw("Error running 'hciconfig' to get readings. Will not fetch readings again", "error", err.Error())
		return nil, nil
	}

	hciconfigUsable = true // Assume hciconfig usable upon error.
	return map[string]interface{}{"hciconfig_output": out.String()}, nil
}

func main() {
	resource.RegisterComponent(sensor.API, socksForwarderSensorModel,
		resource.Registration[resource.Resource, resource.NoNativeConfig]{
			Constructor: newSocksForwarderSensor,
		})
	module.ModularMain("socks-forwarder-module", resource.APIModel{sensor.API, socksForwarderSensorModel})
}
