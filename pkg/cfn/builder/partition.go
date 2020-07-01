package builder

import (
	"fmt"

	gfnv4 "github.com/awslabs/goformation/v4/cloudformation"
	gfn "github.com/weaveworks/goformation/cloudformation"
)

var servicePrincipalPartitionMappings = map[string]map[string]string{
	"aws": {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	"aws-us-gov": {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	"aws-cn": {
		"EC2":            "ec2.amazonaws.com.cn",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
}

const servicePrincipalPartitionMapName = "ServicePrincipalPartitionMap"

// MakeServiceRef returns a reference to an intrinsic map function that looks up the servicePrincipalName
// in servicePrincipalPartitionMappings
func MakeServiceRef(servicePrincipalName string) string {
	return gfnv4.FindInMap(servicePrincipalPartitionMapName, gfnv4.Ref(gfn.Partition), servicePrincipalName)
}

func makePolicyARNs(policyNames ...string) []string {
	policyARNs := make([]string, len(policyNames))
	for i, policy := range policyNames {
		policyARNs[i] = gfnv4.Sub(fmt.Sprintf("arn:${%s}:iam::aws:policy/%s", gfn.Partition, policy))
	}
	return policyARNs
}

func addARNPartitionPrefix(s string) *gfn.Value {
	return gfn.MakeFnSubString(fmt.Sprintf("arn:${%s}:%s", gfn.Partition, s))
}
