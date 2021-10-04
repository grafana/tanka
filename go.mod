module github.com/grafana/tanka

go 1.16

require (
	github.com/Masterminds/semver v1.5.0
	github.com/Masterminds/sprig/v3 v3.2.2
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/go-clix/cli v0.2.0
	github.com/gobwas/glob v0.2.3
	github.com/google/go-cmp v0.5.6
	github.com/google/go-jsonnet v0.17.0
	github.com/google/uuid v1.1.2 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/karrick/godirwalk v1.16.1
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pkg/errors v0.9.1
	github.com/posener/complete v1.2.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/objx v0.3.0
	github.com/stretchr/testify v1.7.0
	github.com/thoas/go-funk v0.9.1
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.4.0
	// yaml.v3 should be bumped with care. The new versions change all list indents
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c
	k8s.io/apimachinery v0.20.0-beta.1
	k8s.io/klog/v2 v2.9.0 // indirect
	sigs.k8s.io/yaml v1.3.0
)
