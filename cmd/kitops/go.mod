module github.com/300481/kitops/cmd/kitops

go 1.14

require (
	github.com/300481/kitops/pkg/kitops v0.0.0-20200803105738-85e607a14515
	github.com/urfave/cli/v2 v2.2.0
)

replace github.com/300481/kitops/pkg/kitops => ../../pkg/kitops
