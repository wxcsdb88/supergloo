# CLI for Supergloo
## Installation
```
make install-cli    # run from the project root directory
```
## Commands

### Help
Lists the available commands.
#### Usage
```
supergloo help
```

### Get
Displays one or many supergloo resources in table format.
#### Usage
```bash
supergloo get RESOURCE_TYPE [RESOURCE_NAME] [-o|--output OUTPUT_TYPE]
```
#### Options
| name | required | default | description |
| ---- |   ----   |   ----  |    ----     |
| output | N | | Output format. Currently only the `wide` option, which causes additional columns to be displayed, is supported. |
##### Example
```bash
supergloo get meshes my-mesh -o wide
```

### Create 
Create a resource from stdin.
#### Routing rule
Create a routing rule. A routing rule controls how requests are routed within the target service mesh.
##### Usage
```bash
supergloo create routingrule --mesh TARGET_MESH
    [--namespace NAMESPACE]
    [--sources NAMESPACE:NAME[,NAMESPACE:NAME...]
    [--destinations NAMESPACE:NAME[,NAMESPACE:NAME...]
    [--matchers prefix=PREFIX|methods=METHODS[,prefix=PREFIX|methods=METHODS...]
    [--override true|false]
```
##### Options
| name | required | default | description |
| ---- |   ----   |   ----  |    ----     |
| mesh | Y | | The mesh that will be the target for this rule. |
| namespace | N | default | The namespace this routing rule will be created in. Defaults to "default". |
| sources | N | | Source upstreams for this rule. The value for this option is a comma-separated list of upstreams. Each entry consists of an upstream namespace and and upstream name, separated by a colon. |
| destinations | N | | Destination upstreams for this rule. Same format as `sources`. |
| matchers | N | | Matchers determine which source requests the routing rule get applied to. |
| override | N | false | If false, the operation will fail if a routing rule with the given name exists in the given namespace. |
##### Example
```bash
supergloo create routingrule my-rule --mesh my-mesh 
    --namespace my-ns-1
    --sources my-ns-2:my-upstream-2,my-ns-3:my-upstream-3
    --destinations some-other-ns:some-other-upstream
    --matchers "prefix=/some/path,methods=get|post"
    --override true
    
```