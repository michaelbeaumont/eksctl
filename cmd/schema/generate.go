package main

import (
	"fmt"
	"io/ioutil"
	"os"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io"
	v1alpha5 "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	//"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

/*func typeMapper(t reflect.Type) *jsonschema.Type {
	if t == reflect.TypeOf(&ipnet.IPNet{}) {
		return &jsonschema.Type{Type: "string"}
	}
	return nil
}*/

func main() {
	if len(os.Args) != 2 {
		panic("expected one argument with the output file")
	}
	outputFile := os.Args[1]

	schema, err := GenerateSchema("../../../..", "v1alpha5", "ClusterConfig")
	if err != nil {
		panic(err)
	}

	// We add some examples and blacklist some descriptions
	if t, ok := schema.Definitions["ClusterConfig"].Properties["kind"]; ok {
		t.Examples = []string{"ClusterConfig"}
		t.Description = ""
	}
	if t, ok := schema.Definitions["ClusterConfig"].Properties["apiVersion"]; ok {
		t.Examples = []string{fmt.Sprintf("%s/%s", api.GroupName, v1alpha5.CurrentGroupVersion)}
		t.Description = ""
	}

	if err != nil {
		fmt.Println(err)
	}
	bytes, err := ToJSON(schema)
	if err != nil {
		fmt.Println(err)
	}
	err = ioutil.WriteFile(outputFile, bytes, 0755)

}
