# JADING-GO

## How to deploy

1. Run on local m1: `make build`
2. Run on ssh:
  - `cd ~/app`
  - `sudo chmod 755 main`
  - `sudo systemctl restart jading.service`
  - `journalctl -u jading -f`
