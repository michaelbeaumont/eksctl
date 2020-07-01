package builder

import (
	"fmt"
	"strings"

	cfn "github.com/awslabs/goformation/v4/cloudformation"
	ec2 "github.com/awslabs/goformation/v4/cloudformation/ec2"
	cfnt "github.com/awslabs/goformation/v4/cloudformation/tags"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/vpc"
)

var internetCIDR = "0.0.0.0/0"

const (
	cfnControlPlaneSGResource         = "ControlPlaneSecurityGroup"
	cfnSharedNodeSGResource           = "ClusterSharedNodeSecurityGroup"
	cfnIngressClusterToNodeSGResource = "IngressDefaultClusterToNodeSG"
)

func (c *ClusterResourceSet) addSubnets(refRT string, topology api.SubnetTopology, subnets map[string]api.Network) {
	var subnetIndexForIPv6 int
	if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
		// this is same kind of indexing we have in vpc.SetSubnets
		switch topology {
		case api.SubnetTopologyPrivate:
			subnetIndexForIPv6 = len(c.spec.AvailabilityZones)
		case api.SubnetTopologyPublic:
			subnetIndexForIPv6 = 0
		}
	}

	for az, subnet := range subnets {
		alias := string(topology) + strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))
		subnet := &ec2.Subnet{
			AvailabilityZone: az,
			CidrBlock:        subnet.CIDR.String(),
			VpcId:            c.vpc,
		}

		switch topology {
		case api.SubnetTopologyPrivate:
			// Choose the appropriate route table for private subnets
			refRT = cfn.Ref("PrivateRouteTable" + strings.ToUpper(strings.Join(strings.Split(az, "-"), "")))
			subnet.Tags = []cfnt.Tag{{
				Key:   "kubernetes.io/role/internal-elb",
				Value: "1",
			}}
		case api.SubnetTopologyPublic:
			subnet.Tags = []cfnt.Tag{{
				Key:   "kubernetes.io/role/elb",
				Value: "1",
			}}
			subnet.MapPublicIpOnLaunch = true
		}
		refSubnet := c.newResourceV4("Subnet"+alias, subnet)
		c.newResourceV4("RouteTableAssociation"+alias, &ec2.SubnetRouteTableAssociation{
			SubnetId:     refSubnet,
			RouteTableId: refRT,
		})

		if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
			// get 8 of /64 subnets from the auto-allocated IPv6 block,
			// and pick one block based on subnetIndexForIPv6 counter;
			// NOTE: this is done inside of CloudFormation using Fn::Cidr,
			// we don't slice it here, just construct the JSON expression
			// that does slicing at runtime.
			refAutoAllocateCIDRv6 := cfn.Select(
				0, []string{cfn.GetAtt("VPC", "Ipv6CidrBlocks")},
			)
			refSubnetSlices := cfn.CIDR(
				refAutoAllocateCIDRv6, 8, 64,
			)
			c.newResourceV4(alias+"CIDRv6", &ec2.SubnetCidrBlock{
				SubnetId:      refSubnet,
				Ipv6CidrBlock: cfn.Select(subnetIndexForIPv6, []string{refSubnetSlices}),
			})
			subnetIndexForIPv6++
		}

		c.subnets[topology] = append(c.subnets[topology], refSubnet)
	}
}

//nolint:interfacer
func (c *ClusterResourceSet) addResourcesForVPC() error {

	c.vpc = c.newResourceV4("VPC", &ec2.VPC{
		CidrBlock:          c.spec.VPC.CIDR.String(),
		EnableDnsSupport:   true,
		EnableDnsHostnames: true,
	})

	if api.IsEnabled(c.spec.VPC.AutoAllocateIPv6) {
		c.newResourceV4("AutoAllocatedCIDRv6", &ec2.VPCCidrBlock{
			VpcId:                       c.vpc,
			AmazonProvidedIpv6CidrBlock: true,
		})
	}

	c.subnets = make(map[api.SubnetTopology][]string)

	refIG := c.newResourceV4("InternetGateway", &ec2.InternetGateway{})
	vpcGA := "VPCGatewayAttachment"
	c.newResourceV4(vpcGA, &ec2.VPCGatewayAttachment{
		InternetGatewayId: refIG,
		VpcId:             c.vpc,
	})

	refPublicRT := c.newResourceV4("PublicRouteTable", &ec2.RouteTable{
		VpcId: c.vpc,
	})

	c.newResourceV4("PublicSubnetRoute", &ec2.Route{
		RouteTableId:               refPublicRT,
		DestinationCidrBlock:       internetCIDR,
		GatewayId:                  refIG,
		AWSCloudFormationDependsOn: []string{vpcGA},
	})

	c.addSubnets(refPublicRT, api.SubnetTopologyPublic, c.spec.VPC.Subnets.Public)

	if err := c.addNATGateways(); err != nil {
		return err
	}

	c.addSubnets("", api.SubnetTopologyPrivate, c.spec.VPC.Subnets.Private)
	return nil
}

