package client

import (
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"log"
	"strings"
)

const JOB_STATUS_COMPLETE = "Completed"
const JOB_STATUS_FAILED = "Failed"
const JOB_STATUS_CANCELED = "Canceled"
const JOB_STATUS_POST_ACTION_PENDING = "PostActionPending"
const JOB_STATUS_FAILED_OOST_ACTION = "FailedPostActionPending"

func isStateActive(state string) bool {
	return state != JOB_STATUS_COMPLETE &&
		state != JOB_STATUS_FAILED &&
		state != JOB_STATUS_POST_ACTION_PENDING &&
		state != JOB_STATUS_CANCELED &&
		state != JOB_STATUS_FAILED_OOST_ACTION
}

func RunNow(ssc *SscClient, args *Arguments) error {
	_, err := ssc.Client.ProjectApi.RunProjectNow(*ssc.Context, args.ProjectName)
	if err == nil {
		fmt.Printf("Successfully set Run Now for %s\n", args.ProjectName)
		return nil
	}
	return fmt.Errorf("run now failed for project %s\n%v\n",  args.ProjectName, err)
}

func CancelJob(ssc *SscClient, args *Arguments) error {
	_, err := ssc.Client.ProjectApi.CancelJob(*ssc.Context, args.Job)
	if err == nil {
		fmt.Printf("Canceled %s\n", args.Job)
		return nil
	}
	return fmt.Errorf("cancel failed for job %s\n%v\n",  args.Job, err)
}

func displayJobs(jobs openapi.ApiJobPaginator) error {
	for jobIndex := range jobs.Data {
		job := jobs.Data[jobIndex]
		fmt.Printf("Job: %s, Description: %s\n", *job.Name, *job.Description)
	}
	return nil
}

func displayJobsWithStatus(jobs openapi.ApiJobsWithStatusPaginator, match string, activeOnly bool) error {
	for jobIndex := range jobs.Data {
		job := jobs.Data[jobIndex]
		if len(match) == 0 || strings.Contains(*job.Name, match) {
			if !activeOnly || isStateActive(*job.State) {
				fmt.Printf("Job: %s, State: %s, Complete: %.2f\n", *job.Name, *job.State, *job.PercentComplete)
			}
		}
	}
	return nil
}

func formatJobStatus(job openapi.ApiJobStatus) string {
	return fmt.Sprintf("Job: %s, State: %s, Complete: %.2f\n", *job.Name, *job.State, *job.PercentComplete)
}

func formatJobWithState(job openapi.ApiJobWithState) string {
	return fmt.Sprintf("Job: %s, State: %s, Complete: %.2f\n", *job.Name, *job.State, *job.PercentComplete)
}

func displayJobStatus(job openapi.ApiJobStatus) error {
	fmt.Printf(formatJobStatus(job))
	return nil
}

func listJobs(ssc *SscClient, args *Arguments) error {
	response, _, err := ssc.Client.ProjectApi.SearchJobs(*ssc.Context, args.ProjectName, nil)
	if err != nil {
		return fmt.Errorf("search jobs for project name %s failed %v\n", args.ProjectName, err)
	}

	return displayJobs(response)
}

func GetJobStatus(ssc *SscClient, args *Arguments) error {
	response, _, err := ssc.Client.ProjectApi.GetJobStatus(*ssc.Context, args.Job)
	if err != nil {
		return fmt.Errorf("get job status for job name %s failed %v\n", args.Job, err)
	}
	return displayJobStatus(response)
}

func GetAllJobs(ssc *SscClient, args *Arguments) error {
	jobType := ""
	if args.Command == "get_all_restore_jobs" ||  args.Command == "get_active_restore_jobs"{
		jobType = "Restore"
	}
	activeOnly := args.Command == "get_active_restore_jobs"
	response, _, err := ssc.Client.ProjectApi.GetAllJobs(*ssc.Context, jobType, nil)
	if err != nil {
		return fmt.Errorf("get jobs failed %v\n", err)
	}

	//  cancel (--cancel) or just list
	if args.Cancel {
		for jobIndex := range response.Data {
			job := response.Data[jobIndex]
			if isStateActive(*job.State) {
				_, err = ssc.Client.ProjectApi.CancelJob(*ssc.Context, *job.Name)
				if err == nil {
					fmt.Printf("Cancel job %s\n", *job.Name)
				} else {
					// best effort to cancel -- keep trying.
					fmt.Printf("Error: cancel job %s failed %v\n", *job.Name, err)
				}
			}
		}
		return nil
	}

	return displayJobsWithStatus(response, args.Prefix, activeOnly)
}

