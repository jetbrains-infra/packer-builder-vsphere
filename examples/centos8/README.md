1. This assumes you have uplaoded CentOS-8-x86_64-1905-dvd1.iso into an ISOs folder on vsphere.  Ensure the MD5sum matches, or update the MD5sum in packer.
2. EDIT CentoOS8_build.json and httpks/centos8.ks, replacing all strings of: "SET_THIS_VALUE" to the correct value for your vsphere.
3. COPY creds.json.example to creds.json and update it to have your vsphere credentials
4. Ensure your vsphere can connect on port 8053 to the server you are running packer on, and no local firewall is blocking that port.
5. RUN: packer build -var-file creds.json CentOS8_build.json
6. ADD: provisions as required.