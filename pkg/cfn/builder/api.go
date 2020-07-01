package builder

import (
	"encoding/json"
	"fmt"
	"reflect"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfnv4 "github.com/awslabs/goformation/v4/cloudformation"
	gfncfn "github.com/awslabs/goformation/v4/cloudformation/cloudformation"
	"github.com/awslabs/goformation/v4/intrinsics"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	gfn "github.com/weaveworks/goformation/cloudformation"
)

const (
	clusterTemplateDescription   = "EKS cluster"
	nodeGroupTemplateDescription = "EKS nodes"
	templateDescriptionSuffix    = "[created and managed by eksctl]"
)

type awsCloudFormationResource struct {
	Type         string
	Properties   map[string]interface{}
	UpdatePolicy map[string]map[string]string `json:",omitempty"`
	DependsOn    []string                     `json:",omitempty"`
}

func (r *awsCloudFormationResource) AWSCloudFormationType() string {
	return r.Type
}

// ResourceSet is an interface which cluster and nodegroup builders
// must implement
type ResourceSet interface {
	AddAllResources() error
	WithIAM() bool
	WithNamedIAM() bool
	RenderJSON() ([]byte, error)
	GetAllOutputs(cfn.Stack) error
}

type resourceSet struct {
	template     *gfnv4.Template
	outputs      *outputs.CollectorSet
	withIAM      bool
	withNamedIAM bool
}

func newResourceSet() *resourceSet {
	return &resourceSet{
		template: gfnv4.NewTemplate(),
		outputs:  outputs.NewCollectorSet(nil),
	}
}

// makeName is syntactic sugar for {"Fn::Sub": "${AWS::Stack}-<name>"}
func makeName(suffix string) string {
	return gfnv4.Sub(fmt.Sprintf("${%s}-%s", gfn.StackName, suffix))
}

// makeSlice makes a slice from a list of arguments
func makeSlice(i ...string) []string {
	return i
}

// makeAutoNameTag create a new Name tag in the following format:
// {Key: "Name", Value: !Sub "${AWS::StackName}/<logicalResourceName>"}
func makeAutoNameTag(suffix string) gfncfn.Tag {
	return gfncfn.Tag{
		Key:   "Name",
		Value: gfnv4.Sub(fmt.Sprintf("${%s}/%s", gfn.StackName, suffix)),
	}
}

func makeAttrAccessor(resource, attr string) string {
	return fmt.Sprintf("%s.%s", resource, attr)
}

// maybeSetNameTag adds a Name tag to any resource that supports tags
// it calls makeAutoNameTag to format the tag value
func maybeSetNameTag(name string, resource interface{}) {
	e := reflect.ValueOf(resource).Elem()
	if e.Kind() == reflect.Struct {
		f := e.FieldByName("Tags")
		if f.IsValid() && f.CanSet() {
			tag := reflect.ValueOf(makeAutoNameTag(name))
			if f.Type() == reflect.ValueOf([]gfncfn.Tag{}).Type() {
				f.Set(reflect.Append(f, tag))
			}
		}
	}
}

// newResourceV4 adds a resource, and adds Name tag if possible, it returns a reference
func (r *resourceSet) newResourceV4(name string, resource gfnv4.Resource) string {
	r.template.Resources[name] = resource
	return gfnv4.Ref(name)
}

// renderJSON renders template as JSON
func (r *resourceSet) renderJSON() ([]byte, error) {
	j, err := json.Marshal(r.template)
	if err != nil {
		return []byte{}, err
	}
	opts := intrinsics.ProcessorOptions{
		IntrinsicHandlerOverrides: gfnv4.EncoderIntrinsics,
	}
	return intrinsics.ProcessJSON(j, &opts)
}
