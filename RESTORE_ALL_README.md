# StorCycle® Client and test utility

## Overview
This utility provides a CLI list/restore of all StorCycle jobs associated with a project.

## list_all_jobs 
- List all jobs for a project (--project_name parameter)

## restore_objects 
- Get all jobs for a project (--project_name parameter)
- Create one Restore Project per migrate job
- Run job now

## Operation
Unpack the executable onto a machine with network access to the StorCycle server.

Display available commands:
```shell
$ ssc-cli --command list_commands
```

Display available parameters:
```shell
$ ssc-cli --help
```
###List all jobs example:
```shell
$ ssc-cli --url https://localhost/openapi --name Administartor --password spectra --ignore_cert --command list_all_jobs --project_name TestArchive --verbose  
```
###Restore example:
```shell
$ ssc-cli --url https://localhost/openapi --name Administartor --password spectra  --ignore_cert --command restore_all_jobs --project_name TestArchive --share SSCTEST-LOCALD --verbose 
```

The project_name and share parameters can be obtained from the StorCycle GUI, or obtained from the GUI:
```shell
$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command get_locations 

$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command get_migrate_projects 
```
 

## Author
developer@spectralogic.com

