package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/mail"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

func makeStartFileName(jobName string) string {
	return jobName + "_start.txt"
}

func makeSuccessFileName(jobName string) string {
	return jobName + "_succeeded.txt"
}

func makeWarningFilename(jobName string) string {
	return jobName + "_warning.txt"
}
func makeErrorFilename(jobName string) string {
	return jobName + "_error.txt"
}

func createStartFile(jobName string) error {
	filename := makeStartFileName(jobName)
	message := "started"
	return createTimestampFile(filename, message)
}

func startFileExists(jobName string) bool {
	filename := makeStartFileName(jobName)
	_, err := os.Stat(filename)
	return err == nil
}

func warningFileExists(jobName string) bool {
	filename := makeWarningFilename(jobName)
	_, err := os.Stat(filename)
	return err == nil
}

func successFileExists(jobName string) bool {
	filename := makeSuccessFileName(jobName)
	_, err := os.Stat(filename)
	return err == nil
}

func errorFileExists(jobName string) bool {
	filename := makeErrorFilename(jobName)
	_, err := os.Stat(filename)
	return err == nil
}

func createSuccessFile(jobName string, count int64) error {
	filename := makeSuccessFileName(jobName)
	message := fmt.Sprintf("wrote %d links", count)
	return createTimestampFile(filename, message)
}

func createWarningFile(jobName string, message string) error {
	filename := makeWarningFilename(jobName)
	return createTimestampFile(filename, message)
}

func createErrorFile(jobName string, jobError error) error {
	filename := makeErrorFilename(jobName)

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create success file %s\n%v", filename, err)
	}
	defer file.Close()

	// Get the current time and format it as a timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	body := fmt.Sprintf("%s: ERROR %v\n", timestamp, jobError)
	_, err = file.WriteString(body)
	if err != nil {
		return fmt.Errorf("Failed to write to error file %s\n%v", filename, err)
	}
	return nil
}

func createTimestampFile(filename string, message string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed to create file %s\n%v", filename, err)
	}
	defer file.Close()

	// Get the current time and format it as a timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	_, err = file.WriteString(timestamp + ": " + message + "\n")
	if err != nil {
		return fmt.Errorf("Failed to write to file %s\n%v", filename, err)
	}
	return nil
}

func breadcrumbReport(ssc *SscClient, args *Arguments) error {
	mail := args.Command == "mail_breadcrumb_report"
	subject := "Breadcrumb Report: All Jobs Succeeded"
	report, err := doBreadCrumbReport(ssc, args)
	if err != nil {
		subject = "Breadcrumb Report: Errors found"
	}
	if mail {
		return mailReport(args, subject, report)
	}
	return printReport(args, subject, report)
}

func doBreadCrumbReport(ssc *SscClient, args *Arguments) (*mail.BreadcrumbReport, error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	ret := &mail.BreadcrumbReport{
		Header:      "Breadcrumb Report: " + timestamp,
		Report:      make([]string, 0),
		Attachments: make([]string, 0),
	}
	var line string
	allSucceeded := true
	verbose := args.Verbose

	projects, err := getJobsToProcess(ssc, args)
	if err != nil {
		line = fmt.Sprintf("could not get projects %v", err)
		ret.Append(line)
		return ret, fmt.Errorf("%s", line)
	}
	ret.Append(fmt.Sprintf("Verifying %d projects", len(*projects)))
	for projectName := range *projects {
		ret.Append(fmt.Sprintf("Verifying project: %s", projectName))
		for _, jobName := range (*projects)[projectName] {
			if verbose {
				ret.Append(fmt.Sprintf("Processing job: %s", jobName))
			}
			if successFileExists(jobName) {
				if warningFileExists(jobName) {
					ret.Append(fmt.Sprintf("Job %s for %s has warnings", projectName, jobName))
					ret.Attachments = append(ret.Attachments, makeWarningFilename(jobName))
				}
				if verbose {
					ret.Append(fmt.Sprintf("Job %s for %s has completed", projectName, jobName))
				}
				continue
			} else {
				if errorFileExists(jobName) {
					ret.Append(fmt.Sprintf("Job %s for %s has an error", projectName, jobName))
					allSucceeded = false
					ret.Attachments = append(ret.Attachments, makeErrorFilename(jobName))
					continue
				}
				if startFileExists(jobName) {
					ret.Append(fmt.Sprintf("Job %s for %s did not complete", projectName, jobName))
					allSucceeded = false
					continue
				}
				ret.Append(fmt.Sprintf("Job %s for %s has not started", projectName, jobName))
			}
		}
	}
	if allSucceeded {
		ret.Append(fmt.Sprintf("All jobs succeeded"))
		return ret, nil
	}
	return ret, fmt.Errorf("not all jobs succeeded")
}

func mailReport(args *Arguments, subject string, report *mail.BreadcrumbReport) error {
	// get the cruise/mail config from the yaml file
	configYaml, err := os.ReadFile(mailConfigFile)
	if err != nil {
		return fmt.Errorf("Could not open config file %s\n%v\n", mailConfigFile, err)
	}
	config := &mail.Config{}

	err = yaml.Unmarshal(configYaml, config)
	if err != nil {
		return fmt.Errorf("Could not parse config file %s\n%v\n", mailConfigFile, err)
	}
	config.Message.Subject = subject
	return mail.Mail(config, report)
}

func printReport(args *Arguments, subject string, report *mail.BreadcrumbReport) error {
	// print to console
	fmt.Printf("%s\n", subject)
	fmt.Printf("%s\n", report.Header)
	for _, line := range report.Report {
		fmt.Println(line)
	}
	return nil
}
