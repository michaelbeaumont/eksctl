package builder

import (
	"fmt"
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/weaveworks/goformation/cloudformation"
	gfnv4 "github.com/awslabs/goformation/v4/cloudformation"
	ec2 "github.com/awslabs/goformation/v4/cloudformation/ec2"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

// NodeGroupResourceSet stores the resource information of the nodegroup
type NodeGroupResourceSet struct {
	rs                   *resourceSet
	clusterSpec          *api.ClusterConfig
	spec                 *api.NodeGroup
	supportsManagedNodes bool
	provider             api.ClusterProvider
	clusterStackName     string
	nodeGroupName        string
	instanceProfileARN   string
	securityGroups       []string
	vpc                  string
	userData             string
}

// NewNodeGroupResourceSet returns a resource set for a nodegroup embedded in a cluster config
func NewNodeGroupResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, clusterStackName string, ng *api.NodeGroup,
	supportsManagedNodes bool) *NodeGroupResourceSet {
	return &NodeGroupResourceSet{
		rs:                   newResourceSet(),
		clusterStackName:     clusterStackName,
		nodeGroupName:        ng.Name,
		supportsManagedNodes: supportsManagedNodes,
		clusterSpec:          spec,
		spec:                 ng,
		provider:             provider,
	}
}

// AddAllResources adds all the information about the nodegroup to the resource set
func (n *NodeGroupResourceSet) AddAllResources() error {
	n.rs.template.Description = fmt.Sprintf(
		"%s (AMI family: %s, SSH access: %v, private networking: %v) %s",
		nodeGroupTemplateDescription,
		n.spec.AMIFamily, api.IsEnabled(n.spec.SSH.Allow), n.spec.PrivateNetworking,
		templateDescriptionSuffix)

	n.Template().Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeaturePrivateNetworking, n.spec.PrivateNetworking, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureSharedSecurityGroup, n.spec.SecurityGroups.WithShared, false)
	n.rs.defineOutputWithoutCollector(outputs.NodeGroupFeatureLocalSecurityGroup, n.spec.SecurityGroups.WithLocal, false)

	n.vpc = makeImportValueV4(n.clusterStackName, outputs.ClusterVPC)

	userData, err := nodebootstrap.NewUserData(n.clusterSpec, n.spec)
	if err != nil {
		return err
	}
	n.userData = userData

	// Ensure MinSize is set, as it is required by the ASG cfn resource
	if n.spec.MinSize == nil {
		if n.spec.DesiredCapacity == nil {
			defaultNodeCount := api.DefaultNodeCount
			n.spec.MinSize = &defaultNodeCount
		} else {
			n.spec.MinSize = n.spec.DesiredCapacity
		}
		logger.Info("--nodes-min=%d was set automatically for nodegroup %s", *n.spec.MinSize, n.nodeGroupName)
	} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity < *n.spec.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes=%d at the same time", *n.spec.MinSize, *n.spec.DesiredCapacity)
	}

	// Ensure MaxSize is set, as it is required by the ASG cfn resource
	if n.spec.MaxSize == nil {
		if n.spec.DesiredCapacity == nil {
			n.spec.MaxSize = n.spec.MinSize
		} else {
			n.spec.MaxSize = n.spec.DesiredCapacity
		}
		logger.Info("--nodes-max=%d was set automatically for nodegroup %s", *n.spec.MaxSize, n.nodeGroupName)
	} else if n.spec.DesiredCapacity != nil && *n.spec.DesiredCapacity > *n.spec.MaxSize {
		return fmt.Errorf("cannot use --nodes-max=%d and --nodes=%d at the same time", *n.spec.MaxSize, *n.spec.DesiredCapacity)
	} else if *n.spec.MaxSize < *n.spec.MinSize {
		return fmt.Errorf("cannot use --nodes-min=%d and --nodes-max=%d at the same time", *n.spec.MinSize, *n.spec.MaxSize)
	}

	if err := n.addResourcesForIAM(); err != nil {
		return err
	}
	n.addResourcesForSecurityGroups()

	return n.addResourcesForNodeGroup()
}

