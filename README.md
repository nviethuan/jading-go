# JADING-GO

## How to deploy

### New deploy

1. Download db file from s3: move to path: `~/app/data`
1. Build `ss3` from `./tools` and move to `/usr/local/bin`
1. Config whitelist IP from Binance,...
1. Build app and copy main file to `~/app`
1. Prepare `jading.service` file
1. Read and run `./scripts/setup.sh`

## Re-deploy

1. Config/update `$JADINGIP` in local machine: `~/.bashrc`
1. Run on local m1: `make build`
1. Run on ssh:
  - `cd ~/app`
  - `sudo chmod 755 main`
  - `sudo systemctl restart jading.service`
  - `journalctl -u jading -f`


## IP Change

In case the EC2 instance is unexpectedly restarted and a new public IP address is assigned, there are 3 things to do:

1. Update IP whitelist from binance
2. ...
3. ...
