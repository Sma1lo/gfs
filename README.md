# gfs

[![Go Version](https://img.shields.io/github/go-mod/go-version/Sma1lo/gfs)](https://github.com/Sma1lo/gfs)
[![License](https://img.shields.io/github/license/Sma1lo/gfs)](https://github.com/Sma1lo/gfs/blob/main/LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/Sma1lo/gfs?include_prereleases)](https://github.com/Sma1lo/gfs/releases)
[![Release Date](https://img.shields.io/github/release-date/Sma1lo/gfs)](https://github.com/Sma1lo/gfs/releases)
[![Downloads](https://img.shields.io/github/downloads/Sma1lo/gfs/total)](https://github.com/Sma1lo/gfs/releases)
[![Last Commit](https://img.shields.io/github/last-commit/Sma1lo/gfs)](https://github.com/Sma1lo/gfs/commits)
[![Commit Activity](https://img.shields.io/github/commit-activity/m/Sma1lo/gfs)](https://github.com/Sma1lo/gfs/graphs/commit-activity)
[![Repo Size](https://img.shields.io/github/repo-size/Sma1lo/gfs)](https://github.com/Sma1lo/gfs)
[![Issues](https://img.shields.io/github/issues/Sma1lo/gfs)](https://github.com/Sma1lo/gfs/issues)
[![Pull Requests](https://img.shields.io/github/issues-pr/Sma1lo/gfs)](https://github.com/Sma1lo/gfs/pulls)

A lightweight, high-performance file transfer tool designed for local networks. gfs uses UDP broadcast for automatic host discovery and TCP for reliable file transfer, allowing you to quickly share files between devices on the same network without any configuration. Simply run the command, select the target host, and transfer files with real-time progress feedback.
> [!NOTE]
>If you want to transfer files over the internet rather than within a local network, you will need a [Radmin VPN](https://www.radmin-vpn.com). gfs works only within a LAN by default.

### Usage

```bash
gfs send <path>  # Start hosting a file
gfs get [dest]   # Scan and download a file
```
