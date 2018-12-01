provider "oraclerdbms" {}

resource "oraclerdbms_profile" "profile" {
  profile = "DEMO01"
}

resource "oraclerdbms_profile_limit" "profile_idle_time" {
  resource_name = "IDLE_TIME"
  value         = "33"
  profile       = "${oraclerdbms_profile.profile.id}"
}

resource "oraclerdbms_user" "example_user" {
  username           = "EXAMPLE_USER"
  default_tablespace = "USERS"
  profile            = "${oraclerdbms_profile.profile.id}"
}

resource "oraclerdbms_grant_system_privilege" "syspriv_create_session" {
  grantee   = "${oraclerdbms_user.example_user.id}"
  privilege = "CREATE SESSION"
}

resource "oraclerdbms_role" "example_role" {
  role = "EXAMPLE_ROLE"
}

resource "oraclerdbms_grant_role_privilege" "grantroleprivs" {
  grantee = "${oraclerdbms_user.example_user.id}"
  role    = "${oraclerdbms_role.example_role.id}"
}

resource "oraclerdbms_stats" "granularity_auto" {
  preference_name  = "GRANULARITY"
  preference_value = "AUTO"
}

resource "oraclerdbms_parameter" "testatistics_levelst" {
  name           = "statistics_level"
  value          = "ALL"
  update_comment = "setting statistics level to all"
  scope          = "BOTH"
}