func (c *ClusterResourceSet) addNATGateways() error {

	switch *c.spec.VPC.NAT.Gateway {

	case api.ClusterHighlyAvailableNAT:
		c.haNAT()
	case api.ClusterSingleNAT:
		c.singleNAT()
	case api.ClusterDisableNAT:
		c.noNAT()
	default:
		// TODO validate this before starting to add resources
		return fmt.Errorf("%s is not a valid NAT gateway mode", *c.spec.VPC.NAT.Gateway)
	}
	return nil
}

func (c *ClusterResourceSet) importResourcesForVPC() {
	c.vpc = cfn.Ref(c.spec.VPC.ID)
	c.subnets = make(map[api.SubnetTopology][]string)
	for _, subnet := range c.spec.PrivateSubnetIDs() {
		c.subnets[api.SubnetTopologyPrivate] = append(c.subnets[api.SubnetTopologyPrivate], subnet)
	}
	for _, subnet := range c.spec.PublicSubnetIDs() {
		c.subnets[api.SubnetTopologyPublic] = append(c.subnets[api.SubnetTopologyPublic], subnet)
	}

}

func (c *ClusterResourceSet) addOutputsForVPC() {
	if c.spec.VPC == nil {
		c.spec.VPC = &api.ClusterVPC{}
	}
	c.rs.defineOutput(outputs.ClusterVPC, c.vpc, true, func(v string) error {
		c.spec.VPC.ID = v
		return nil
	})
	if c.spec.VPC.NAT != nil {
		c.rs.defineOutputWithoutCollector(outputs.ClusterFeatureNATMode, c.spec.VPC.NAT.Gateway, false)
	}
	if refs, ok := c.subnets[api.SubnetTopologyPrivate]; ok {
		c.rs.defineJoinedOutputV4(outputs.ClusterSubnetsPrivate, refs, true, func(v string) error {
			return vpc.ImportSubnetsFromList(c.provider, c.spec, api.SubnetTopologyPrivate, strings.Split(v, ","))
		})
	}
	if refs, ok := c.subnets[api.SubnetTopologyPublic]; ok {
		c.rs.defineJoinedOutputV4(outputs.ClusterSubnetsPublic, refs, true, func(v string) error {
			return vpc.ImportSubnetsFromList(c.provider, c.spec, api.SubnetTopologyPublic, strings.Split(v, ","))
		})
	}
}

var (
	sgProtoTCP           = "tcp"
	sgSourceAnywhereIPv4 = "0.0.0.0/0"
	sgSourceAnywhereIPv6 = "::/0"

	sgPortZero    = 0
	sgMinNodePort = 1025
	sgMaxNodePort = 65535

	sgPortHTTPS = 443
	sgPortSSH   = 22
)

