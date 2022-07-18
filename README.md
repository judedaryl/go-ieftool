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

# Getting started

> Credentials are set using environment variables

```sh
export B2C_TENANT_ID=mytenant.onmicrosoft.com
export B2C_CLIENT_ID=00000000-0000-0000-0000-000000000000
export B2C_CLIENT_SECRET=some_secret

ieftool deploy {POLICY_PATH}
```

