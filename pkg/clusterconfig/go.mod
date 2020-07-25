module github.com/300481/kitops/pkg/clusterconfig

go 1.14

replace github.com/300481/kitops/pkg/apiresource => ../apiresource

replace github.com/300481/kitops/pkg/sourcerepo => ../sourcerepo

require (
	github.com/300481/kitops/pkg/apiresource v0.0.0-20200725194204-1d3a9134b56b
	github.com/300481/kitops/pkg/sourcerepo v0.0.0-20200725194204-1d3a9134b56b
)