// RenderJSON returns the rendered JSON
func (n *NodeGroupResourceSet) RenderJSON() ([]byte, error) {
	return n.rs.renderJSON()
}

// Template returns the CloudFormation template
func (n *NodeGroupResourceSet) Template() gfnv4.Template {
	return *n.rs.template
}

func (n *NodeGroupResourceSet) newResourceV4(name string, resource gfnv4.Resource) string {
	return n.rs.newResourceV4(name, resource)
}
func (n *NodeGroupResourceSet) addResourcesForNodeGroup() error {
	launchTemplateName := gfnv4.Sub(fmt.Sprintf("${%s}", gfn.StackName))
	launchTemplateData := newLaunchTemplateData(n)

	if n.spec.SSH != nil && api.IsSetAndNonEmptyString(n.spec.SSH.PublicKeyName) {
		launchTemplateData.KeyName = *n.spec.SSH.PublicKeyName
	}

	if volumeSize := n.spec.VolumeSize; volumeSize != nil && *volumeSize > 0 {
		var (
			kmsKeyID   string
			volumeIOPS int
		)
		if api.IsSetAndNonEmptyString(n.spec.VolumeKmsKeyID) {
			kmsKeyID = *n.spec.VolumeKmsKeyID
		}

		if *n.spec.VolumeType == api.NodeVolumeTypeIO1 {
			volumeIOPS = *n.spec.VolumeIOPS
		}

		launchTemplateData.BlockDeviceMappings = []ec2.LaunchTemplate_BlockDeviceMapping{{
			DeviceName: *n.spec.VolumeName,
			Ebs: &ec2.LaunchTemplate_Ebs{
				VolumeSize: *volumeSize,
				VolumeType: *n.spec.VolumeType,
				Encrypted:  *n.spec.VolumeEncrypted,
				KmsKeyId:   kmsKeyID,
				Iops:       volumeIOPS,
			},
		}}
	}

	n.newResourceV4("NodeGroupLaunchTemplate", &ec2.LaunchTemplate{
		LaunchTemplateName: launchTemplateName,
		LaunchTemplateData: launchTemplateData,
	})

	vpcZoneIdentifier, err := AssignSubnets(n.spec.AvailabilityZones, n.clusterStackName, n.clusterSpec, n.spec.PrivateNetworking)
	if err != nil {
		return err
	}

	tags := []map[string]interface{}{
		{
			"Key":               "Name",
			"Value":             n.generateNodeName(),
			"PropagateAtLaunch": "true",
		},
		{
			"Key":               "kubernetes.io/cluster/" + n.clusterSpec.Metadata.Name,
			"Value":             "owned",
			"PropagateAtLaunch": "true",
		},
	}
	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.AutoScaler) {
		tags = append(tags,
			map[string]interface{}{
				"Key":               "k8s.io/cluster-autoscaler/enabled",
				"Value":             "true",
				"PropagateAtLaunch": "true",
			},
			map[string]interface{}{
				"Key":               "k8s.io/cluster-autoscaler/" + n.clusterSpec.Metadata.Name,
				"Value":             "owned",
				"PropagateAtLaunch": "true",
			},
		)
	}

	asg := nodeGroupResource(launchTemplateName, vpcZoneIdentifier, tags, n.spec)
	n.newResourceV4("NodeGroup", asg)

	return nil
}

// generateNodeName formulates the name based on the configuration in input
func (n *NodeGroupResourceSet) generateNodeName() string {
	name := []string{}
	if n.spec.InstancePrefix != "" {
		name = append(name, n.spec.InstancePrefix, "-")
	}
	// this overrides the default naming convention
	if n.spec.InstanceName != "" {
		name = append(name, n.spec.InstanceName)
	} else {
		name = append(name, fmt.Sprintf("%s-%s-Node", n.clusterSpec.Metadata.Name, n.nodeGroupName))
	}
	return strings.Join(name, "")
}

