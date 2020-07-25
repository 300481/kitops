module github.com/300481/kitops/pkg/clusterconfig

go 1.14

replace github.com/300481/kitops/pkg/sourcerepo => ../sourcerepo

require (
	github.com/300481/kitops/pkg/sourcerepo v0.0.0-20200725203057-20c2740f1cd6
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776
)
