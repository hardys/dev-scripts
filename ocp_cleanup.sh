#!/bin/bash
set -x

source logging.sh
source common.sh
source ocp_install_env.sh

sudo systemctl stop fix_certs.timer
systemctl is-failed fix_certs.service >/dev/null && sudo systemctl reset-failed fix_certs.service

if [ -d ocp ]; then
    ocp/openshift-install --dir ocp --log-level=debug destroy bootstrap
    ocp/openshift-install --dir ocp --log-level=debug destroy cluster
    rm -rf ocp
fi

if [ -d ocp2 ]; then
    ocp2/openshift-install --dir ocp2 --log-level=debug destroy bootstrap
    ocp2/openshift-install --dir ocp2 --log-level=debug destroy cluster
    rm -rf ocp2
fi

sudo rm -rf /etc/NetworkManager/dnsmasq.d/openshift.conf

# Cleanup ssh keys for baremetal network
if [ -f $HOME/.ssh/known_hosts ]; then
    sed -i "/^192.168.111/d" $HOME/.ssh/known_hosts
    sed -i "/^api.${CLUSTER_DOMAIN}/d" $HOME/.ssh/known_hosts
fi

if test -f assets/templates/99_master-chronyd-redhat.yaml ; then
    rm -f assets/templates/99_master-chronyd-redhat.yaml
fi
if test -f assets/templates/99_worker-chronyd-redhat.yaml ; then
    rm -f assets/templates/99_worker-chronyd-redhat.yaml
fi

# If the installer fails before terraform completes the destroy bootstrap
# cleanup doesn't clean up the VM/volumes created..
for vm in $(sudo virsh list --all --name | grep "^${CLUSTER_NAME}.*bootstrap"); do
  sudo virsh destroy $vm
  sudo virsh undefine $vm --remove-all-storage
done

if [ -d assets/generated ]; then
  rm -rf assets/generated
fi
