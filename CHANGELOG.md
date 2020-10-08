# Terraform oraclerdbms provider changelog

## 0.5.3 (October 8, 2020)

* Bumping terraform-oracle-rdbms-helper to v0.5.0
* Dependency bumping terraform-plugin-sdk

## 0.5.2 (January 3, 2020)

* Bumping terraform-oracle-rdbms-helper to v0.4.2

## 0.5.1 (January 2, 2020)

* Bumping terraform-oracle-rdbms-helper to v0.4.1

## 0.5.0 (September 27, 2019)

* As of September 2019, Terraform provider developers importing the Go module github.com/hashicorp/terraform, known as Terraform Core, should switch to github.com/hashicorp/terraform-plugin-sdk, the Terraform Plugin SDK, instead.

## 0.4.0 (September 24, 2019)

* Adding support for terraform 0.12

## 0.3.2 (February, 14, 2019)

* Bumping terraform-oracle-rdbms-helper to v0.3.1
* Add resource oraclerdbms_block_change_tracking.
  * Does not support OMF.
* Switch to go mod

## 0.3.1 (December, 08, 2018)

* Adding database resource to handle force logging and flashback
* Adding autotask resource to enable/disable sql tuning advisor, auto optimizer stats collection, auto space advisor

## 0.3.0 (November 27, 2018)

* Bumping terraform-oracle-rdbms-helper to v0.2.8.1
* Adding support for user account_status , open, locked, expired
* Adding support for user resource quota
* Removing password attribute from resource user

## 0.2.9 (November 15, 2018)

* Bumping terraform-oracle-rdbms-helper to v0.2.7
* Adding support to delete resource operation.
* objects_sha256 & privs_sha256 from string to map to track sha256 per privilege

## 0.2.8 (November 13, 2018)

* Updating vendoring for terraform-oracle-rdbms-helper to v0.2.5
* Implement resource oraclerdbms_stats

## 0.2.7 (November 12, 2018)

* Changing name on struct that olding the OracleApi from providerConfiguration to oracleHelperType
* Generate a diff if privs_sha256 is not equal to objects_sha256

## 0.2.6 (November 9, 2018)

* using the GetHashSchemaPrivsToUser to get diff

## 0.2.5 (November 9, 2018)

* Updating vendoring for terraform-oracle-rdbms-helper to 0.2.3

## 0.2.4 (November 9, 2018)

* Adding some logging for debuging

## 0.2.3 (November 8, 2018)

* Updating vendoring for terraform-oracle-rdbms-helper to 0.2.2
* First alpha release of grant schema to user. Just support tables and now it revokes/grants everything the grant if there is something diff

## 0.2.2 (November 2, 2018)

* Updating vendoring for terraform-oracle-rdbms-helper to 0.2.1

## 0.2.1 (October 24, 2018)

NOTES:

* Implement Importer for all the resources

## 0.2.0 (October 22, 2018)

NOTES:

### resource_profile_limite

* Store in upercase in the state
* Changes attribute profile_id => profile
* Changes attribute limit => resource_name

## 0.1.1 (October 18, 2018)

NOTES:

uppdating go vendoring

## 0.1.0 (October 15, 2018)

NOTES:

* resource_grant_object_privileges
* resource_grant_role_privilege
* resource_grant_system_privilege
* resource_parameter
* resource_profile
* resource_profile_limit
* resource_role
* resource_user