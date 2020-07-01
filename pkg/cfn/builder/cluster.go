package builder

import (
	"encoding/base64"
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfnv4 "github.com/awslabs/goformation/v4/cloudformation"
	gfneks "github.com/awslabs/goformation/v4/cloudformation/eks"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	gfn "github.com/weaveworks/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// ClusterResourceSet stores the resource information of the cluster
type ClusterResourceSet struct {
	rs                   *resourceSet
	spec                 *api.ClusterConfig
	provider             api.ClusterProvider
	supportsManagedNodes bool
	vpc                  string
	subnets              map[api.SubnetTopology][]string
	securityGroups       []string
}

// NewClusterResourceSet returns a resource set for the new cluster
func NewClusterResourceSet(provider api.ClusterProvider, spec *api.ClusterConfig, supportsManagedNodes bool, existingStack *gjson.Result) *ClusterResourceSet {
	if existingStack != nil {
		unsetExistingResources(existingStack, spec)
	}
	return &ClusterResourceSet{
		rs:                   newResourceSet(),
		spec:                 spec,
		provider:             provider,
		supportsManagedNodes: supportsManagedNodes,
	}
}

// unsetExistingResources unsets fields for CloudFormation resources that were created by eksctl (and not user-supplied)
// in order to trigger execution of code that relies on these fields
func unsetExistingResources(existingStack *gjson.Result, clusterConfig *api.ClusterConfig) {
	controlPlaneSG := existingStack.Get(cfnControlPlaneSGResource)
	if controlPlaneSG.Exists() {
		clusterConfig.VPC.SecurityGroup = ""
	}
	sharedNodeSG := existingStack.Get(cfnSharedNodeSGResource)
	if sharedNodeSG.Exists() {
		clusterConfig.VPC.SharedNodeSecurityGroup = ""
	}
}

// AddAllResources adds all the information about the cluster to the resource set
func (c *ClusterResourceSet) AddAllResources() error {
	dedicatedVPC := c.spec.VPC.ID == ""

	if err := c.spec.HasSufficientSubnets(); err != nil {
		return err
	}

	if dedicatedVPC {
		if err := c.addResourcesForVPC(); err != nil {
			return errors.Wrap(err, "error adding VPC resources")
		}
	} else {
		c.importResourcesForVPC()
	}
	c.addOutputsForVPC()

	c.addResourcesForSecurityGroups()
	c.addResourcesForIAM()
	c.addResourcesForControlPlane()

	if len(c.spec.FargateProfiles) > 0 {
		c.addResourcesForFargate()
	}

	c.rs.defineOutput(outputs.ClusterStackName, gfnv4.Ref(gfn.StackName), false, func(v string) error {
		if c.spec.Status == nil {
			c.spec.Status = &api.ClusterStatus{}
		}
		c.spec.Status.StackName = v
		return nil
	})

	c.Template().Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	c.rs.template.Description = fmt.Sprintf(
		"%s (dedicated VPC: %v, dedicated IAM: %v) %s",
		clusterTemplateDescription,
		dedicatedVPC, c.rs.withIAM,
		templateDescriptionSuffix)

	return nil
}

// RenderJSON returns the rendered JSON
func (c *ClusterResourceSet) RenderJSON() ([]byte, error) {
	return c.rs.renderJSON()
}

// Template returns the CloudFormation template
func (c *ClusterResourceSet) Template() gfnv4.Template {
	return *c.rs.template
}

// HasManagedNodesSG reports whether the stack has the security group required for communication between
// managed and unmanaged nodegroups
func HasManagedNodesSG(stackResources *gjson.Result) bool {
	return stackResources.Get(cfnIngressClusterToNodeSGResource).Exists()
}

func (c *ClusterResourceSet) newResourceV4(name string, resource gfnv4.Resource) string {
	return c.rs.newResourceV4(name, resource)
}

type encryptionProvider struct {
	KeyArn string `json:"KeyArn"`
}

type encryptionConfig struct {
	Provider  *encryptionProvider `json:"Provider"`
	Resources []string            `json:"Resources"`
}

type awsEKSCluster gfneks.Cluster

func (c *ClusterResourceSet) addResourcesForControlPlane() {
	clusterVPC := &gfneks.Cluster_ResourcesVpcConfig{
		SecurityGroupIds: c.securityGroups,
	}
	for topology := range c.subnets {
		var tops []string
		for _, c := range c.subnets[topology] {
			tops = append(tops, c)
		}
		clusterVPC.SubnetIds = append(clusterVPC.SubnetIds, tops...)
	}

	serviceRoleARN := gfnv4.GetAtt("ServiceRole", "Arn")
	if api.IsSetAndNonEmptyString(c.spec.IAM.ServiceRoleARN) {
		serviceRoleARN = *c.spec.IAM.ServiceRoleARN
	}

	var encryptionConfigs []gfneks.Cluster_EncryptionConfig
	if c.spec.SecretsEncryption != nil && c.spec.SecretsEncryption.KeyARN != nil {
		encryptionConfigs = []gfneks.Cluster_EncryptionConfig{
			{
				Resources: []string{"secrets"},
				Provider: &gfneks.Cluster_Provider{
					KeyArn: *c.spec.SecretsEncryption.KeyARN,
				},
			},
		}
	}

	c.newResourceV4("ControlPlane", &gfneks.Cluster{
		Name:               c.spec.Metadata.Name,
		RoleArn:            serviceRoleARN,
		Version:            c.spec.Metadata.Version,
		ResourcesVpcConfig: clusterVPC,
		EncryptionConfig:   encryptionConfigs,
	})

	if c.spec.Status == nil {
		c.spec.Status = &api.ClusterStatus{}
	}

	c.rs.defineOutputFromAttV4(outputs.ClusterCertificateAuthorityData, "ControlPlane", "CertificateAuthorityData", false, func(v string) error {
		caData, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return errors.Wrap(err, "decoding certificate authority data")
		}
		c.spec.Status.CertificateAuthorityData = caData
		return nil
	})
	c.rs.defineOutputFromAttV4(outputs.ClusterEndpoint, "ControlPlane", "Endpoint", true, func(v string) error {
		c.spec.Status.Endpoint = v
		return nil
	})
	c.rs.defineOutputFromAttV4(outputs.ClusterARN, "ControlPlane", "Arn", true, func(v string) error {
		c.spec.Status.ARN = v
		return nil
	})

	if c.supportsManagedNodes {
		// This exports the cluster security group ID that EKS creates by default. To enable communication between both
		// managed and unmanaged nodegroups, they must share a security group.
		// EKS attaches this to Managed Nodegroups by default, but we need to add this for unmanaged nodegroups.
		// This exported value is imported in the CloudFormation resource for unmanaged nodegroups
		c.rs.defineOutputFromAttV4(outputs.ClusterDefaultSecurityGroup, "ControlPlane", "ClusterSecurityGroupId",
			true, func(s string) error {
				return nil
			})
	}
}

func (c *ClusterResourceSet) addResourcesForFargate() {
	_ = AddResourcesForFargate(c.rs, c.spec)
}

// GetAllOutputs collects all outputs of the cluster
func (c *ClusterResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return c.rs.GetAllOutputs(stack)
}
