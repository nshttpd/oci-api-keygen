### oci-api-keygen

In order to use the [Oracle Compute Infrastructure](https://cloud.oracle.com/en_US/infrastructure/compute) API a
user must generate a key to upload into the web console. This key is used to authenticate to the API for a specific
[tenancy](https://docs.us-phoenix-1.oraclecloud.com/Content/GSG/Concepts/settinguptenancy.htm) for tools such as 
[Terraform](https://www.terraform.io/) or other third-party or home build tools.

The [instructions for creating](https://docs.us-phoenix-1.oraclecloud.com/Content/API/Concepts/apisigningkey.htm) 
these keys is a multi step process involving openssl and command line while keeping track of the files that are 
generated. This tool handles all of that for you. It'll generate the private and public key along with keeping 
track of the tenancies they are created for and the fingerprints of the keys.

#### Installation

The usual.

`go get -u github.com/nshttpd/oci-api-keygen`

At some point in the future binaries for different platorms will be provided.

#### Usage

The default location for the public and private keys along with a configuration file will be `~/.oci/` where most 
other things associated with OCI are stored. This can be overridden with a `--config` parameter and the keys 
will be stored in the same directory as where the config file is set to.

The basics are :

**Create**

create a set of keys for a tenancy

> oci-api-keygen create [tenancy]

**List**

List all of the tenancies that have had keys generated for them

> oci-api-keygen list

Show fingerprint for a tenancy

> oci-api-keygen list [tenancy]

