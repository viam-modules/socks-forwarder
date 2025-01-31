# Viam SOCKS forwarder module

The SOCKS forwarder module serves a generic service that can interface with and update the
`socks-forwarder` systemd service. `socks-forwarder` is one piece of the BLE-SOCKS bridge;
you can read more about the BLE-SOCKS bridge [here](https://github.com/viam-labs/ble-managed/tree/main).

The generic service served by this module can start, stop or restart the service through
`DoCommand`.

The module will also download and install the `socks-forwarder` systemd process. Updating
the module will download and install any new version of `socks-forwarder` (potentially
using the currently running `socks-forwarder`) and restart the systemd service with the
new version.

## Usage

Here are more details on how to use this module. Note that the `viam-server` process that
uses this module _must_ be run as root or with `sudo` in order to control the state of the
`socks-fowarder` systemd process.

### Fragment and Sample JSON

A Viam fragment is available
[here](https://app.viam.com/fragment/c799e8c9-3a8a-4df4-8c6d-1b9851fcd529/json) to
automatically include the module and its associated service in your config.

If you would like to use raw JSON, here is a snippet you can include in your Viam machine
config to install the module and construct its associated service. Note that the `version`
field should be updated as necessary.

```json
{
  "services": [
    {
      "name": "socks-forwarder-controller",
      "type": "generic",
      "model": "viam:ble-socks:controller",
    }
  ],
  "modules": [
    {
      "module_id": "viam:socks-forwarder",
      "name": "socks-forwarder",
      "type": "registry",
      "version": "~0.0.7"
    }
  ]
}
```


### Interacting with the generic service through DoCommand

Here are two examples of using `DoCommand` to start, stop, and restart the
`socks-forwarder` service. Note that stopping the service may make the Viam machine
unreachable through the internet, so you may not be able to start it again without direct
access to the machine through ethernet, wifi, or serial (mouse, keyboard, and monitor.)

#### Python

TODO

#### Golang

TODO

## Development

### dpkg updates

The .deb file (dpkg) must be manually updated by running `make dpkg` in [this
repo](https://github.com/viam-labs/ble-managed/tree/9aca1c2a0709056b442c408e34c8dc5f01d392b6/socks-forwarder)
and copying the created .deb file to this repo. `make dpkg` must be run on the appropriate
architecture (likely aarch64) to compile the .deb correctly.

### Updating the module

To update the module:
- Commit any changes directly to the `main` branch of this repository
- Run `go build .` to to build the `socks-forwarder-module` binary
- Ensure the `viam` CLI tool [installed](https://docs.viam.com/dev/tools/cli/#install)
- Ensure you are logged in with `viam login`
    - You must have write privileges for the `socks-forwarder` module
- Run `viam module upload --version "[new-version]" --platform "[platform]" meta.json`
  from the top of this directory
