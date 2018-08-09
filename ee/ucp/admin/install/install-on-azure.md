---
title: Install UCP on Azure
description: Learn how to install Docker Universal Control Plane in a Microsoft Azure environment.
keywords: Universal Control Plane, UCP, install, Docker EE, Azure, Kubernetes
---

Docker UCP configures the Azure IPAM module for Kubernetes to allocate
IP addresses to Kubernetes pods. The Azure IPAM module requires each Azure
VM that's part of the Kubernetes cluster to be configured with a pool of
IP addresses.

You have two options for deploying the VMs for the Kubernetes cluster on Azure:
- Install the cluster on Azure stand-alone virtual machines. Docker UCP provides
  an [automated mechanism](#configure-ip-pools-for-azure-stand-alone-vms)
  to configure and maintain IP pools for stand-alone Azure VMs.
- Install the cluster on an Azure virtual machine scale set. Configure the
  IP pools by using an ARM template like [this one](#set-up-ip-configurations-on-an-azure-virtual-machine-scale-set).

The steps for setting up IP address management are different in the two
environments. If you're using a scale set, you set up `ipConfigurations`
in an ARM template. If you're using stand-alone VMs, you set up IP pools
for each VM by using a utility container that's configured to run as a
global Swarm service, which Docker provides.

## Considerations for size of IP pools

The subnet and the virtual network associated with the primary interface of
the Azure VMs need to be configured with a large enough address prefix/range. 
The number of required IP addresses depends on the number of pods running
on each node and the number of nodes in the cluster.

For example, in a cluster of 256 nodes, to run a maximum of 128 pods
concurrently on a node, make sure that the address space of the subnet and the
virtual network can allocate at least 128 * 256 IP addresses, _in addition to_
initial IP allocations to VM NICs during Azure resource creation.

Accounting for IP addresses that are allocated to NICs during VM bring-up, set
the address space of the subnet and virtual network to 10.0.0.0/16. This
ensures that the network can dynamically allocate at least 32768 addresses,
plus a buffer for initial allocations for primary IP addresses.

> Azure IPAM, UCP, and Kubernetes
> 
> The Azure IPAM module queries an Azure virtual machine's metadata to obtain
> a list of IP addresses that are assigned to the virtual machine's NICs. The
> IPAM module allocates these IP addresses to Kubernetes pods. You configure the
> IP addresses as `ipConfigurations` in the NICs associated with a virtual machine
> or scale set member, so that Azure IPAM can provide them to Kubernetes when
> requested.
{: .important}

## Configure IP pools for Azure stand-alone VMs

Follow these steps when the cluster is deployed using stand-alone Azure VMs.

### Create an Azure resource group

Create an Azure resource group with VMs representing the nodes of the cluster
by using the Azure Portal, CLI, or ARM template.

[Install UCP on the cluster](#install-ucp-on-the-cluster).

## Set up IP configurations on an Azure virtual machine scale set

Configure IP Pools for each member of the VM scale set during provisioning by
associating multiple `ipConfigurations` with the scale set’s
`networkInterfaceConfigurations`. Here's an example `networkProfile`
configuration for an ARM template that configures pools of 32 IP addresses
for each VM in the VM scale set.

```json
"networkProfile": {
  "networkInterfaceConfigurations": [
    {
      "name": "[variables('nicName')]",
      "properties": {
        "ipConfigurations": [
          {
            "name": "[variables('ipConfigName1')]",
            "properties": {
              "primary": "true",
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              },
              "loadBalancerBackendAddressPools": [
                {
                  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('loadBalancerName'), '/backendAddressPools/', variables('bePoolName'))]"
                }
              ],
              "loadBalancerInboundNatPools": [
                {
                  "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/loadBalancers/', variables('loadBalancerName'), '/inboundNatPools/', variables('natPoolName'))]"
                }
              ]
            }
          },
          {
            "name": "[variables('ipConfigName2')]",
            "properties": {
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              }
            }
          }
          .
          .
          .
          {
            "name": "[variables('ipConfigName32')]",
            "properties": {
              "subnet": {
                "id": "[concat('/subscriptions/', subscription().subscriptionId,'/resourceGroups/', resourceGroup().name, '/providers/Microsoft.Network/virtualNetworks/', variables('virtualNetworkName'), '/subnets/', variables('subnetName'))]"
              }
            }
          }
        ],
        "primary": "true"
      }
    }
  ]
}
```

## Install UCP on the cluster

Use the following command to install UCP on the manager node.
The `--pod-cidr` option maps to the IP address range that you configured for
the subnets in the previous sections, and the `--host-address` maps to the
IP address of the master node.

```bash
docker container run --rm -it \
  --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  {{ page.ucp_org }}/{{ page.ucp_repo }}:{{ page.ucp_version }} install \
  --host-address <ucp-ip> \
  --interactive \
  --swarm-port 3376 \
  --pod-cidr <ip-address-range> \
  --cloud-provider Azure
```
