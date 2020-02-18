# build_rpi_image.sh
This script creates a new Raspbian image, on a secondary SDCARD/USB drive, with everithing that is needed to run *weqalign*.
The script needs to be started from a working Raspbian installation.

* identify the secondary drive
```
dmesg | grep -i 'Attached .*removable disk'
```

* If the secondary drive is detected (for example) as *sda*, make sure that the card is empty (if the card has any partitions, the script will refuse to run)
```
sudo fdisk -l /dev/sda
```

* delete any existing partitions from the card
```
sudo fdisk /dev/sda
```
> **sudo fdisk /dev/sda**
> 
> Welcome to fdisk (util-linux 2.29.2).
> Changes will remain in memory only, until you decide to write them.
> Be careful before using the write command.
> 
> 
> Command (m for help): **p**
> Disk /dev/sda: 58 GiB, 62277025792 bytes, 121634816 sectors
> Units: sectors of 1 * 512 = 512 bytes
> Sector size (logical/physical): 512 bytes / 512 bytes
> I/O size (minimum/optimal): 512 bytes / 512 bytes
> Disklabel type: dos
> Disk identifier: 0x00000000
> 
> Device     Boot Start       End   Sectors Size Id Type
> /dev/sda1       32768 121634815 121602048  58G  7 HPFS/NTFS/exFAT
> 
> Command (m for help): **d**
> Selected partition 1
> Partition 1 has been deleted.
> 
> Command (m for help): **p**
> Disk /dev/sda: 58 GiB, 62277025792 bytes, 121634816 sectors
> Units: sectors of 1 * 512 = 512 bytes
> Sector size (logical/physical): 512 bytes / 512 bytes
> I/O size (minimum/optimal): 512 bytes / 512 bytes
> Disklabel type: dos
> Disk identifier: 0x00000000
> 
> Command (m for help): **w**
> The partition table has been altered.
> Calling ioctl() to re-read partition table.
> Syncing disks.

* prefetch downloadable items
```
SKIP_SDCARD=1 SKIP_INSTALL=1 ./build_rpi_image.sh
```

* install the image the the new card (the script is working in a chrooted environment).
By default, the script creates a new wifi access point with SSID[weq] and password[pi=3.1415926+]; you can change those defaults through the env variables WIFI_SSID and WIFI_PASSWORD
```
sudo SDCARD=sda WIFI_SSID=someName WIFI_PASSWORD=secretPassword SKIP_GET=1 ./build_rpi_image.sh
```

* boot the new image, connect to the new AP, open a ssh session to 192.168.5.1 (pi/raspberry), enjoy
