# StorCycle® Client and test utility

## Overview
This utility provides a CLI to interact with the StorCycle REST interface and
to interact directly with the underlying database. All commands requiring db access
(have db in the name) must be run on the StorCyle server.

If --in is specified, it will read filenames from the CSV file
(First column, second row). Else the project must be specified as --project_name

This file can be the output of the get_migrate_projects command:
``` shell
--url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command get_migrate_projects --out projectList.csv
```


## search_db_project
- Search for all objects in all jobs of project (--project_name parameter)
- Output to csv (console or file name specified by --out parameter)

## verify_db_project
- Search for all objects in all jobs of project (--project_name parameter)
- Verify that object exists (has been restored to) a source (--share parameter) and directory (--directory parameter)
- Output to csv (console or file name specified by --out parameter)

## assert_db_project
- Search for all objects in all jobs of project (--project_name parameter)
- Verify that object exists (has been restored to) a source (--share parameter) and directory (--directory parameter)
- Output to csv ONLY IF FILE IS NOT FOUND (console or file name specified by --out parameter)

## Operation
Unpack the executable onto a machine with network access to the StorCycle server.

On Windows, A directory under C:\StorCycle is recommended, e.g., C:\StorCycle\verify

### Display available commands:
```shell
$ ssc-cli --command list_commands
```

### Display available parameters:
```shell
$ ssc-cli --help
```
### Display all migrate projects:
```shell
$ ssc-cli --name Administrator --password spectra --verbose  --ignore_cert  --command get_migrate_projects
```
### Search db example:
```shell
$ ./ssc-cli --name Administrator --password spectra --verbose  --ignore_cert  --command search_db_project --project_name picnic --out searchObjects.csv
```
### Verify db example:
```shell
$ ./ssc-cli --name Administrator --password spectra --verbose  --ignore_cert  --command verify_db_project --in projectList.csv --share shareOne --directory restore --out verifyObjects.csv
```
Include --verbose to write log output to the console. 

Include --jobToPath to append manifest name to share path.
(If --directory is specified, it will appear after share path and job name).

If --out is not specified, it will print to the console.

NOTE: files will be written to csv in batches of 500,000 (can be overwritten with --count parameter); additional pages will be named with the offset:
```
searchOutput.csv
searchOutput_500000.csv
searchOutput_1000000.csv
searchOutput_1500000.csv
...
```

## Author
developer@spectralogic.com
