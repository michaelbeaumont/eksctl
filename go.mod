// Make sure to run the following commands after changes to this file are made:
// ` make -f Makefile.docker update-build-image-tag && make -f Makefile.docker push-build-image`
module github.com/weaveworks/eksctl

go 1.15

require (
	github.com/Djarvur/go-err113 v0.1.0 // indirect
	github.com/aws/aws-sdk-go v1.38.35
	github.com/benjamintf1/unmarshalledmatchers v0.0.0-20190408201839-bb1c1f34eaea
	github.com/blang/semver v3.5.1+incompatible
	github.com/bxcodec/faker v2.0.1+incompatible
	github.com/cloudflare/cfssl v1.5.0
	github.com/dave/jennifer v1.4.1
	github.com/dlespiau/kube-test-harness v0.0.0-20200915102055-a03579200ae8
	github.com/evanphx/json-patch/v5 v5.2.0
	github.com/fluxcd/flux/pkg/install v0.0.0-20201001122558-cb08da1b356a // flux 1.21.0
	github.com/fluxcd/go-git-providers v0.0.3
	github.com/fluxcd/helm-operator/pkg/install v0.0.0-20200729150005-1467489f7ee4 // helm-operator 1.2.0
	github.com/github-release/github-release v0.10.0
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/flock v0.8.0
	github.com/golangci/golangci-lint v1.37.1
	github.com/golangci/misspell v0.3.5 // indirect
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e // indirect
	github.com/goreleaser/goreleaser v0.172.1
	github.com/gostaticanalysis/analysisutil v0.6.1 // indirect
	github.com/instrumenta/kubeval v0.0.0-20190918223246-8d013ec9fc56
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/justinbarrick/go-k8s-portforward v1.0.4-0.20200904152830-b575325c1855
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kevinburke/go-bindata v3.22.0+incompatible
	github.com/kevinburke/rest v0.0.0-20210106114233-22cd0577e450 // indirect
	github.com/kubicorn/kubicorn v0.0.0-20180829191017-06f6bce92acc
	github.com/lithammer/dedent v1.1.0
	github.com/matoous/godox v0.0.0-20200801072554-4fb83dc2941e // indirect
	github.com/maxbrunsfeld/counterfeiter/v6 v6.3.0
	github.com/nxadm/tail v1.4.6 // indirect
	github.com/onsi/ginkgo v1.15.0
	github.com/onsi/gomega v1.10.5
	github.com/pelletier/go-toml v1.8.1
	github.com/pkg/errors v0.9.1
	github.com/quasilyte/regex/syntax v0.0.0-20200805063351-8f842688393c // indirect
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/spf13/afero v1.5.1
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tdakkota/asciicheck v0.0.0-20200416200610-e657995f937b // indirect
	github.com/tidwall/gjson v1.6.8
	github.com/tidwall/sjson v1.1.5
	github.com/timakin/bodyclose v0.0.0-20200424151742-cb6215831a94 // indirect
	github.com/tomarrell/wrapcheck v0.0.0-20201130113247-1683564d9756 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/vektra/mockery v1.1.2
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2 // indirect
	github.com/weaveworks/goformation/v4 v4.10.2-0.20210202192510-c984c16fe84b
	github.com/weaveworks/launcher v0.0.2-0.20200715141516-1ca323f1de15
	github.com/weaveworks/logger v0.0.0-20210210175120-de9359622dfc
	github.com/whilp/git-urls v0.0.0-20191001220047-6db9661140c0
	golang.org/x/tools v0.1.0
	k8s.io/api v0.19.5
	k8s.io/apiextensions-apiserver v0.19.5
	k8s.io/apimachinery v0.19.5
	k8s.io/cli-runtime v0.19.5
	k8s.io/client-go v0.19.5
	k8s.io/cloud-provider v0.19.5
	k8s.io/code-generator v0.19.5
	k8s.io/kops v1.19.0
	k8s.io/kubelet v0.19.5
	k8s.io/kubernetes v1.19.5
	k8s.io/legacy-cloud-providers v0.19.5
	sigs.k8s.io/aws-iam-authenticator v0.5.2
	sigs.k8s.io/mdtoc v1.0.1
	sigs.k8s.io/yaml v1.2.0
)

replace (
	github.com/aws/aws-sdk-go => github.com/weaveworks/aws-sdk-go v0.0.0-20210212091355-35b293563a18
	// Used to pin the k8s library versions regardless of what other dependencies enforce
	k8s.io/api => k8s.io/api v0.19.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.5
	k8s.io/apiserver => k8s.io/apiserver v0.19.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.5
	k8s.io/client-go => k8s.io/client-go v0.19.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.5
	k8s.io/code-generator => k8s.io/code-generator v0.19.5
	k8s.io/component-base => k8s.io/component-base v0.19.5
	k8s.io/cri-api => k8s.io/cri-api v0.19.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.5
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.5
	k8s.io/kubectl => k8s.io/kubectl v0.19.5
	k8s.io/kubelet => k8s.io/kubelet v0.19.5
	k8s.io/kubernetes => k8s.io/kubernetes v1.19.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.5
	k8s.io/metrics => k8s.io/metrics v0.19.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.5
)
