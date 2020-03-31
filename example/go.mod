module gglmm-account-example

go 1.13

replace github.com/weihongguo/gglmm => ../../gglmm

replace github.com/weihongguo/gglmm-redis => ../../gglmm-redis

replace github.com/weihongguo/gglmm-account => ../

require (
	github.com/jinzhu/gorm v1.9.12
	github.com/weihongguo/gglmm v0.0.0-20200226150144-384f169aa64a
	github.com/weihongguo/gglmm-account v0.0.0-20200317144519-84e6a1300420
	github.com/weihongguo/gglmm-redis v0.0.0-00010101000000-000000000000
)
