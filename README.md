# gfs

A lightweight, high-performance file transfer tool designed for local networks. gfs uses UDP broadcast for automatic host discovery and TCP for reliable file transfer, allowing you to quickly share files between devices on the same network without any configuration. Simply run the command, select the target host, and transfer files with real-time progress feedback.
> [!WARNING]
>If you want to transfer files over the internet rather than within the same local network, you will need a VPN such as [Radmin VPN](https://www.radmin-vpn.com). gfs works only within a LAN by default.

### Usage

```bash
gfs send <path>  # Start hosting a file
gfs get [dest]   # Scan and download a file
```
