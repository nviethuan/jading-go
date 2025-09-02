#!/bin/bash

sudo cp jading.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable jading.service
sudo systemctl start jading.service
sudo systemctl status jading.service

