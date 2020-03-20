module github.com/NKztq/spider

go 1.13

require (
	bou.ke/monkey v0.0.0-00010101000000-000000000000
	github.com/baidu/go-lib v0.0.0-20191217050907-c1bbbad6b030
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20180821023952-922f4815f713
	gopkg.in/gcfg.v1 v1.2.3
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

replace (
	bou.ke/monkey => github.com/bouk/monkey v1.0.3-0.20191209094521-b118a1738765
	golang.org/x/net v0.0.0-20180821023952-922f4815f713 => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
)
