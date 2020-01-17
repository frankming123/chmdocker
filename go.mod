module chmdocker

go 1.13

replace (
	github.com/Sirupsen/logrus v1.0.5 => github.com/sirupsen/logrus v1.0.5
	github.com/Sirupsen/logrus v1.3.0 => github.com/Sirupsen/logrus v1.0.6
	github.com/Sirupsen/logrus v1.4.0 => github.com/sirupsen/logrus v1.0.6
)

require (
	github.com/Sirupsen/logrus v1.4.0
	github.com/onsi/ginkgo v1.11.0 // indirect
	github.com/onsi/gomega v1.8.1 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect
	github.com/stretchr/testify v1.4.0 // indirect
	github.com/urfave/cli v1.22.2
	github.com/xianlubird/mydocker v0.0.0-20180315123543-9ea8dbc2b308
	golang.org/x/crypto v0.0.0-20191227163750-53104e6ec876 // indirect
	golang.org/x/sys v0.0.0-20191228213918-04cbcbbfeed8
	gopkg.in/airbrake/gobrake.v2 v2.0.9 // indirect
	gopkg.in/gemnasium/logrus-airbrake-hook.v2 v2.1.2 // indirect
)
