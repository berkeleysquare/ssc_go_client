# Automating the Breadcrumber

Two new commands have been added to the utility to support automation.
## process_projects
This command will read a list of projects from a CSV file
and get a list of jobs for each project (newest -> oldest).
It will go through the list until it finds a job that has completed.

As jobs process, text files are left in the home directory named by job ID and status.

To generate the CSV file, use the following command:
```shell
$ ssc-cli --command get_migrate_projects --out projects.csv --url "https://localhost/openapi" --name "Administrator" --password spectra --ignore_cert 
```
This will create projects.csv with all migrate projects. Edit the file to include 
only the projects you want to process.

Then execute (or schedule) the process_projects command:
```shell
$ ssc-cli --command "process_projects" --in projects.csv --template "templates/YNH.html" --log breadcrumb_ --suffix ".ARCHIVE.html" --url "https://localhost/openapi" --start 0  --name "Administrator" --password spectra --ignore_cert --verbose
```

## mail_breadcrumb_report
To ensure that the operations are completing, a report can be generated and emailed.

First, edit the mail_config.yaml file to include your email server and credentials. 
The sample uses an internal SMTP server in our sandbox
```
email:
  server: emailsmarthost.spectralogic.com
  port: 25
  from: noreply@spectralogic.com
  password:
  authorization:
message:
  subject: Breadcrumb Report
  template: templates/YNH_report.html
  to:
   - johnk@spectralogic.com
```

Then execute the following command:
```shell
$ ssc-cli --command mail_breadcrumb_report --in projects.csv --url "https://localhost/openapi" --name "Administrator" --password spectra --ignore_cert  --verbose --name "Administrator" --password spectra --ignore_cert
```

When connectivity has been established and the mail format approved,
this job can be scheduled to run after the process_projects command is expected to have completed.

If the --log and --verbose options are used, the utility will write a log file to the home directory.

To rerun a job, delete the <job_name>_succeeded.txt file from the home directory. 
Then run the process_projects command or the original:
```shell
 $ makeLinks.BAT <job_name> <password>
```
or:
```shell
$ ssc-cli --url "https://localhost/openapi" --name "Administrator" --password spectra --ignore_cert --command "write_breadcrumbs" --job ArkHive-1 --in "templates/YNH.html" --start 0 --log ArkHive-1 --suffix ".ARCHIVE.html" --verbose
```
