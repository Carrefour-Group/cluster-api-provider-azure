/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package azure

import (
	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha4"
)

// PublicIPSpec defines the specification for a Public IP.
type PublicIPSpec struct {
	Name    string
	DNSName string
	IsIPv6  bool
}

// NICSpec defines the specification for a Network Interface.
type NICSpec struct {
	Name                      string
	MachineName               string
	SubnetName                string
	VNetName                  string
	VNetResourceGroup         string
	StaticIPAddress           string
	PublicLBName              string
	PublicLBAddressPoolName   string
	PublicLBNATRuleName       string
	InternalLBName            string
	InternalLBAddressPoolName string
	PublicIPName              string
	VMSize                    string
	AcceleratedNetworking     *bool
	IPv6Enabled               bool
	EnableIPForwarding        bool
}

// DiskSpec defines the specification for a Disk.
type DiskSpec struct {
	Name string
}

// LBSpec defines the specification for a Load Balancer.
type LBSpec struct {
	Name              string
	Role              string
	Type              infrav1.LBType
	SKU               infrav1.SKU
	SubnetName        string
	BackendPoolName   string
	FrontendIPConfigs []infrav1.FrontendIP
	APIServerPort     int32
}

// RouteTableRole defines the unique role of a route table.
type RouteTableRole string

// RouteTableSpec defines the specification for a Route Table.
type RouteTableSpec struct {
	Name   string
	Subnet infrav1.SubnetSpec
}

// InboundNatSpec defines the specification for an inbound NAT rule.
type InboundNatSpec struct {
	Name             string
	LoadBalancerName string
}

// SubnetSpec defines the specification for a Subnet.
type SubnetSpec struct {
	Name              string
	CIDRs             []string
	VNetName          string
	RouteTableName    string
	SecurityGroupName string
	Role              infrav1.SubnetRole
}

// VNetSpec defines the specification for a Virtual Network.
type VNetSpec struct {
	ResourceGroup string
	Name          string
	CIDRs         []string
}

// RoleAssignmentSpec defines the specification for a Role Assignment.
type RoleAssignmentSpec struct {
	MachineName  string
	Name         string
	ResourceType string
}

// ResourceType defines the type azure resource being reconciled.
// Eg. Virtual Machine, Virtual Machine Scale Sets
type ResourceType string

const (

	// VirtualMachine ...
	VirtualMachine = "VirtualMachine"

	// VirtualMachineScaleSet ...
	VirtualMachineScaleSet = "VirtualMachineScaleSet"
)

// NSGSpec defines the specification for a Security Group.
type NSGSpec struct {
	Name         string
	IngressRules infrav1.IngressRules
}

// VMSpec defines the specification for a Virtual Machine.
type VMSpec struct {
	Name                   string
	Role                   string
	NICNames               []string
	SSHKeyData             string
	Size                   string
	Zone                   string
	Identity               infrav1.VMIdentity
	OSDisk                 infrav1.OSDisk
	DataDisks              []infrav1.DataDisk
	UserAssignedIdentities []infrav1.UserAssignedIdentity
	SpotVMOptions          *infrav1.SpotVMOptions
	SecurityProfile        *infrav1.SecurityProfile
}

// BastionSpec defines the specification for bastion host.
type BastionSpec struct {
	Name         string
	SubnetName   string
	PublicIPName string
	VNetName     string
}

// ScaleSetSpec defines the specification for a Scale Set.
type ScaleSetSpec struct {
	Name                         string
	Size                         string
	Capacity                     int64
	SSHKeyData                   string
	OSDisk                       infrav1.OSDisk
	DataDisks                    []infrav1.DataDisk
	SubnetName                   string
	VNetName                     string
	VNetResourceGroup            string
	PublicLBName                 string
	PublicLBAddressPoolName      string
	AcceleratedNetworking        *bool
	TerminateNotificationTimeout *int
	Identity                     infrav1.VMIdentity
	UserAssignedIdentities       []infrav1.UserAssignedIdentity
	SecurityProfile              *infrav1.SecurityProfile
	SpotVMOptions                *infrav1.SpotVMOptions
	FailureDomains               []string
}

// TagsSpec defines the specification for a set of tags.
type TagsSpec struct {
	Scope      string
	Tags       infrav1.Tags
	Annotation string
}

// PrivateDNSSpec defines the specification for a private DNS zone.
type PrivateDNSSpec struct {
	ZoneName          string
	VNetName          string
	VNetResourceGroup string
	LinkName          string
	Records           []infrav1.AddressRecord
}

// AvailabilitySetSpec defines the specification for an availability set.
type AvailabilitySetSpec struct {
	Name string
}

// VMExtensionSpec defines the specification for a VM extension.
type VMExtensionSpec struct {
	Name              string
	VMName            string
	Publisher         string
	Version           string
	ProtectedSettings map[string]string
}

// VMSSExtensionSpec defines the specification for a VMSS extension.
type VMSSExtensionSpec struct {
	Name              string
	ScaleSetName      string
	Publisher         string
	Version           string
	ProtectedSettings map[string]string
}
