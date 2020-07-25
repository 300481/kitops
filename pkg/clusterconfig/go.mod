module github.com/300481/kitops/pkg/clusterconfig

go 1.14

replace github.com/300481/kitops/pkg/apiresource => ../apiresource

replace github.com/300481/kitops/pkg/sourcerepo => ../sourcerepo

require (
	github.com/300481/kitops/pkg/apiresource v0.0.0-20200725192506-e6155c3e640b
	github.com/300481/kitops/pkg/sourcerepo v0.0.0-20200725192506-e6155c3e640b
)
