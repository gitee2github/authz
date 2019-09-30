module isula.org/authz

replace (
	github.com/docker/docker => github.com/docker/engine v0.0.0-20190717160951-456712c5b8d9
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190701094942-4def268fd1a4
	golang.org/x/net => github.com/golang/net v0.0.0-20190813141303-74dc4d7220e7
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190813064441-fde4db37ae7a
	golang.org/x/text => github.com/golang/text v0.3.2
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190816200558-6889da9d5479
	golang.org/x/xerrors => github.com/golang/xerrors v0.0.0-20190717185122-a985d3407aa7
	gopkg.in/yaml.v2 => github.com/go-yaml/yaml v2.1.0+incompatible
)

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/containerd/continuity v0.0.0-20190815185530-f2a389ac0a02 // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/gorilla/mux v1.7.3
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli v1.21.0
	golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2 // indirect
)
