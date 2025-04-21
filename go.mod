module flkcli

go 1.22.1

require (
	github.com/spf13/cobra v1.8.0
	gopkg.in/masci/flickr.v3 v3.0.0-20230428071620-b971d524ac6f 
)

replace (
	gopkg.in/masci/flickr.v3 v3.0.0-20230428071620-b971d524ac6f  => ../flickr
)

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v2 v2.4.0
)
