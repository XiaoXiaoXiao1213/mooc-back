module management

go 1.12

//无法访问的原因，替换golang.org源为github.com源
replace (
	cloud.google.com/go => github.com/googleapis/google-cloud-go v0.37.2
	golang.org/x/build => github.com/golang/build v0.0.0-20190327004547-5a2224f3eb52
	golang.org/x/crypto => github.com/golang/crypto v0.0.0-20190325154230-a5d413f7728c
	golang.org/x/exp => github.com/golang/exp v0.0.0-20190321205749-f0864edee7f3
	golang.org/x/image => github.com/golang/image v0.0.0-20190321063152-3fc05d484e9f
	golang.org/x/lint => github.com/golang/lint v0.0.0-20190313153728-d0100b6bd8b3
	golang.org/x/mobile => github.com/golang/mobile v0.0.0-20190319155245-9487ef54b94a
	golang.org/x/net => github.com/golang/net v0.0.0-20190327025741-74e053c68e29
	golang.org/x/oauth2 => github.com/golang/oauth2 v0.0.0-20190319182350-c85d3e98c914
	golang.org/x/perf => github.com/golang/perf v0.0.0-20190312170614-0655857e383f
	golang.org/x/sync => github.com/golang/sync v0.0.0-20190227155943-e225da77a7e6
	golang.org/x/sys => github.com/golang/sys v0.0.0-20190322080309-f49334f85ddc
	golang.org/x/text => github.com/golang/text v0.3.0
	golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
	golang.org/x/tools => github.com/golang/tools v0.0.0-20190327011446-79af862e6737
	google.golang.org/api => github.com/googleapis/google-api-go-client v0.3.0
	google.golang.org/appengine => github.com/golang/appengine v1.5.0
	google.golang.org/genproto => github.com/google/go-genproto v0.0.0-20190321212433-e79c0c59cdb5
	google.golang.org/grpc => github.com/grpc/grpc-go v1.19.1

)

require (
	github.com/Joker/jade v1.0.0 // indirect
	github.com/develop1024/jwt-go v3.2.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/flosch/pongo2 v0.0.0-20181225140029-79872a7b2769 // indirect
	github.com/gavv/monotime v0.0.0-20171021193802-6f8212e8d10d // indirect
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/go-playground/locales v0.12.1
	github.com/go-playground/universal-translator v0.16.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/iris-contrib/blackfriday v2.0.0+incompatible // indirect
	github.com/iris-contrib/formBinder v0.0.0-20190104093907-fbd5963f41e1 // indirect
	github.com/iris-contrib/httpexpect v0.0.0-20180314041918-ebe99fcebbce // indirect
	github.com/iris-contrib/middleware/jwt v0.0.0-20210110101738-6d0a4d799b5d // indirect
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/k0kubun/colorstring v0.0.0-20150214042306-9440f1994b88 // indirect
	github.com/kataras/iris v11.1.1+incompatible
	github.com/leodido/go-urn v1.1.0 // indirect
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/mattn/go-colorable v0.1.1
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/moul/http2curl v1.0.0 // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/ryanuber/columnize v0.0.0-20160712163229-9b3edd62028f // indirect
	github.com/sirupsen/logrus v1.4.1
	github.com/smartystreets/goconvey v0.0.0-20190306220146-200a235640ff // indirect
	github.com/tebeka/strftime v0.0.0-20140926081919-3f9c7761e312 // indirect
	github.com/tietang/dbx v0.0.0-20190314092248-c144fca0b077
	github.com/tietang/go-utils v0.1.2
	github.com/tietang/props v2.3.0+incompatible
	github.com/valyala/fasttemplate v1.0.1 // indirect
	github.com/x-cray/logrus-prefixed-formatter v0.5.2
	github.com/yudai/pp v2.0.1+incompatible // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v9 v9.27.0
)
