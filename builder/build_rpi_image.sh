#!/bin/bash

############################
# cd BUILD_RPI
# SDCARD=sda WIFI_SSID=name WIFI_PASSWORD=pass SKIP_GET=1 SKIP_SDCARD=1 SKIP_INSTALL=1 ./build_rpi_image.sh
############################

WIFI_SSID=${WIFI_SSID:-weq}
WIFI_PASSWORD=${WIFI_PASSWORD:-pi=3.1415926+}
SDCARD=${SDCARD:-sda}
SDCARD=${SDCARD##*/}
DISK=/dev/$SDCARD
EXTRA_ROOT_SPACE_GB=${EXTRA_ROOT_SPACE_GB:-1} # extra rootFS space in GB

############################

if [ "`uname --machine`" != "armv7l" ] || ! grep -qis raspbian /etc/os-release ]; then
	echo "This script can only be ran from a Raspbian/armv7l instance"
	if [ -f /etc/os-release ] && grep -qis ^NAME /etc/os-release; then
		echo Your OS is `grep ^NAME= /etc/os-release | cut -d= -f 2`/`uname --machine`
	fi
	exit 1
fi

echo -e "\nCurrent PARAMS:\n\tSDCARD[$SDCARD] EXTRA_ROOT_SPACE_GB[$EXTRA_ROOT_SPACE_GB]\n\tWIFI_SSID[$WIFI_SSID] WIFI_PASSWORD[$WIFI_PASSWORD]\n\tSKIP_GET[$SKIP_GET] SKIP_SDCARD[$SKIP_SDCARD] SKIP_INSTALL[$SKIP_INSTALL]\n"
echo Press ENTER to continue, or CTRL-C to abort
read

############################


if [ "$SKIP_GET" = 1 ]; then
	echo skipping get phase...
	sleep 1
fi
if [ "$SKIP_SDCARD" = 1 ]; then
	echo skipping SDCARD building phase...
	sleep 1
fi
if [ "$SKIP_INSTALL" = 1 ]; then
	echo skipping installation phase...
	sleep 1
fi
############################
isMounted=''

set -e
#set -x

do_exit() {
	if [ "$isMounted" != "" ]; then
		umount ROOT/{dev/pts,dev,sys,proc,boot,tmp,DATA,}
		#umount $(mount|grep sda|cut -d' ' -f 1)
	fi
}
trap do_exit EXIT

do_err () {
	echo >&2 ERROR: $*
	echo >&2 ... bailing out
	echo >&2 ""
	
	exit 1
}

mk_sdcard() {
	DISK=/dev/$SDCARD

	[ ! -b $DISK ] && do_err DISK [$DISK] does not exist
	[ "`ls ${DISK}* | wc -l`" != '1' ] && do_err [$DISK] is not empty\; please use [fdisk] to delete all partitions, then try again

	dd if=raspberry.img of=${DISK} status=progress
	return 0
}

get_raspberry() {
	[ -f raspberry.img ] && return 0
	echo retrieving raspberry.img
	wget --continue https://downloads.raspberrypi.org/raspbian_lite_latest
	imgFile=$(unzip -l raspbian_lite_latest |grep '[.]img'|sed -e 's/^.*[ \t]//')
	[ -f "$imgFile" ] || unzip raspbian_lite_latest
	ln "$imgFile" raspberry.img
	rm -f raspbian_lite_latest
}

get_golang() {
	mkdir -p DATA
	if [ ! -f DATA/go.tgz ]; then
		wget --continue -P DATA/ https://dl.google.com/go/go1.13.8.linux-armv6l.tar.gz
		ln DATA/go1.13.8.linux-armv6l.tar.gz DATA/go.tgz
	fi
}
get_weqalign() {
	mkdir -p DATA
	if [ ! -d DATA/weqalign ]; then
		git -C DATA/ clone --depth=1 https://github.com/kairos10/weqalign
	fi
}

get_astrometry_src() {
	mkdir -p DATA/ASTROMETRY_DATA
	if [ ! -d DATA/astrometry.net ]; then
		git -C DATA/ clone --depth=1 https://github.com/dstndstn/astrometry.net
	fi

	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4219.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4218.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4217.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4216.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4215.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4214.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4213.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4212.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4211.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4210.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4209.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4208.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-11.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-10.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-09.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-08.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-07.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-06.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-05.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-04.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-03.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-02.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-01.fits
	wget --continue -P DATA/ASTROMETRY_DATA/ http://data.astrometry.net/4200/index-4207-00.fits

}

