module ehang.io/nps

go 1.15

require (
	ehang.io/nps-mux v0.0.0-20210407130203-4afa0c10c992
	fyne.io/fyne/v2 v2.0.2
	github.com/astaxie/beego v1.12.0
	github.com/bradfitz/iter v0.0.0-20191230175014-e8f45d346db8 // indirect
	github.com/c4milo/unpackit v0.0.0-20170704181138-4ed373e9ef1c
	github.com/caddyserver/caddy/v2 v2.4.6
	github.com/ccding/go-stun v0.0.0-20180726100737-be486d185f3d
	github.com/dsnet/compress v0.0.1 // indirect
	github.com/golang/snappy v0.0.3
	github.com/hooklift/assert v0.0.0-20170704181755-9d1defd6d214 // indirect
	github.com/kardianos/service v1.2.0
	github.com/klauspost/pgzip v1.2.1 // indirect
	github.com/klauspost/reedsolomon v1.9.12 // indirect
	github.com/panjf2000/ants/v2 v2.4.2
	github.com/pkg/errors v0.9.1
	github.com/shiena/ansicolor v0.0.0-20151119151921-a422bbe96644 // indirect
	github.com/shirou/gopsutil/v3 v3.21.3
	github.com/templexxx/cpufeat v0.0.0-20180724012125-cef66df7f161 // indirect
	github.com/templexxx/xor v0.0.0-20191217153810-f85b25db303b // indirect
	github.com/tidwall/gjson v1.12.1
	github.com/tidwall/sjson v1.2.4
	github.com/tjfoc/gmsm v1.4.0 // indirect
	github.com/xtaci/kcp-go v5.4.20+incompatible
	github.com/xtaci/lossyconn v0.0.0-20190602105132-8df528c0c9ae // indirect
	golang.org/x/net v0.0.0-20210913180222-943fd674d43e
)

replace github.com/astaxie/beego => github.com/exfly/beego v1.12.0-export-init

replace github.com/caddyserver/caddy/v2 => github.com/dounx/caddy/v2 v2.4.7
