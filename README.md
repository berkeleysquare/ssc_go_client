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
$ ssc-cli --command list_commands
```

Display available parameters:
```shell
$ ssc-cli --help
```
###Full cycle test

The full cycle test verifies StorCycle operation and demonstrates full persistence through to tape 
as well as data integrity. It must communicate both with StorCycle and the BlackPearl,
so both credentials need to be provided. StorCycle credentials are passed in as CLI parameters
 (e.g., --name Administrator --password spectra). The BlackPearl credentials can be set as environment variables or passed in 
to ssc-cli (e.g., --endpoint http://10.85.41.36 --secret_key btkDKJBd --access_key ams=)
* `DS3_ENDPOINT` - The URL to the DS3 Endpoint 
* `DS3_ACCESS_KEY` - The DS3 access key
* `DS3_SECRET_KEY` - The DS3 secret key
* `http_proxy` - If set, the `Client` instance will proxy through this URL

 
###Run full cycle verify example:
First create a BlackPearl target storage location using the GUI
or API. The test location will clone that to acquire credentials. 
In the example, bp-sandbox is the location to be cloned.
```shell
$ ssc-cli  --url https://localhost/openapi --name Administrator --password spectra --ignore_cert --command full_verify --directory C:\\StorCycle\\verify --clone bp-sandbox --endpoint http://10.85.41.36 --secret_key btkDKJBd --access_key ams=
```
The full_cycle command will write status to the console at each task. 
Steps which must wait for tasks to complete  will print a period (.) to the console each attempt.
```shell
$ ssc-cli.exe --command directory_checksum --directory C:\\StorCycle\\shares\\two --out shares_two.csv
```

### generate checksums
Two commands generate crc64 checksums to match those tracked in StorCycle. To generate a checksum for a single file:
```shell
$ ssc-cli.exe --command checksum_test_file --file_name verify-test-file220526083234.txt --directory C:\\StorCycle\\verify\\verify-test-source  
```

To recurse an entire directory and output to a .csv file:
```shell
$ ssc-cli.exe --command directory_checksum --directory C:\\StorCycle\\shares\\two --out shares_two.csv
```
If --out is not specified, it will print to the console.

## Author

developer@spectralogic.com

