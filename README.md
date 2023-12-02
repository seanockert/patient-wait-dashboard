# Patient Wait Dashboard

## Build webserver

`go build -ldflags="-s -w"`

This will compile the webserver into a single file: `patient-wait-dashboard`

## Run webserver

In the terminal, run:

`./patient-wait-dashboard`

This will start a local webserver running on port 8090

Now, find out the IP address of the host computer on the network. Usually something like `192.168.1.2`.

Next, open your web browser and go to `http://192.168.1.2:8090`. You should see the dashboard.

Go to `http://192.168.1.2:8090` on another device connected on the network.
