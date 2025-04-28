# StorCycle® Client and test utility

## Overview
This utility provides a CLI search/restore for Spectra Logic's StorCycle
matching the directory or file name. 
It will search for the string and filter results by file extensions. 
Matching objects can be batched into restore projects by migrate job.

NOTE: Search must be from the first character of a directory or file name. 
Trailing wildcard is assumed but leading wildcards are not supported.

## search_objects 
- Search for all matching objects (--file_name parameter)
- Filter by file extensions (--ext parameter: comma-delimited)
- Output to csv (console or file name specified by --out parameter)

## restore_objects 
- Search for all matching objects (--file_name parameter)
- Filter by file extensions (--ext parameter: comma-delimited)
- Output to csv (only if file name specified by --out parameter)
- Restore to original location (default) or specify --share (the name of a StorCycle Storage Location) and --directory (a subdirectory in the location -- will be created if not extant)
- Create one Restore Project per migrate job
- Run job now

## Operation
Unpack the executable onto a machine with network access to the StorCycle server.

On Windows, A directory under C:\StorCycle is recommended, e.g., C:\StorCycle\verify

Display available commands:
```shell
$ ssc-cli --command list_commands
```

Display available parameters:
```shell
$ ssc-cli --help
```
###Search example:
```shell
$ ssc-cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command search_objects --file_name picnic -ext mp4 
```

If --in is specified, it will read filenames from the CSV file
(First column, second row). Else the search string must be specified as --file_name
```shell
$ ssc-cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command search_objects --in dir_list.csv -ext jpg,gif 
```
Include --verbose to write log output to the console. 
Include --verbose and --log <logfile> to capture output to a file.

No --exts parameter, or --exts "*" will return all files matching the search string:
```shell
$ ssc-cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command search_objects --in dir_list.csv -ext "*" 
-- or --
$ ssc-cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command search_objects --in dir_list.csv 
```

###Restore example:
```shell
$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_objects --file_name picnic --ext mp4,mp3,jpg --share Restorey --directory /testAuto --out myFiles.csv
```
Will output CSV only if --out is specified.

Include --verbose to write log output to the console. 

If --out is not specified, it will print to the console.

If --in is specified, it will read filenames from the CSV file
(First column, second row). Else the search string must be specified as --file_name
```shell
$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_objects --in directory_list --ext mp4,mp3,jpg --share Restorey --directory /testAuto --out myFiles.csv
```
No --exts parameter, or --exts "*" will return all files matching the search string.
```shell
$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_objects --file_name picnic --ext "*" --share Restorey --directory /testAuto --out myFiles.csv
-- or --
$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_objects --file_name picnic --share Restorey --directory /testAuto --out myFiles.csv
```
### Restore status
All projects created have a tag of "Restore &lt;caseid&gt;" ("Restore" + space + the search/restore term).
Jobs can be queried using this tag to check status
```shell
>> .\ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_jobs_by_tag --tag "Restore MonaLisa"  --verbose
Job: Restore_MonaLisa__tuesdays+CHild-2_23-02-02-08-47-35.321-1, State: Completed, Complete: 100.00
Job: Restore_MonaLisa__full-1_23-02-02-08-47-35.314-1, State: Active, Complete: 37.50
```
 
And a command can block until the jobs are complete.
```shell
>> .\ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command wait_for_restore_jobs_by_tag --tag "Restore MonaLisa"  --verbose
2023/02/02 08:48:28 1 jobs not complete
2023/02/02 08:49:23 1 jobs not complete
Job: Restore_MonaLisa__tuesdays+CHild-2_23-02-02-08-47-35.321-1, State: Completed, Complete: 100.00
Job: Restore_MonaLisa__full-1_23-02-02-08-47-35.314-1, State: Completed, Complete: 100.00
```

##Search and Restore Database Direct
Accessing the database directly improves performances and enables full regEx search across the entire path. 

NOTE: for security, db search/restore can only be run on the Storcycle Server.

### Search DB
Same syntax as search_objects, but the command is search_db:
```shell
$ ssc-cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command search_db --file_name picnic -ext mp4 
```
### Restore DB
Same syntax as restore_objects, but the command is restore_db_objects:
```shell
$ ssc-cli --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_db_objects --file_name picnic -ext mp4 --share Restorey --directory /testAuto --out myFiles.csv
```
NOTE: --ext * or not including the --ext parameter will return all that match the file name string 
on the search_db and restore_db_objects commands.

### Fast (pinned) Search and Restore
Searching a large database for an occurrence of a string anywhere 
in the path can be slow. "Pinning" the search to the start of a 
component path (directory name or file name) can speed up the search
-- seconds vs. tens of minutes -- in large DBs.

The search_db and restore_db_objects commands accept a --fast_search 
parameter. If included, the search will be pinned to the start of
a component path.

E.g., --file_name "picnic" --fast_search will  match
\\share\users\jk\picnic\recipes.txt
\\share\users\jk\family\picnic.jpg

But not
\\share\users\jk\family_picnic\recipes.txt

Omitting the --fast_search parameter will match all of the above.

```shell





### Find and cancel restore jobs
Find all active restore jobs
```nashorn js
ssc-cli --url https://localhost/openapi --name administrator --password spectra --ignore_cert --command get_active_restore_jobs
```
Add the flag --cancel to cancel all returned jibs

Find active restore jobs by tag
```nashorn js
--url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command restore_jobs_by_tag --tag "Restore Mona"  --verbose
````
Add the flag --cancel to cancel all returned jibs.
To cancel all active restore jobs:
```nashorn js
--url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command get_active_restore_jobs --cancel
````


## Author

developer@spectralogic.com

