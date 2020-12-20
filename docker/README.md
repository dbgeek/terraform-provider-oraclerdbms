# terraform oraclerdbms in Docker

## build

```sh
./build
```

## example run

Example how to run oracle_rdbms using docker. This is an raw example. To make more usable we need to build some bash script wrapper around the the docker run commands.
We need to have ORACLE_SERVICE set to the tnsnames.ora entry we want to use and we need also set ORACLE_USERNAME & ORACLE_PASSWORD for which db user and password we should use.
`TNS_NAMES` need also be exported to the the path where your tnsnames.ora is.

### terraform init

```text
docker run \
        -v $(pwd):/terraform \
        -v "$TNS_ADMIN/tnsnames.ora":/home/terraform/tnsnames.ora \
        -e ORACLE_SERVICE \
        -e ORACLE_USERNAME \
        -e ORACLE_PASSWORD \
         tf-oraclerdbms:0.5.3.1 init

Initializing the backend...

Initializing provider plugins...
- Finding local/tf/oraclerdbms versions matching "~> 0.1"...
- Installing local/tf/oraclerdbms v0.5.3...
- Installed local/tf/oraclerdbms v0.5.3 (unauthenticated)

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

### terraform plan

```text
docker run \
        -v $(pwd):/terraform \
        -v $TNS_ADMIN/tnsnames.ora:/home/terraform/tnsnames.ora \
        -e ORACLE_SERVICE \
        -e ORACLE_USERNAME \
        -e ORACLE_PASSWORD \
        tf-oraclerdbms:0.5.3.1 plan -out=terraform.tfplan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # oraclerdbms_profile.app_profile will be created
  + resource "oraclerdbms_profile" "app_profile" {
      + id      = (known after apply)
      + profile = "APP_PROFILE"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

This plan was saved to: terraform.tfplan

To perform exactly these actions, run the following command to apply:
    terraform apply "terraform.tfplan"
```

### terraform apply

```text
docker run \
        -v $(pwd):/terraform \
        -v $TNS_ADMIN/tnsnames.ora:/home/terraform/tnsnames.ora \
        -e ORACLE_SERVICE \
        -e ORACLE_USERNAME \
        -e ORACLE_PASSWORD \
        tf-oraclerdbms:0.5.3.1 apply "terraform.tfplan"
oraclerdbms_profile.app_profile: Creating...
oraclerdbms_profile.app_profile: Creation complete after 0s [id=APP_PROFILE]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

The state of your infrastructure has been saved to the path
below. This state is required to modify and destroy your
infrastructure, so keep it safe. To inspect the complete state
use the `terraform show` command.

State path: terraform.tfstate
```
