module github.com/300481/kitops/pkg/clusterconfig

go 1.14

replace github.com/300481/kitops/pkg/apiresource => ../apiresource

replace github.com/300481/kitops/pkg/sourcerepo => ../sourcerepo

require (
	github.com/300481/kitops/pkg/apiresource v0.0.0-20200723202845-cae2a42451aa
	github.com/300481/kitops/pkg/sourcerepo v0.0.0-20200722201655-f3681c684206
)
