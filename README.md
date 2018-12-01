Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Terraform team at [HashiCorp](https://www.hashicorp.com/).

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) 0.10.x
- [Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)go-oci8[()]
- [go-oci8](https://github.com/mattn/go-oci8)
- [terraform-oracle-rdbms-helper](https://github.com/dbgeek/terraform-oracle-rdbms-helper)

Usage
---------------------

```terraform
# For example, restrict oraclerdbms version in 0.1.x
provider "oraclerdbms" {
  version = "~> 0.1"
}

resource "oraclerdbms_profile" "test" {
    profile = "TEST01"
}

resource "oraclerdbms_profile_limit" "test1" {
  resource_name = "IDLE_TIME"
  value         = "33"
  profile       = "${oraclerdbms_profile.test.id}"

}

resource "oraclerdbms_user" "testuser" {
  username            = "TESTUSER"
  default_tablespace  = "USERS"
  profile             = "${oraclerdbms_profile.test.id}
}

resource "oraclerdbms_grant_system_privilege" "grantsysprivs" {
  grantee   = "${oraclerdbms_user.testuser.id}"
  privilege = "CREATE SESSION"
}

resource "oraclerdbms_role" "roletest" {
  role = "TESTROLE"
}

resource "oraclerdbms_grant_role_privilege" "grantroleprivs" {
  grantee = "${oraclerdbms_user.testuser.id}"
  role    = "${oraclerdbms_role.roletest.id}"
}
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-oraclerdbms`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-oraclerdbms
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-oraclerdbms
$ make build
```

Using the provider
----------------------
## Fill in for each provider

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-oraclerdbms
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```