extend_fs() {
	numSectors=$(parted $DISK -ms unit s p | grep ^$DISK | cut -d: -f2 | sed -e 's/[^0-9]*//g')
	numParts=$(parted $DISK -ms unit s p | grep '^[0-9]'|wc -l)
	[ $numParts = 2 ] || do_err Unknown disk layout on [$DISK][$numParts partitions]

	p2Start=$(fdisk -l $DISK | grep ^${DISK}2 | sed -e 's/ \{1,\}/:/g' | cut -d: -f2)
	p2End=$(fdisk -l $DISK | grep ^${DISK}2 | sed -e 's/ \{1,\}/:/g' | cut -d: -f3)
	p2Sectors=$(fdisk -l $DISK | grep ^${DISK}2 | sed -e 's/ \{1,\}/:/g' | cut -d: -f4)
	p2Type=$(fdisk -l $DISK | grep ^${DISK}2 | sed -e 's/ \{1,\}/:/g' | cut -d: -f6)
	#p2NewEnd=$(( p2End + ( numSectors - p2End ) / 2 ))
	p2NewEnd=$(( p2End + $EXTRA_ROOT_SPACE_GB * ( 1024 * 1024 * 2 ) ))
	[ $pNewEnd -lt $pEnd ] || do_err "The card is too small [trying to extend the ROOT partition by 1G]"
	
	fdisk $DISK <<-EOF
		p
		d
		2
		n
		p
		2
		$p2Start
		$p2NewEnd
		t
		2
		$p2Type
		p
		w
		q
	EOF
	e2fsck -f ${DISK}2
	resize2fs ${DISK}2
}

mount_sdcard() {
	mkdir -p ROOT
	mount ${DISK}2 ROOT
	mount ${DISK}1 ROOT/boot
	mount --bind /dev ROOT/dev/
	mount --bind /sys ROOT/sys/
	mount --bind /proc ROOT/proc/
	mount --bind /dev/pts ROOT/dev/pts
	mount --bind /tmp ROOT/tmp
	mkdir -p ROOT/DATA
	mount --bind DATA ROOT/DATA

	isMounted=1
	mount |grep $DISK
}

prepare_sys_core() {
	chroot ROOT /bin/bash <<-EOF
		# base config
		echo 'zero' > /etc/hostname
		if ! grep -qs '^CONF_SWAPFILE=/dev/null' /etc/dphys-swapfile; then
			echo removing SWAP...
			echo CONF_SWAPFILE=/dev/null >> /etc/dphys-swapfile
			rm -f /var/swap
		fi
		if ! grep -qs '^gpu_mem=' /boot/config.txt; then
			echo gpu_mem=128 >> /boot/config.txt
		fi
		if ! grep -qs '^tmpfs' /etc/fstab; then
			cat >> /etc/fstab <<-XX
				tmpfs /tmp                tmpfs    defaults,noatime,nosuid,size=200m                  0 0
				tmpfs /var/tmp            tmpfs    defaults,noatime,nosuid,size=30m                   0 0
				tmpfs /var/log            tmpfs    defaults,noatime,nosuid,mode=0755,size=100m        0 0
				tmpfs /run                tmpfs    defaults,noatime,nosuid,mode=0755,size=2m          0 0
				tmpfs /var/spool/mqueue   tmpfs    defaults,noatime,nosuid,mode=0700,gid=12,size=30m  0 0
			XX
		fi

		# freeup the UART port
		if ! grep -qs '^dtoverlay=disable-bt' /boot/config.txt; then
			echo dtoverlay=disable-bt >> /boot/config.txt
		fi
		systemctl disable hciuart

		systemctl unmask ssh
		systemctl enable ssh
		apt -y install tmux git

		#update
		apt-get update
		apt-get -y dist-upgrade
		apt-get -y autoremove

	EOF
}
prepare_sys_wifiAP() {
	chroot ROOT /bin/bash <<-EOF
		# wifi AP
		rm -f /etc/systemd/network/99-default.link
		ln -s /dev/null /etc/systemd/network/99-default.link
		apt install hostapd dnsmasq

		if ! grep -qs '^denyinterfaces wlan0' /etc/dhcpcd.conf; then
			echo "denyinterfaces wlan0" >> /etc/dhcpcd.conf
		fi
		if ! grep -qs '^DAEMON_CONF=' /etc/default/hostapd; then
			echo 'DAEMON_CONF="/etc/hostapd/hostapd.conf"' >> /etc/default/hostapd
		fi

		cat > /etc/network/interfaces.d/wlan0 <<-XX
		allow-hotplug wlan0
		iface wlan0 inet static
    			address 192.168.5.1
    			netmask 255.255.255.0
    			network 192.168.5.0
    			broadcast 192.168.5.255
		XX

		cat > /etc/hostapd/hostapd.conf <<-XX
			interface=wlan0
			ssid=$WIFI_SSID
			wpa_passphrase=$WIFI_PASSWORD
			driver=nl80211
			country_code=RO
			hw_mode=g
			ieee80211n=1
			channel=1
			wpa=2
			wpa_key_mgmt=WPA-PSK
			rsn_pairwise=CCMP
			auth_algs=1
			macaddr_acl=0
			wmm_enabled=1
			ignore_broadcast_ssid=0
		XX

		cat > /etc/dnsmasq.conf <<-XX
			interface=wlan0
			listen-address=192.168.5.1
			bind-dynamic
			server=192.168.5.1
			domain-needed
			bogus-priv
			dhcp-range=192.168.5.100,192.168.5.200,24h
		XX

		systemctl unmask hostapd
		systemctl enable hostapd

	EOF
}

