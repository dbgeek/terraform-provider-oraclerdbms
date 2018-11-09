# Terraform oraclerdbms provider changelog

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