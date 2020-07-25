module github.com/300481/kitops/pkg/kitops

go 1.14

require (
	github.com/300481/kitops/pkg/queue v0.0.0-20200725203232-1022066be267
	github.com/300481/kitops/pkg/sourcerepo v0.0.0-20200725203232-1022066be267
	github.com/gorilla/mux v1.7.4
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)

replace (
	github.com/300481/kitops/pkg/queue => ../queue
	github.com/300481/kitops/pkg/sourcerepo => ../sourcerepo
)
