# StorCycle® Client and test utility

## Overview
This utility provides a limited CLI for Spectra Logic's StorCycle, 
including a full-cycle verification.

## full-verify 

- Create test file
- Generate checksum
- Create source storage location
- Clone BlackPearl target location
- Create Migrate project; run mow
- Wait for file to be placed in cache
- Wait for file to be placed on tape
- Display tape location (id and barcode)
- Create Restore project; run now
- Wait for restored file
- Generate checksum on restored file
- Compare checksums on source and restored file.

## Operation
Unpack the executable onto a machine with network access to the StorCycle server.

On Windows, A directory under C:\StorCycle is recommended, e.g., C:\StorCycle\verify

Display available commands:
```shell
$ ssc_cli --command list_commands
```

Display available parameters:
```shell
$ ssc_cli --help
```

Run full cycle verify example:
First create a BlackPearl target storage location using the GUI
or API. The test location will clone that to acquire credentials. 
In the example, bp-sandbox is the location to be cloned.
```shell
$ ssc_cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command full_verify --directory C:\\StorCycle\\verify --clone bp-sandbox
```
The full_cycle command will write status to the console at each task. 
Steps which must wait for tasks to complete  will print a period (.) to the console each attempt.

## Author

developer@spectralogic.com

