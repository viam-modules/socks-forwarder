# Viam SOCKS forwarder module

The SOCKS forwarder module serves a sensor that can interface with the
socks-forwarder systemd service. The sensor reports the tx/rx of the
bluetooth adapter through `Readings`. It can also start, stop or restart
the service through `DoCommand`.

The .deb file must be manually updated by running `make dpkg` in [this
repo](https://github.com/viam-labs/ble-managed/tree/9aca1c2a0709056b442c408e34c8dc5f01d392b6/socks-forwarder)
and copying the created .deb file here. `make dpkg` must be run on the
appropriate architecture to compile the .deb correctly.