func (c *ClusterResourceSet) addResourcesForSecurityGroups() {
	var refControlPlaneSG, refClusterSharedNodeSG string

	if c.spec.VPC.SecurityGroup == "" {
		refControlPlaneSG = c.newResourceV4(cfnControlPlaneSGResource, &ec2.SecurityGroup{
			GroupDescription: "Communication between the control plane and worker nodegroups",
			VpcId:            c.vpc,
		})
	} else {
		refControlPlaneSG = c.spec.VPC.SecurityGroup
	}
	c.securityGroups = []string{refControlPlaneSG} // only this one SG is passed to EKS API, nodes are isolated

	if c.spec.VPC.SharedNodeSecurityGroup == "" {
		refClusterSharedNodeSG = c.newResourceV4(cfnSharedNodeSGResource, &ec2.SecurityGroup{
			GroupDescription: "Communication between all nodes in the cluster",
			VpcId:            c.vpc,
		})
		c.newResourceV4("IngressInterNodeGroupSG", &ec2.SecurityGroupIngress{
			GroupId:               refClusterSharedNodeSG,
			SourceSecurityGroupId: refClusterSharedNodeSG,
			Description:           "Allow nodes to communicate with each other (all ports)",
			IpProtocol:            "-1",
			FromPort:              sgPortZero,
			ToPort:                sgMaxNodePort,
		})
		if c.supportsManagedNodes {
			// To enable communication between both managed and unmanaged nodegroups, this allows ingress traffic from
			// the default cluster security group ID that EKS creates by default
			// EKS attaches this to Managed Nodegroups by default, but we need to handle this for unmanaged nodegroups
			c.newResourceV4(cfnIngressClusterToNodeSGResource, &ec2.SecurityGroupIngress{
				GroupId:               refClusterSharedNodeSG,
				SourceSecurityGroupId: cfn.GetAtt("ControlPlane", outputs.ClusterDefaultSecurityGroup),
				Description:           "Allow managed and unmanaged nodes to communicate with each other (all ports)",
				IpProtocol:            "-1",
				FromPort:              sgPortZero,
				ToPort:                sgMaxNodePort,
			})
			c.newResourceV4("IngressNodeToDefaultClusterSG", &ec2.SecurityGroupIngress{
				GroupId:               cfn.GetAtt("ControlPlane", outputs.ClusterDefaultSecurityGroup),
				SourceSecurityGroupId: refClusterSharedNodeSG,
				Description:           "Allow unmanaged nodes to communicate with control plane (all ports)",
				IpProtocol:            "-1",
				FromPort:              sgPortZero,
				ToPort:                sgMaxNodePort,
			})
		}
	} else {
		refClusterSharedNodeSG = c.spec.VPC.SharedNodeSecurityGroup
	}

	if c.spec.VPC == nil {
		c.spec.VPC = &api.ClusterVPC{}
	}
	c.rs.defineOutput(outputs.ClusterSecurityGroup, refControlPlaneSG, true, func(v string) error {
		c.spec.VPC.SecurityGroup = v
		return nil
	})
	c.rs.defineOutput(outputs.ClusterSharedNodeSecurityGroup, refClusterSharedNodeSG, true, func(v string) error {
		c.spec.VPC.SharedNodeSecurityGroup = v
		return nil
	})
}

