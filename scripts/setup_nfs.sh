#!/bin/bash
set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}=== NAS NFS Setup ===${NC}"
echo ""

if [ "$(id -u)" -ne 0 ]; then
    echo -e "${RED}Please run as root: sudo bash setup_nfs.sh${NC}"
    exit 1
fi

EXPORTS_FILE="/etc/exports"
NAS_DATA_DIR="${NAS_DATA_DIR:-/aaa}"

if ! command -v exportfs &> /dev/null; then
    echo -e "${YELLOW}Installing NFS server...${NC}"
    apt-get update -qq && apt-get install -y nfs-kernel-server > /dev/null 2>&1
    echo -e "${GREEN}NFS server installed.${NC}"
fi

if [ ! -f "${EXPORTS_FILE}.bak" ]; then
    cp "$EXPORTS_FILE" "${EXPORTS_FILE}.bak" 2>/dev/null || true
fi

grep -q "$NAS_DATA_DIR" "$EXPORTS_FILE" 2>/dev/null && \
    sed -i "\|${NAS_DATA_DIR}|d" "$EXPORTS_FILE"

echo "${NAS_DATA_DIR} *(rw,sync,no_subtree_check,no_root_squash,insecure)" >> "$EXPORTS_FILE"

echo -e "${GREEN}NFS export added to $EXPORTS_FILE${NC}"

exportfs -ra
systemctl enable nfs-kernel-server 2>/dev/null || true
systemctl restart nfs-kernel-server 2>/dev/null || service nfs-kernel-server restart 2>/dev/null || true

echo ""
echo -e "${GREEN}=== NFS Setup Complete ===${NC}"
echo ""
echo -e "Mount on Linux:   ${YELLOW}sudo mount -t nfs $(hostname -I | awk '{print $1}'):${NAS_DATA_DIR} /mnt/nas${NC}"
echo -e "Mount on macOS:   ${YELLOW}sudo mount -t nfs -o resvport $(hostname -I | awk '{print $1}'):${NAS_DATA_DIR} /Volumes/nas${NC}"
echo -e "Add to /etc/fstab: ${YELLOW}$(hostname -I | awk '{print $1}'):${NAS_DATA_DIR} /mnt/nas nfs defaults 0 0${NC}"