# How to run on Gitpod

## Install vsftpd

```sh
$ sudo apt-get update
$ sudo apt-get install -y vsftpd
```

## Configure vsftpd

Edit ```/etc/vsftpd.conf```

```ini
anonymous_enable=YES
local_enable=YES
write_enable=YES
local_umask=022
chroot_local_user=NO
```

## Start vsftpd

```sh
$ sudo mkdir -p /var/run/vsftpd/empty
$ sudo vsftpd /etc/vsftpd.conf
```

## Run example

```sh
$ task applayer-ftp
```

## Connect server

```sh
$ sudo apt-get install -y lftp
$ lftp -u anonymous, anonymous@localhost
```
