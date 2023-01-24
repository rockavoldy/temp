# CPU Temp InfluxDB Publisher

This is a simple script that reads the CPU temperature from a debian/ubuntu, and publishes it to an InfluxDB v2 database.

## How to run
Set your credentials to the InfluxDB first by export the env variables
```bash
export INFLUX_URL=http://localhost:8086
export INFLUX_TOKEN=<your-influxdb-token>
export INFLUX_BUCKET=<bucket-name>
export INFLUX_ORG=<org-name>
```
Then run the app
```bash
./temp
```
if you want to run it in the background, you can create a systemd service file inside `/etc/systemd/system/` and run it as a service
```bash
[Unit]
Description=InfluxDB Publisher
After=network.target

[Service]
Type=simple
User=root
Restart=always
RestartSec=5s
EnvironmentFile=/etc/default/temp/.env
ExecStart=/usr/bin/temp
PermissionsStartOnly=true
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=sleepservice

[Install]
WantedBy=multi-user.target
```