prepare_sys_install_golang() {
	chroot ROOT/ /bin/bash <<-EOF
		if [ ! -d /usr/local/go ]; then
			mkdir /usr/local
			cd /usr/local
			tar zxf /DATA/go.tgz

			cat > /etc/profile.d/golang.sh <<-XX
				if [ -d /usr/local/go ]; then
        				PATH="\\\$PATH:/usr/local/go/bin"
				fi
			XX
		fi
	EOF
}
prepare_sys_install_weqalign() {
	chroot ROOT/ /bin/bash <<-EOF
		su - pi <<-XX
			if [ ! -d weqalign ]; then
				git clone /DATA/weqalign
			fi
			cd weqalign
			go get golang.org/x/net/webdav
			go get -u golang.org/x/sys/...
			go get github.com/fsnotify/fsnotify
			git fetch
			git pull
			make
		XX
	EOF
}

prepare_sys_install_astrometry() {
	chroot ROOT/ /bin/bash <<-EOF
		if [ ! -d /usr/local/astrometry ]; then
			apt -y install libcfitsio-dev
			apt -y install libcairo2-dev
			apt -y install libnetpbm10-dev netpbm
			apt -y install zlib1g-dev libbz2-dev
			apt -y install libjpeg-dev
			apt -y install swig
			#apt -y install python-dev python-pyfits python-numpy
			apt -y install libpng-dev

			apt -y install python3-pyfits python3-numpy python3-dev

			cd /DATA/astrometry.net
			make install
			#make clean

			cp /DATA/ASTROMETRY_DATA/*.fits /usr/local/astrometry/data/

			echo Building astrometry.cfg
			cat > /usr/local/astrometry/etc/astrometry.cfg <<-XX
				inparallel
				minwidth 0.8
				maxwidth 10
				depths 10,35
				cpulimit 300
				add_path /usr/local/astrometry/data
				#autoindex
				#index index-4212
				#index index-4213
				#index index-4214
				#index index-4215
			XX
			for fFits in /usr/local/astrometry/data/*.fits; do
				fBase=\${fFits##*/}
				fIndex=\${fBase%%.*}
				echo "index \${fIndex}" >> /usr/local/astrometry/etc/astrometry.cfg
			done

			cat > /etc/profile.d/astrometry-det.sh <<-XX
			if [ -d /usr/local/astrometry/bin ]; then
       				PATH="\\\$PATH:/usr/local/astrometry/bin"
			fi
			XX
		fi
	EOF
}
prepare_sys_summary() {
	chroot ROOT/ /bin/bash <<-EOF
		df -h /
	EOF
}


##########################################
##########################################
##########################################

pwd0=`pwd`
cd -P "`dirname "${BASH_SOURCE[0]}"`"
if [ "$pwd0" != "`pwd`" ]; then echo Working folder changed to `pwd`; fi

mkdir -p BUILD
cd BUILD/

if [ "$SKIP_GET" != 1 ]; then
	get_raspberry
	get_golang
	get_weqalign
	get_astrometry_src
fi
if [ "$SKIP_SDCARD" != 1 ]; then
	mk_sdcard $1
	extend_fs
	#DISK=/dev/sda
fi
if [ "$SKIP_INSTALL" != 1 ]; then
	mount_sdcard
	prepare_sys_core
	prepare_sys_wifiAP
	prepare_sys_install_golang
	prepare_sys_install_astrometry
	prepare_sys_install_weqalign

	prepare_sys_summary
fi

