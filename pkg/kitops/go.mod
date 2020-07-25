module github.com/300481/kitops/pkg/kitops

go 1.14

require (
	github.com/300481/kitops/pkg/clusterconfig v0.0.0-20200725192506-e6155c3e640b
	github.com/300481/kitops/pkg/queue v0.0.0-20200725192506-e6155c3e640b
	github.com/300481/kitops/pkg/sourcerepo v0.0.0-20200725192506-e6155c3e640b
	github.com/gorilla/mux v1.7.4
)

replace (
	github.com/300481/kitops/pkg/clusterconfig => ../clusterconfig
	github.com/300481/kitops/pkg/queue => ../queue
	github.com/300481/kitops/pkg/sourcerepo => ../sourcerepo
)
