package builder

import (
	"fmt"
	"reflect"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfnv4 "github.com/awslabs/goformation/v4/cloudformation"
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
	template     *gfn.Template
	outputs      *outputs.CollectorSet
	withIAM      bool
	withNamedIAM bool
}

func newResourceSet() *resourceSet {
	return &resourceSet{
		template: gfn.NewTemplate(),
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

// makeSlice makes a slice from a list of string arguments
func makeStringSlice(s ...string) []*gfn.Value {
	slice := []*gfn.Value{}
	for _, i := range s {
		slice = append(slice, gfn.NewString(i))
	}
	return slice
}

// makeAutoNameTag create a new Name tag in the following format:
// {Key: "Name", Value: !Sub "${AWS::StackName}/<logicalResourceName>"}
func makeAutoNameTag(suffix string) gfn.Tag {
	return gfn.Tag{
		Key:   gfn.NewString("Name"),
		Value: gfn.MakeFnSubString(fmt.Sprintf("${%s}/%s", gfn.StackName, suffix)),
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
			if f.Type() == reflect.ValueOf([]gfn.Tag{}).Type() {
				f.Set(reflect.Append(f, tag))
			}
		}
	}
}

// newResource adds a resource, and adds Name tag if possible, it returns a reference
func (r *resourceSet) newResource(name string, resource interface{}) *gfn.Value {
	maybeSetNameTag(name, resource)
	r.template.Resources[name] = resource
	return gfn.MakeRef(name)
}

// newResourceV4 adds a resource, and adds Name tag if possible, it returns a reference
func (r *resourceSet) newResourceV4(name string, resource interface{}) string {
	r.template.Resources[name] = resource
	return gfnv4.Ref(name)
}

// renderJSON renders template as JSON
func (r *resourceSet) renderJSON() ([]byte, error) {
	js, err := r.template.JSON()
	if err != nil {
		return []byte{}, err
	}
	opts := intrinsics.ProcessorOptions{
		IntrinsicHandlerOverrides: gfnv4.EncoderIntrinsics,
	}
	return intrinsics.ProcessJSON(js, &opts)
}
