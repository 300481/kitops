module github.com/300481/kitops/cmd/kitops

go 1.14

require (
	github.com/300481/kitops/pkg/kitops v0.0.0-20200723202845-cae2a42451aa
	github.com/urfave/cli/v2 v2.2.0
)

replace github.com/300481/kitops/pkg/kitops => ../../pkg/kitops