// AssignSubnets subnets based on the specified availability zones
func AssignSubnets(availabilityZones []string, clusterStackName string, clusterSpec *api.ClusterConfig, privateNetworking bool) (interface{}, error) {
	// currently goformation type system doesn't allow specifying `VPCZoneIdentifier: { "Fn::ImportValue": ... }`,
	// and tags don't have `PropagateAtLaunch` field, so we have a custom method here until this gets resolved

	if numNodeGroupsAZs := len(availabilityZones); numNodeGroupsAZs > 0 {
		subnets := clusterSpec.VPC.Subnets.Private
		if !privateNetworking {
			subnets = clusterSpec.VPC.Subnets.Public
		}
		makeErrorDesc := func() string {
			return fmt.Sprintf("(subnets=%#v AZs=%#v)", subnets, availabilityZones)
		}
		if len(subnets) < numNodeGroupsAZs {
			return nil, fmt.Errorf("VPC doesn't have enough subnets for nodegroup AZs %s", makeErrorDesc())
		}
		subnetIDs := make([]string, numNodeGroupsAZs)
		for i, az := range availabilityZones {
			subnet, ok := subnets[az]
			if !ok {
				return nil, fmt.Errorf("VPC doesn't have subnets in %s %s", az, makeErrorDesc())
			}

			subnetIDs[i] = subnet.ID
		}
		return subnetIDs, nil
	}

	var subnets string
	if privateNetworking {
		subnets = makeImportValueV4(clusterStackName, outputs.ClusterSubnetsPrivate)
	} else {
		subnets = makeImportValueV4(clusterStackName, outputs.ClusterSubnetsPublic)
	}

	return gfnv4.Split(",", subnets), nil
}

// GetAllOutputs collects all outputs of the nodegroup
func (n *NodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack)
}

func newLaunchTemplateData(n *NodeGroupResourceSet) *ec2.LaunchTemplate_LaunchTemplateData {
	launchTemplateData := &ec2.LaunchTemplate_LaunchTemplateData{
		IamInstanceProfile: &ec2.LaunchTemplate_IamInstanceProfile{
			Arn: n.instanceProfileARN,
		},
		ImageId:  n.spec.AMI,
		UserData: n.userData,
		NetworkInterfaces: []ec2.LaunchTemplate_NetworkInterface{{
			// Explicitly un-setting this so that it doesn't get defaulted to true
			AssociatePublicIpAddress: false,
			DeviceIndex:              0,
			Groups:                   n.securityGroups,
		}},
		MetadataOptions: &ec2.LaunchTemplate_MetadataOptions{
			HttpPutResponseHopLimit: 2,
		},
	}
	if !api.HasMixedInstances(n.spec) {
		launchTemplateData.InstanceType = n.spec.InstanceType
	} else {
		launchTemplateData.InstanceType = n.spec.InstancesDistribution.InstanceTypes[0]
	}
	if n.spec.EBSOptimized != nil {
		launchTemplateData.EbsOptimized = *n.spec.EBSOptimized
	}

	return launchTemplateData
}

