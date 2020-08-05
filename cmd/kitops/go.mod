module github.com/300481/kitops/cmd/kitops

go 1.14

require (
	github.com/300481/kitops/pkg/kitops v0.0.0-20200805123234-032e5f213d70
	github.com/urfave/cli/v2 v2.2.0
)

replace github.com/300481/kitops/pkg/kitops => ../../pkg/kitops
