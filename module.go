package main

import (
	"context"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/module"

	"go.viam.com/utils"
)

func main() {
	utils.ContextualMain(mainWithArgs, module.NewLoggerFromArgs("socks-forwarder"))
}

func mainWithArgs(ctx context.Context, args []string, logger logging.Logger) error {
	sfModule, err := module.NewModuleFromArgs(ctx, logger)
	if err != nil {
		return err
	}

	err = sfModule.AddModelFromRegistry(ctx, sensor.API, sensorModel)
	if err != nil {
		return err
	}

	err = sfModule.Start(ctx)
	if err != nil {
		return err
	}
	defer sfModule.Close(ctx)

	<-ctx.Done()
	return nil
}
