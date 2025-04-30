# Buf 
* `buf generate` will geneate the pb files
* buf.yaml is the mail file for ocnfiguration
* You can define modules in `buf.yaml` and split your protos into modules.
  * Here it's storage and proto module each having it's configuration.
  * Alternatively you can use `buf.work.yaml` to define these workspaces and have `buf.yaml` inside each workspace (module specific)
* `out=source_relative` will put the generated files in the same structure of protos and doesn't follow the `go_package` option
  * If you don't specify it will follow `go_package` strucutre