func includesTag(tags []string, match string) bool {
	for tagIndex := range tags {
		if tags[tagIndex] == match {
			return true
		}
	}
	return false
}

func GetRestoreJobsByTag(ssc *SscClient, args *Arguments) error {
	activeOnly := args.Command == "active_restore_jobs_by_tag" || args.Command == "wait_for_restore_jobs_by_tag"
	matchingJobs, err := doRestoreJobsByTag(ssc, args.Tag, activeOnly)
	if err != nil {
		return fmt.Errorf("getRestoreJobsByTag failed %v\n", err)
	}
	//  cancel (--cancel) or just list
	if args.Cancel {
		for jobIndex := range matchingJobs {
			job := matchingJobs[jobIndex]
			if isStateActive(*job.State) {
				_, err = ssc.Client.ProjectApi.CancelJob(*ssc.Context, *job.Name)
				if err == nil {
					fmt.Printf("Cancel job %s\n", *job.Name)
				} else {
					// best effort to cancel -- keep trying.
					fmt.Printf("Error: cancel job %s failed %v\n", *job.Name, err)
				}
			}
		}
		return nil
	}

	for jobIndex := range matchingJobs {
		fmt.Printf(formatJobWithState(matchingJobs[jobIndex]))
	}
	return nil
}


func doRestoreJobsByTag(ssc *SscClient, match string, activeOnly bool) ([]openapi.ApiJobWithState, error) {
	jobType := "Restore"
	if len(match) == 0 {
		return nil, fmt.Errorf("must supply --tag parameter")
	}
	jobs, _, err := ssc.Client.ProjectApi.GetAllJobs(*ssc.Context, jobType, nil)
	if err != nil {
		return nil, fmt.Errorf("get jobs failed %v\n", err)
	}
	// filter by tag
	var matchingJobs []openapi.ApiJobWithState
	for jobIndex := range jobs.Data {
		job := jobs.Data[jobIndex]
		tags := job.Tags
		if tags != nil && includesTag(*tags, match) {
			if !activeOnly || isStateActive(*job.State) {
				matchingJobs = append(matchingJobs, job)
			}
		}
	}
	return matchingJobs, nil
}

func WaitForRestoreJobsByTag(ssc *SscClient, args *Arguments) error {
	activeOnly := args.Command == "active_restore_jobs_by_tag" || args.Command == "wait_for_restore_jobs_by_tag"
	err := tryRestoreJobsByTag(ssc, args.Tag, activeOnly,  args.Verbose, 13, 21)
	if err != nil {
		return fmt.Errorf("tryRestoreJobsByTags() failed %v", err)
	}

	// all are complete, print all
	allMatchingJobs, err := doRestoreJobsByTag(ssc, args.Tag, false)
	if err != nil {
		return fmt.Errorf("getRestoreJobsByTag failed %v\n", err)
	}
	for jobIndex := range allMatchingJobs {
		fmt.Printf(formatJobWithState(allMatchingJobs[jobIndex]))
	}
	return nil
}

func tryRestoreJobsByTag(ssc *SscClient, match string, activeOnly bool, verbose bool, fib1 int, fib2 int) error {
	matchingJobs, err := doRestoreJobsByTag(ssc, match, activeOnly)
	if err != nil {
		return fmt.Errorf("doRestoreJobsByTag failed %v\n", err)
	}
	if len(matchingJobs) == 0 {
		// all done
		return nil
	}
	// wait
	fib3 := staggeredWait(fib1, fib2)
	// try again
	if verbose {
		log.Printf("%d jobs not complete", len(matchingJobs))
	}
	return tryRestoreJobsByTag(ssc, match, activeOnly, verbose, fib2, fib3)
}
