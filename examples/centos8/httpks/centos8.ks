#version=RHEL8
ignoredisk --only-use=sda
#autopart --type=lvm

zerombr
clearpart --all --initlabel
clearpart --all --drives=sda
ignoredisk --only-use=sda
part /boot --fstype="xfs" --ondisk=sda --size=512
part / --fstype="xfs" --ondisk=sda --grow --size=1

# Partition clearing information
clearpart --none --initlabel
# Use graphical install
text
repo --name="AppStream" --baseurl=file:///run/install/repo/AppStream
# Use CDROM installation media
cdrom
# Keyboard layouts
keyboard --vckeymap=us --xlayouts='us'
# System language
lang en_US.UTF-8

firewall --disabled
selinux --enforcing

# Network information
network  --bootproto=dhcp --device=enp0s3 --ipv6=auto --activate
network  --hostname=centos8
# Root password
rootpw SET_THIS_VALUE
# Run the Setup Agent on first boot
firstboot --disabled

reboot
# Do not configure the X Window System
skipx
# System services
services --disabled="chronyd"
# System timezone
timezone Pacific/Auckland --isUtc --nontp

%packages
@^server-product-environment
#@container-management
@performance
#@remote-system-management
#@rpm-development-tools
@security-tools
@system-tools

%end

%addon com_redhat_kdump --disable --reserve-mb='auto'

%end

%anaconda
pwpolicy root --minlen=6 --minquality=1 --notstrict --nochanges --notempty
pwpolicy user --minlen=6 --minquality=1 --notstrict --nochanges --emptyok
pwpolicy luks --minlen=6 --minquality=1 --notstrict --nochanges --notempty
%end

%post
# Install open-vm-tools, required to detect IP when building on ESXi
yum -y install open-vm-tools
systemctl enable vmtoolsd
%end
