# Identity Experience Framework tool

This is a port of the nodejs version https://github.com/judedaryl/ieftool which removes any dependency on nodejs and external libraries and reduces the file size to 9mb.

This tool enhances the development experience with B2C policies, policies can now be ``multi-environment`` by introducing different variable values depending on the environment and helps you upload your policies to Azure B2C seamlessly.

### Variables

B2C policies are built on xml and has no support for variables, ieftool introduces a build command that lets you inject variables to your policies either through a configuration file or environment variables. See the [build command](https://github.com/judedaryl/go-ieftool/blob/main/README.md#build) below for more information.

### Uploads
Policies are uploaded in-order based on the inheritance of a policy. Uploads are also faster because policies are uploaded by batch depending on its position on the inheritance tree.


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

# Commands

## Usage:
```bash
ieftool
Tooling for Azure B2C Identity Experience Framework

Usage:
  ieftool [command]

Available Commands:
  build       Build
  completion  Generate completion script
  deploy      Deploy b2c policies.
  help        Help about any command
  list        List remote b2c policies.
  remove      Delete remote b2c policies.

Flags:
  -h, --help   help for ieftool
```

### Example config

```yaml
# config.yaml
- name: test
  tenant: test.onmicrosoft.com
  tenantId: aaaaaaaa-cccc-dddd-eeee-ffffffffffff
  clientId: aaaaaaaa-cccc-dddd-eeee-ffffffffffff
  settings:
    IdentityExperienceFrameworkAppId: aaaaaaaa-cccc-dddd-eeee-ffffffffffff
    ProxyIdentityExperienceFrameworkAppId: aaaaaaaa-cccc-dddd-eeee-ffffffffffff
    AADCommonClientID: aaaaaaaa-cccc-dddd-eeee-ffffffffffff
    AADCommonObjectID: aaaaaaaa-cccc-dddd-eeee-ffffffffffff
```

> Secrets for remote operations need to be set using environment variables `B2C_CLIENT_SECRET_<environent>`
```bash
export B2C_CLIENT_SECRET_TEST=mysecret
```
> Check for exact variable name
```bash
ieftool list -c test/fixtures/config.yaml -e test
2024/01/30 09:57:04 Failed to list policies could not create client credentials. Did you send the env var B2C_CLIENT_SECRET_TEST?: secret can't be empty string
```

The required variables in the above example would be 
- B2C_CLIENT_SECRET_TEST





## Build
```bash
ieftool build -h 
Build source policies and replacing template variables for given environments.

Usage:
  ieftool build [flags]

Flags:
  -c, --config string        Path to the ieftool configuration file (default "./config.yaml")
  -d, --destination string   Destination directory (default "./build")
  -e, --environment string   Environment to deploy (default: all environments)
  -h, --help                 help for build
  -s, --source string        Source directory (default "./src")
```

## List
```bash
ieftool list -h 
List remote b2c policies from B2C identity experience framework.

Usage:
  ieftool list [path to policies] [flags]

Flags:
  -c, --config string        Path to the ieftool configuration file (default "./config.yaml")
  -e, --environment string   Environment to deploy (default: all environments)
  -h, --help                 help for list
```

> Secret needs to be set using environment variable `B2C_CLIENT_SECRET_<environent>`

## Remove
```bash
ieftool remove -h 
Delete remote b2c policies from B2C identity experience framework.

Usage:
  ieftool remove [flags]

Flags:
  -c, --config string        Path to the ieftool configuration file (default "./config.yaml")
  -e, --environment string   Environment to deploy (default: all environments)
  -h, --help                 help for remove
```

> Secret needs to be set using environment variable `B2C_CLIENT_SECRET_<environent>`


## Deploy
```bash
ieftool deploy -h 
Deploy b2c policies to B2C identity experience framework.

Usage:
  ieftool deploy [path to policies] [flags]

Flags:
  -b, --build-dir string     Build directory (default "./build")
  -c, --config string        Path to the ieftool configuration file (default "./config.yaml")
  -e, --environment string   Environment to deploy (default: all environments)
  -h, --help                 help for deploy
```

> Secret needs to be set using environment variable `B2C_CLIENT_SECRET_<environent>`


