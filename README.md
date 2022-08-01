# UPI Recon CLI
![](images/logo.png)

**This tool doesn't work right now, as it depends on 3rd party APIs that aren't working the way they were.**

A command line tool for reconnaissance using virtual payment address (VPA).
This tool leverages the openness available with the UPI platform to find :
1. UPI ID and name associated with a mobile number
2. UPI ID and name associated with a gmail account
3. UPI ID and name associated with a vehicle registration number. Leveraging UPI id associated with a fastag.

This project is a golang port of [upi-recon](https://github.com/qurbat/upi-recon/) by [@squeal](https://twitter.com/squeal).

# Overview

```sh

                _    _ _____ _____   _____                         _____ _      _____ 
                | |  | |  __ \_   _| |  __ \                       / ____| |    |_   _|
                | |  | | |__) || |   | |__) |___  ___ ___  _ __   | |    | |      | |  
                | |  | |  ___/ | |   |  _  // _ \/ __/ _ \| '_ \  | |    | |      | |  
                | |__| | |    _| |_  | | \ \  __/ (_| (_) | | | | | |____| |____ _| |_ 
                \____/|_|   |_____| |_|  \_\___|\___\___/|_| |_|  \_____|______|_____|

                        #  Author: Aseem Shrey (@aseemshrey)
                        #  URL: https://github.com/LuD1161/upi-recon-cli
                        #  Website : https://aseemshrey.in
                        #  YouTube : https://www.youtube.com/c/HackingSimplifiedAS

Check virtual payment address corresponding to a mobile number, email address and get user's name as well.

Usage:
  upi-recon-cli PHONE_NUMBER [flags]
  upi-recon-cli [command]

Available Commands:
  checkAll    Check a particular number against all UPI identifiers.
  checkFastag Check Fast tag suffixes for vehicle registration number.
  checkGpay   Check gmail id corresponding to GPay suffixes.
  help        Help about any command

Flags:
  -h, --help            help for upi-recon-cli
  -t, --threads int     No of threads (default 100)
      --timeout int     Timeout for requests (default 15)
  -v, --version         version for upi-recon-cli

Use "upi-recon-cli [command] --help" for more information about a command.
```

### Checking a Mobile number for the Owner's name and UPI IDs
```sh
./upi-recon-cli <MOBILE_NUMBER_HERE>
```
![](images/usage-mobile-number.png)

### Checking a Vehicle Number for the Owner's name and UPI IDs
```sh
./upi-recon-cli checkFastag <VEHICLE_NUMBER>
```
![](images/usage-fastag.png)


### Checking a Gmail ID for the Owner's name and UPI IDs
```sh
./upi-recon-cli checkGpay <GMAIL_ID>
```
![](images/usage-google.png)


## Installation

1. Download the binaries for your platform from [releases page](https://github.com/LuD1161/upi-recon-cli/releases).
2. Extract the `tar.gz` file.
You'd find the following file strucuture inside the extracted folder : 
```sh
.
├── LICENSE
├── README.md
├── data
│   ├── all_suffixes.txt
│   ├── fastag_suffixes.txt
│   ├── gpay_suffixes.txt
│   └── mobile_suffixes.txt
└── upi-recon-cli

1 directory, 7 files
```
3. That's it. You're ready to go 🎉🚀

## Run with Gitpod

 Click this button to run your project on [Gitpod](https://gitpod.io/) which comes pre-configured with the go environment you need 🔥 

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/LuD1161/upi-recon-cli)


## Acknowledgements

- Karan S ([@squeal](https://twitter.com/squeal)): authored upi-recon
- Srikanth L ([@logic](https://twitter.com/logic)): contributed suffix files, introduced support for Google Pay & Fast tag addresses

## 🚀 About Me

This is [Aseem](https://aseemshrey.in). I'm a security engineer from India 🇮🇳.<br/>I am always curious about learning and building new things. Teaching security stuff through my youtube channel. Ping me up for anything related to security 🙌


<a href="https://twitter.com/intent/follow?screen_name=AseemShrey" target="_blank"><img src="https://img.shields.io/twitter/follow/AseemShrey?style=social&logo=twitter" alt="follow on Twitter"></a>
<a href="https://youtube.com/c/HackingSimplifiedAS?sub_confirmation=1" target="_blank"><img src="https://img.shields.io/youtube/channel/subscribers/UCARsgS1stRbRgh99E63Q3ng?label=HackingSimplified&style=social" alt="Subscribe on Youtube"></a>

## Disclaimer

Note: Unified Payment Interface ("UPI") Virtual Payment Addresses ("VPAs") do not carry a data security classification by virtue of their usage in practice, and should as such be considered to be public information, similar to how email addresses may be considered to be public information.

This tool allows users to 1) check the existence of UPI payment addresses, and 2) fetch associated information about the account holder, in an automated manner based on provided input. This functionality is already available (however, not in an automated fashion) through most UPI payment applications available on the Android and/or iOS platforms.

This tool is provided "AS IS" without any warranty of any kind, either expressed, implied, or statutory, to the extent permitted by applicable law.
