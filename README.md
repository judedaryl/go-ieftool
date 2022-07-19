# Identity Experience Framework tool

This is a port of the nodejs version https://github.com/judedaryl/ieftool which removes any dependency on nodejs and external libraries and reduces the file size to 9mb.


This tool makes it easier for B2C policies to be uploaded in-order based on the inheritance of a policy. Uploads are also faster because policies are uploaded by batch depending on its position on the inheritance tree.


```pre
src/
├─ social/
│  ├─ base.xml (1A_SBASE)
│  ├─ signupsignin.xml (1A_SSS)
├─ local/
│  ├─ base.xml (1A_LBASE)
│  ├─ signupsignin.xml (1A_LSS)
│  ├─ passwordreset.xml (1A_LPR)
├─ base.xml (1A_BASE)
├─ extension.xml (1A_EXT)

```

The example folder structure above has the following inheritance tree.

```pre
                1A_BASE
                    |
                 1A_EXT
                /      \
          1A_LBASE    1A_SBASE
           /    \        \      
       1A_LSS  1A_LPR    1A_SSS
```

These policies are then batched by their hierarchy in the tree, as well as their parent policy. The order of upload would then be.

1. 1A_Base
2. 1A_EXT
3. 1A_LBASE, 1A_SBASE
4. 1A_LSS, 1A_LPR
5. 1A_LSSS



<br/>
<br/>

# Installation

Install via curl

```
curl https://raw.githubusercontent.com/judedaryl/go-ieftool/main/install.sh | bash
```

# Commands

## Build

Compiles and injects variable values into source IEF policies (.xml). The variables are extracted from a configuration file that you can provide using ``--config`` or ``-c`` (defaults to ``ieftool.config``)

### Usage:
ieftool build [path to source code] [path to target directory] [flags]

### Flags:
|flag|alias|type|description|
|-|-|-|-|
|--config|-c|string|Path to the ieftool configuration file (yaml) (default "ieftool.config")|
|--help|-h|-|help for build|

### Example:

``ieftool.config``
```yaml
tenantId: mytenant.onmicrosoft.com
deploymentMode: Development
```

``src/BasePolicy.xml``
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<TrustFrameworkPolicy 
    ...
    TenantId="{{ tenantId }}"  
    DeploymentMode="{{ deploymentMode }}">
  ...
</xml>
```
Run the build command

```sh
# ieftool build [source dir] [target dir] -c [config path]
ieftool build src output -c ieftool.config
```

The policies are then compiled into

``output/BasePolicy.xml``
```xml
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<TrustFrameworkPolicy 
    ...
    TenantId="mytenant.onmicrosoft.com"  
    DeploymentMode="Development">
  ...
</xml>
```

## Deploy

Deploys your policies into Identity Experience Framework.

### Usage:
ieftool deploy [path to policies] [flags]

### Flags:
|flag|alias|type|description|
|-|-|-|-|
|--help|-h|-|help for build|


> Credentials are set using environment variables

```sh
export B2C_TENANT_ID=mytenant.onmicrosoft.com
export B2C_CLIENT_ID=00000000-0000-0000-0000-000000000000
export B2C_CLIENT_SECRET=some_secret

# ieftool deploy [path to policies]
ieftool deploy {POLICY_PATH}
```

