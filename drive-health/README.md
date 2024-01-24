## 📖 About
Drive Health is a program written in golang to help with tracking and monitoring of your hardware's temperature.

This tool had been conceived with the purpose of installing it on different servers I own with different configurations to help keep track of the temperature of hard-disks, ssds, nvme drives, etc...

### Features
- Disk Listing
- Temperature Graphing
- Disk activity logging
- [API](./lib/web/api.go)

![UI Example](./media/design_v1.webp)

## ❗ Disclaimer
I'm not exactly a linux hardware wizard, so I honestly have no clue about a lot of things and I myself can tell there's a lot to improve upon and that there's a lot of other things missing that are a little bit more obscure, I personally don't currently own any m.2 sata drives to test the code on, or many of the other drive types, I have only tested on HDD, SSD and NVMe drives, any issues opened would help me so much!

## ❗ Requirements
1. A linux machine, this will NOT work on macOS or on Windows, it's meant to be ran on servers as a service with which administrators can privately connect to for temperature logging.

2. Please make sure you have the [**drivetemp kernel drive**](https://docs.kernel.org/hwmon/drivetemp.html) you can check this by running `sudo modprobe drivetemp`.
The program depends on this to be able to log the temperature of your devices.


## 📖 How to use
1. Follow the `Deployment` section instrcutions to launch the program

2. Once the program has launched, access it in your browser

3. Enter the administrative username and password for the simple HTTP Auth

4. You now have access to the application, you can monitor your disk's temperature over a period of time.

## 🐦 Deployment
To deploy the application you have multiple choices, the preffered method should be one which runs the binary directly and not containerization, the `docker` image is taking up a wopping `1Gb+` because I have to include sqlite3-dev and musl-dev dependencies, which sucks, so I whole heartedly recommend just installing this on your system as a binary either with `SystemD` or whichever service manager you are using.

Download binaries from [the releases page](https://github.com/JustKato/drive-health/releases)

### 🐋 Docker
In the project there's a `docker-compose.prod.yml` which you can deploy on your server, you will notice that there's also a "dev" version, this version simply has a `build` instead of `image` property, so feel free to use either.

Please do take notice that I have just fed the `environment file` directly to the service via docker-compose, and I recommend you do the same but please feel free to pass in `environment` variables straight to the process as well.

[Docker Compose File](./docker-compose.prod.yml)
```yaml
version: "3.8"

services:
  drive-health:
    # Latest image pull, mention the specific version here please.
    image: ghcr.io/justkato/drive-health:latest
    # Restart in case of crashing
    restart: unless-stopped
    # Load environment variables from .env file
    env_file:
      - .env
    # Mount the volume to the local drive
    volumes:
      - ./data:/data
    # Setup application ports
    ports:
      - 5003:8080
```

### 💾 SystemD
When running with SystemD or any other service manager, please make sure you have a `.env` inside the `WorkingDirectory` of your runner, in the below example I will simply put my env in `/home/daniel/services/drive-health/.env`

```ini
[Unit]
Description=Drive Health Service
After=network.target

[Service]
Type=simple
User=daniel # Your user here
WorkingDirectory=/home/daniel/services/drive-health # The path to the service's directory
ExecStart=/home/daniel/services/drive-health/drive-health # The path to the binary
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

## ❔ FAQ

### How does it work?
Currently the program does not depend on any go library for hardware detection as I couldn't find anything that would not require root access while giving me the possibility to interrogate the temperature of the drives.

I chose not to depend on `lsblk` either, so how does the program work?
The program currently looks in `/sys/block` and then tries to make sense of the devices, I have had limited testing with my hardware specs, any issues being open in regards to different kinds of hardware would be highly appreciated

### Why not just run as root?
I really, REALLY, **REALLY** want to avoid asking people to run **ANY** program I write as root and even try and prevent that from happening since that's how things can go bad, especially because I am running actions over hardware devices.

## Support & Contribution
For support, bug reports, or feature requests, please open an issue on the [GitHub repository](https://github.com/JustKato/drive-health/issues). Contributions are welcome! Fork the repository, make your changes, and submit a pull request.

## License
This project is licensed under the [Apache License 2.0](./LICENSE).