func (n *NodeGroupResourceSet) addResourcesForSecurityGroups() {
	for _, id := range n.spec.SecurityGroups.AttachIDs {
		n.securityGroups = append(n.securityGroups, id)
	}

	if api.IsEnabled(n.spec.SecurityGroups.WithShared) {
		refClusterSharedNodeSG := makeImportValueV4(n.clusterStackName, outputs.ClusterSharedNodeSecurityGroup)
		n.securityGroups = append(n.securityGroups, refClusterSharedNodeSG)
	}

	if api.IsDisabled(n.spec.SecurityGroups.WithLocal) {
		return
	}

	desc := "worker nodes in group " + n.nodeGroupName

	allInternalIPv4 := n.clusterSpec.VPC.CIDR.String()

	refControlPlaneSG := makeImportValueV4(n.clusterStackName, outputs.ClusterSecurityGroup)

	refNodeGroupLocalSG := n.newResourceV4("SG", &ec2.SecurityGroup{
		VpcId:            makeImportValueV4(n.clusterStackName, outputs.ClusterVPC),
		GroupDescription: "Communication between the control plane and " + desc,
		Tags: []cfnt.Tag{{
			Key:   "kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name,
			Value: "owned",
		}},
	})

	n.securityGroups = append(n.securityGroups, refNodeGroupLocalSG)

	n.newResourceV4("IngressInterCluster", &ec2.SecurityGroupIngress{
		GroupId:               refNodeGroupLocalSG,
		SourceSecurityGroupId: refControlPlaneSG,
		Description:           "Allow " + desc + " to communicate with control plane (kubelet and workload TCP ports)",
		IpProtocol:            sgProtoTCP,
		FromPort:              sgMinNodePort,
		ToPort:                sgMaxNodePort,
	})
	n.newResourceV4("EgressInterCluster", &ec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                "Allow control plane to communicate with " + desc + " (kubelet and workload TCP ports)",
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgMinNodePort,
		ToPort:                     sgMaxNodePort,
	})
	n.newResourceV4("IngressInterClusterAPI", &ec2.SecurityGroupIngress{
		GroupId:               refNodeGroupLocalSG,
		SourceSecurityGroupId: refControlPlaneSG,
		Description:           "Allow " + desc + " to communicate with control plane (workloads using HTTPS port, commonly used with extension API servers)",
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
	n.newResourceV4("EgressInterClusterAPI", &ec2.SecurityGroupEgress{
		GroupId:                    refControlPlaneSG,
		DestinationSecurityGroupId: refNodeGroupLocalSG,
		Description:                "Allow control plane to communicate with " + desc + " (workloads using HTTPS port, commonly used with extension API servers)",
		IpProtocol:                 sgProtoTCP,
		FromPort:                   sgPortHTTPS,
		ToPort:                     sgPortHTTPS,
	})
	n.newResourceV4("IngressInterClusterCP", &ec2.SecurityGroupIngress{
		GroupId:               refControlPlaneSG,
		SourceSecurityGroupId: refNodeGroupLocalSG,
		Description:           "Allow control plane to receive API requests from " + desc,
		IpProtocol:            sgProtoTCP,
		FromPort:              sgPortHTTPS,
		ToPort:                sgPortHTTPS,
	})
	if *n.spec.SSH.Allow {
		if n.spec.PrivateNetworking {
			n.newResourceV4("SSHIPv4", &ec2.SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIp:      allInternalIPv4,
				Description: "Allow SSH access to " + desc + " (private, only inside VPC)",
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
		} else {
			n.newResourceV4("SSHIPv4", &ec2.SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIp:      sgSourceAnywhereIPv4,
				Description: "Allow SSH access to " + desc,
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
			n.newResourceV4("SSHIPv6", &ec2.SecurityGroupIngress{
				GroupId:     refNodeGroupLocalSG,
				CidrIpv6:    sgSourceAnywhereIPv6,
				Description: "Allow SSH access to " + desc,
				IpProtocol:  sgProtoTCP,
				FromPort:    sgPortSSH,
				ToPort:      sgPortSSH,
			})
		}
	}
}

func (c *ClusterResourceSet) haNAT() {

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		// Allocate an EIP
		c.newResourceV4("NATIP"+alphanumericUpperAZ, &ec2.EIP{
			Domain: "vpc",
		})
		// Allocate a NAT gateway in the public subnet
		refNG := c.newResourceV4("NATGateway"+alphanumericUpperAZ, &ec2.NatGateway{
			AllocationId: cfn.GetAtt("NATIP"+alphanumericUpperAZ, "AllocationId"),
			SubnetId:     cfn.Ref("SubnetPublic" + alphanumericUpperAZ),
		})

		// Allocate a routing table for the private subnet
		refRT := c.newResourceV4("PrivateRouteTable"+alphanumericUpperAZ, &ec2.RouteTable{
			VpcId: c.vpc,
		})
		// Create a route that sends Internet traffic through the NAT gateway
		c.newResourceV4("NATPrivateSubnetRoute"+alphanumericUpperAZ, &ec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		// Associate the routing table with the subnet
		c.newResourceV4("RouteTableAssociationPrivate"+alphanumericUpperAZ, &ec2.SubnetRouteTableAssociation{
			SubnetId:     cfn.Ref("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}

}

func (c *ClusterResourceSet) singleNAT() {

	sortedAZs := c.spec.AvailabilityZones
	firstUpperAZ := strings.ToUpper(strings.Join(strings.Split(sortedAZs[0], "-"), ""))

	c.newResourceV4("NATIP", &ec2.EIP{
		Domain: "vpc",
	})
	refNG := c.newResourceV4("NATGateway", &ec2.NatGateway{
		AllocationId: cfn.GetAtt("NATIP", "AllocationId"),
		SubnetId:     cfn.Ref("SubnetPublic" + firstUpperAZ),
	})

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := c.newResourceV4("PrivateRouteTable"+alphanumericUpperAZ, &ec2.RouteTable{
			VpcId: c.vpc,
		})

		c.newResourceV4("NATPrivateSubnetRoute"+alphanumericUpperAZ, &ec2.Route{
			RouteTableId:         refRT,
			DestinationCidrBlock: internetCIDR,
			NatGatewayId:         refNG,
		})
		c.newResourceV4("RouteTableAssociationPrivate"+alphanumericUpperAZ, &ec2.SubnetRouteTableAssociation{
			SubnetId:     cfn.Ref("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}

func (c *ClusterResourceSet) noNAT() {

	for _, az := range c.spec.AvailabilityZones {
		alphanumericUpperAZ := strings.ToUpper(strings.Join(strings.Split(az, "-"), ""))

		refRT := c.newResourceV4("PrivateRouteTable"+alphanumericUpperAZ, &ec2.RouteTable{
			VpcId: c.vpc,
		})
		c.newResourceV4("RouteTableAssociationPrivate"+alphanumericUpperAZ, &ec2.SubnetRouteTableAssociation{
			SubnetId:     cfn.Ref("SubnetPrivate" + alphanumericUpperAZ),
			RouteTableId: refRT,
		})
	}
}
