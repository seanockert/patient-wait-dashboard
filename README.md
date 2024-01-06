# Patient Wait Dashboard

A patient wait dashboard for displaying in a hospital Emergency Department. Show the number of patients waiting and estimate the wait time.

Hosted locally on a Go webserver. Uses a websocket to update all dashboards on the local network when the patient count changes.

## Demo (dashboard only)

https://so-ed-dashboard.surge.sh

<img width="1134" alt="Screenshot 2023-12-03 at 11 12 40 pm" src="https://github.com/seanockert/patient-wait-dashboard/assets/574163/28232849-f995-4a42-92af-9a133871f366">


## Run the web server with Go

If you have the Go language installed, start the webserver directly with:

```
go run main.go
```

Now, find out the IP address of the host computer on the network. Usually something like `192.168.1.2`.

Next, open your web browser and go to `http://192.168.1.2:8090`. You should see the dashboard.

Go to `http://192.168.1.2:8090` on another device connected on the network.

## Build for different environments

If you don't have Go installed on the target device, you can compile a single file to run.

### Build for Windows

```bash
GOOS=windows GOARCH=amd64 go build -o patient_wait_dashboard.exe -ldflags="-s -w"
```

### Build for MacOS (Intel)

```bash
GOOS=darwin GOARCH=amd64 go build -o patient_wait_dashboard_mac_intel -ldflags="-s -w"
```

### Build for MacOS (ARM)

```bash
GOOS=darwin GOARCH=arm64 go build -o patient_wait_dashboard_mac_m1 -ldflags="-s -w"
```

This will compile the webserver into a single file: `patient_wait_dashboard`. You can now run this from anywhere with:

```bash
`./patient_wait_dashboard`
```

Make sure to include the `index.html` and `/assets` folder too.

Note: the `-ldflags="-s -w"` part is just for making a smaller file. The file should be about 6MB in size.

## Admin

Currently you can tap the "Update" button below the patient count to adjust the patients.

There's also a basic Admin page at `http://192.168.1.2:8090/admin` to do the same thing. So we can remove the "Update" button.

## TODO

- Improve the design
