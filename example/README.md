# README

## Pull Docker image

```shell
docker pull ba78/terraform-provider-oraclerdbms:0.3.0
```

## Setup tf alias

```shell
alias tf="docker run \
        -u $(id -u $(whoami)):$(id -g $(whoami)) \
        -v $(pwd):/opt/tf \
        -w /opt/tf \
        --rm \
        -e ORACLE_DBHOST \
        -e ORACLE_DBPORT \
        -e ORACLE_SERVICE \
        -e ORACLE_USERNAME \
        -e ORACLE_PASSWORD \
        -e TF_LOG \
        -e HOME=/home/tf \
        -it \
        ba78/terraform-provider-oraclerdbms:0.3.0"
```

## init

tf init

## plan

```shell
tf plan
```

### output

```
C02WH1TKHTD5:example bjorn.ahl$ tf plan
Refreshing Terraform state in-memory prior to plan...
The refreshed state will be used to calculate this plan, but will not be
persisted to local or remote state storage.


------------------------------------------------------------------------

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + oraclerdbms_grant_role_privilege.grantroleprivs
      id:                   <computed>
      grantee:              "${ORACLERDBMS_USER.EXAMPLE_USER.ID}"
      role:                 "${ORACLERDBMS_ROLE.EXAMPLE_ROLE.ID}"

  + oraclerdbms_grant_system_privilege.syspriv_create_session
      id:                   <computed>
      grantee:              "${ORACLERDBMS_USER.EXAMPLE_USER.ID}"
      privilege:            "CREATE SESSION"

  + oraclerdbms_profile.profile
      id:                   <computed>
      profile:              "DEMO01"

  + oraclerdbms_profile_limit.profile_idle_time
      id:                   <computed>
      profile:              "${oraclerdbms_profile.profile.id}"
      resource_name:        "IDLE_TIME"
      value:                "33"

  + oraclerdbms_role.example_role
      id:                   <computed>
      role:                 "EXAMPLE_ROLE"

  + oraclerdbms_stats.granularity_auto
      id:                   <computed>
      preference_name:      "GRANULARITY"
      preference_value:     "AUTO"

  + oraclerdbms_user.example_user
      id:                   <computed>
      account_status:       "OPEN"
      default_tablespace:   "USERS"
      profile:              "${ORACLERDBMS_PROFILE.PROFILE.ID}"
      temporary_tablespace: <computed>
      username:             "EXAMPLE_USER"


Plan: 7 to add, 0 to change, 0 to destroy.

------------------------------------------------------------------------

Note: You didn't specify an "-out" parameter to save this plan, so Terraform
can't guarantee that exactly these actions will be performed if
"terraform apply" is subsequently run.
```

## Apply

```shell
tf apply
```

### output

```
C02WH1TKHTD5:example bjorn.ahl$ tf apply

An execution plan has been generated and is shown below.
Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  + oraclerdbms_grant_role_privilege.grantroleprivs
      id:                   <computed>
      grantee:              "${ORACLERDBMS_USER.EXAMPLE_USER.ID}"
      role:                 "${ORACLERDBMS_ROLE.EXAMPLE_ROLE.ID}"

  + oraclerdbms_grant_system_privilege.syspriv_create_session
      id:                   <computed>
      grantee:              "${ORACLERDBMS_USER.EXAMPLE_USER.ID}"
      privilege:            "CREATE SESSION"

  + oraclerdbms_profile.profile
      id:                   <computed>
      profile:              "DEMO01"

  + oraclerdbms_profile_limit.profile_idle_time
      id:                   <computed>
      profile:              "${oraclerdbms_profile.profile.id}"
      resource_name:        "IDLE_TIME"
      value:                "33"

  + oraclerdbms_role.example_role
      id:                   <computed>
      role:                 "EXAMPLE_ROLE"

  + oraclerdbms_stats.granularity_auto
      id:                   <computed>
      preference_name:      "GRANULARITY"
      preference_value:     "AUTO"

  + oraclerdbms_user.example_user
      id:                   <computed>
      account_status:       "OPEN"
      default_tablespace:   "USERS"
      profile:              "${ORACLERDBMS_PROFILE.PROFILE.ID}"
      temporary_tablespace: <computed>
      username:             "EXAMPLE_USER"


Plan: 7 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

oraclerdbms_stats.granularity_auto: Creating...
  preference_name:  "" => "GRANULARITY"
  preference_value: "" => "AUTO"
oraclerdbms_role.example_role: Creating...
  role: "" => "EXAMPLE_ROLE"
oraclerdbms_profile.profile: Creating...
  profile: "" => "DEMO01"
oraclerdbms_stats.granularity_auto: Creation complete after 0s (ID: STATS-GLOBAL-GRANULARITY)
oraclerdbms_profile.profile: Creation complete after 0s (ID: DEMO01)
oraclerdbms_profile_limit.profile_idle_time: Creating...
  profile:       "" => "DEMO01"
  resource_name: "" => "IDLE_TIME"
  value:         "" => "33"
oraclerdbms_user.example_user: Creating...
  account_status:       "" => "OPEN"
  default_tablespace:   "" => "USERS"
  profile:              "" => "DEMO01"
  temporary_tablespace: "" => "<computed>"
  username:             "" => "EXAMPLE_USER"
oraclerdbms_profile_limit.profile_idle_time: Creation complete after 0s (ID: DEMO01-IDLE_TIME)
oraclerdbms_role.example_role: Creation complete after 0s (ID: EXAMPLE_ROLE)
oraclerdbms_user.example_user: Creation complete after 0s (ID: EXAMPLE_USER)
oraclerdbms_grant_system_privilege.syspriv_create_session: Creating...
  grantee:   "" => "EXAMPLE_USER"
  privilege: "" => "CREATE SESSION"
oraclerdbms_grant_role_privilege.grantroleprivs: Creating...
  grantee: "" => "EXAMPLE_USER"
  role:    "" => "EXAMPLE_ROLE"
oraclerdbms_grant_system_privilege.syspriv_create_session: Creation complete after 0s (ID: EXAMPLE_USER-CREATE SESSION)
oraclerdbms_grant_role_privilege.grantroleprivs: Creation complete after 0s (ID: EXAMPLE_USER-EXAMPLE_ROLE)

Apply complete! Resources: 7 added, 0 changed, 0 destroyed.
```