func nodeGroupResource(launchTemplateName string, vpcZoneIdentifier interface{}, tags []map[string]interface{}, ng *api.NodeGroup) *awsCloudFormationResource {
	ngProps := map[string]interface{}{
		"VPCZoneIdentifier": vpcZoneIdentifier,
		"Tags":              tags,
	}
	if ng.DesiredCapacity != nil {
		ngProps["DesiredCapacity"] = fmt.Sprintf("%d", *ng.DesiredCapacity)
	}
	if ng.MinSize != nil {
		ngProps["MinSize"] = fmt.Sprintf("%d", *ng.MinSize)
	}
	if ng.MaxSize != nil {
		ngProps["MaxSize"] = fmt.Sprintf("%d", *ng.MaxSize)
	}
	if len(ng.ASGMetricsCollection) > 0 {
		ngProps["MetricsCollection"] = metricsCollectionResource(ng.ASGMetricsCollection)
	}
	if len(ng.ClassicLoadBalancerNames) > 0 {
		ngProps["LoadBalancerNames"] = ng.ClassicLoadBalancerNames
	}
	if len(ng.TargetGroupARNs) > 0 {
		ngProps["TargetGroupARNs"] = ng.TargetGroupARNs
	}
	if api.HasMixedInstances(ng) {
		ngProps["MixedInstancesPolicy"] = *mixedInstancesPolicy(launchTemplateName, ng)
	} else {
		ngProps["LaunchTemplate"] = map[string]interface{}{
			"LaunchTemplateName": launchTemplateName,
			"Version":            gfnv4.GetAtt("NodeGroupLaunchTemplate", "LatestVersionNumber"),
		}
	}

	return &awsCloudFormationResource{
		Type:       "AWS::AutoScaling::AutoScalingGroup",
		Properties: ngProps,
		UpdatePolicy: map[string]map[string]string{
			"AutoScalingRollingUpdate": {
				"MinInstancesInService": "0",
				"MaxBatchSize":          "1",
			},
		},
	}
}

func mixedInstancesPolicy(launchTemplateName string, ng *api.NodeGroup) *map[string]interface{} {
	instanceTypes := ng.InstancesDistribution.InstanceTypes
	overrides := make([]map[string]string, len(instanceTypes))

	for i, instanceType := range instanceTypes {
		overrides[i] = map[string]string{
			"InstanceType": instanceType,
		}
	}
	policy := map[string]interface{}{
		"LaunchTemplate": map[string]interface{}{
			"LaunchTemplateSpecification": map[string]interface{}{
				"LaunchTemplateName": launchTemplateName,
				"Version":            gfnv4.GetAtt("NodeGroupLaunchTemplate", "LatestVersionNumber"),
			},

			"Overrides": overrides,
		},
	}

	instancesDistribution := map[string]string{}

	// Only set the price if it was specified so otherwise AWS picks "on-demand price" as the default
	if ng.InstancesDistribution.MaxPrice != nil {
		instancesDistribution["SpotMaxPrice"] = fmt.Sprintf("%f", *ng.InstancesDistribution.MaxPrice)
	}
	if ng.InstancesDistribution.OnDemandBaseCapacity != nil {
		instancesDistribution["OnDemandBaseCapacity"] = fmt.Sprintf("%d", *ng.InstancesDistribution.OnDemandBaseCapacity)
	}
	if ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity != nil {
		instancesDistribution["OnDemandPercentageAboveBaseCapacity"] = fmt.Sprintf("%d", *ng.InstancesDistribution.OnDemandPercentageAboveBaseCapacity)
	}
	if ng.InstancesDistribution.SpotInstancePools != nil {
		instancesDistribution["SpotInstancePools"] = fmt.Sprintf("%d", *ng.InstancesDistribution.SpotInstancePools)
	}

	if ng.InstancesDistribution.SpotAllocationStrategy != nil {
		instancesDistribution["SpotAllocationStrategy"] = *ng.InstancesDistribution.SpotAllocationStrategy
	}

	policy["InstancesDistribution"] = instancesDistribution

	return &policy
}

func metricsCollectionResource(asgMetricsCollection []api.MetricsCollection) []map[string]interface{} {
	var metricsCollections []map[string]interface{}
	for _, m := range asgMetricsCollection {
		newCollection := make(map[string]interface{})

		if len(m.Metrics) > 0 {
			newCollection["Metrics"] = m.Metrics
		}
		newCollection["Granularity"] = m.Granularity

		metricsCollections = append(metricsCollections, newCollection)
	}
	return metricsCollections
}
