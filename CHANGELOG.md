## 2.0.0 (March 24, 2021)

FEATURES:
* appstream/resource_fleet.go - aws api compatibility fixes
* appstream/resource_stack.go - aws api compatibility fixes
* appstream/resource_stack_attachment.go - new resource to achieve terraform naming standards

ENHANCEMENTS:
* Upgraded sdk modules
* Updated examples

BUGFIXES:
* Multiple api inconsistencies
* Multiple error handling issues

## 1.0.8 (June 15, 2020)

FEATURES:
* appstream/resource_stack.go - user_settings patch from dhruv2511

ENHANCEMENTS:
* Upgraded modules
* Updated examples

BUGFIXES:


## 1.0.7 (June 9, 2020)

FEATURES:
* Added support for role ARN

ENHANCEMENTS:

BUGFIXES:

Patch by: Konstantin Odnoralov <kodnoral@pmintl.net>

## 1.0.6 (May 27, 2020)

FEATURES:
* Added Ability to Remote Image
* Changes to iamge_arn forces new stack

ENHANCEMENTS:
* updated tf lib to 0.12.25

BUGFIXES:
* image_name changed to image_arn

Patch by: Konstantin Odnoralov <kodnoral@pmintl.net>


## 1.0.5 (May 03, 2020)

FEATURES:
* Assume Role authentication

ENHANCEMENTS:
* authentication: Adopted AWS authentication from terraform-provider-aws
* structure: changed structure and build setup of provider

BUGFIXES:


