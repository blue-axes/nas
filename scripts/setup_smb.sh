#!/bin/bash
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== NAS SMB (Samba) Setup ===${NC}"
echo ""

if [ "$(id -u)" -ne 0 ]; then
    echo -e "${RED}Please run as root: sudo bash setup_smb.sh${NC}"
    exit 1
fi

SMB_CONF="/etc/samba/smb.conf"
NAS_DATA_DIR="${NAS_DATA_DIR:-/aaa}"

if ! command -v smbd &> /dev/null; then
    echo -e "${YELLOW}Installing Samba...${NC}"
    apt-get update -qq && apt-get install -y samba > /dev/null 2>&1
    echo -e "${GREEN}Samba installed.${NC}"
fi

if [ ! -f "${SMB_CONF}.bak" ]; then
    cp "$SMB_CONF" "${SMB_CONF}.bak" 2>/dev/null || true
fi

cat > "$SMB_CONF" << EOF
[global]
   workgroup = WORKGROUP
   server string = NAS Server
   server role = standalone server
   security = user
   map to guest = Bad User
   guest account = nobody
   log file = /var/log/samba/log.%m
   max log size = 50
   dns proxy = no

[nas-data]
   comment = NAS Data Directory
   path = ${NAS_DATA_DIR}
   browseable = yes
   read only = no
   guest ok = yes
   create mask = 0666
   directory mask = 0777
   force user = root
EOF

echo -e "${GREEN}Samba config written to $SMB_CONF${NC}"

systemctl enable smbd 2>/dev/null || true
systemctl restart smbd 2>/dev/null || service smbd restart 2>/dev/null || true

echo ""
echo -e "${GREEN}=== SMB Setup Complete ===${NC}"
echo ""
echo -e "Share path:  ${YELLOW}\\\\$(hostname -I | awk '{print $1}')\\nas-data${NC}"
echo -e "Mount on Linux:   ${YELLOW}sudo mount -t cifs //$(hostname -I | awk '{print $1}')/nas-data /mnt/nas -o guest${NC}"
echo -e "Mount on macOS:   ${YELLOW}mount_smbfs //guest@$(hostname -I | awk '{print $1}')/nas-data /Volumes/nas${NC}"
echo -e "Access on Windows: ${YELLOW}\\\\$(hostname -I | awk '{print $1}')\\nas-data${NC}"