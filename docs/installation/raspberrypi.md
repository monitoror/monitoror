# RaspberryPi
## Installation Prerequisites
```shell script
# Chromium
sudo apt-get update
sudo apt-get install chromium-browser --yes

# Disabling screen saver
# Copy the entire block until END tag (included)
sudo tee -a /usr/bin/screensaver_off.sh > /dev/null <<END
#!/bin/bash
sleep 10 &&
sudo xset s 0 0
sudo xset s off
exit 0
END

sudo chmod +x /usr/bin/screensaver_off.sh

# :warning: you need to do this with user starting x server (usually pi)
# Copy the entire block until END tag (included)
tee -a ~/.config/autostart/chromium-browser.desktop > /dev/null <<END
[Desktop Entry]
Type=Application
Exec=/usr/bin/screensaver_off.sh
Hidden=false
X-MATE-Autostart-enabled=true
Name[fr_FR]=screensaver_off
Name=screensaver_off
Comment[fr_FR]=
Comment=
END
```

## Installation of Monitoror
```shell script
# Installing binariy file
sudo mkdir /opt/monitoror
sudo wget https://github.com/monitoror/monitoror/releases/download/{VERSION}/monitoror-linux-arm -P /opt/monitoror
sudo chmod +x /opt/monitoror/monitoror-linux-arm

# Installing backend configuration file
sudo nano /opt/monitoror/.env 

# If you want to run monitoror with another user, you can change owner of binary
```

## Starting Monitoror
### Manually
```shell script
# Starting backend
/opt/monitoror/monitoror-linux-arm

# Starting frontend
chromium-browser --kiosk --password-store=basic --disable-infobars \
  --app=http://localhost:8080/?configUrl={CONFIG_URL}
```

### Automatically
This part explain how to start Monitoror automatically on RaspberryPi startup

```shell script
# Backend startup with systemd
# Copy the entire block until END tag (included)
sudo tee -a /lib/systemd/system/monitoror.service > /dev/null <<END
[Unit]
Description=Monitoror
After=multi-user.target

[Service]
Type=idle
ExecStart=/opt/monitoror/monitoror-linux-arm
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=monitoror

[Install]
WantedBy=multi-user.target
END

systemctl enable monitoror.service
systemctl start monitoror.service
```

```shell script
# rsyslog configuration
# Copy the entire block until END tag (included)
sudo tee -a /etc/rsyslog.d/99-monitoror.conf > /dev/null <<END
if \$programname == 'monitoror' then /var/log/monitoror/monitoror.log
& stop
END

sudo mkdir /var/log/monitoror/
sudo touch /var/log/monitoror/monitoror.log
sudo chown syslog:adm /var/log/monitoror/monitoror.log

sudo service rsyslog restart
```

```shell script
# Frontend startup with autostart
# :warning: you need to do this with user starting x server (usually pi)
# Copy the entire block until END tag (included)
tee -a ~/.config/autostart/chromium-browser.desktop > /dev/null <<END
[Desktop Entry]
Type=Application
Exec=chromium-browser --kiosk --password-store=basic --disable-infobars --app=http://localhost:8080/?configUrl={CONFIG_URL}
Hidden=false
X-MATE-Autostart-enabled=true
Name[fr_FR]=chromium-browse
Name=chromium-browser
Comment[fr_FR]=
Comment=
END
```

