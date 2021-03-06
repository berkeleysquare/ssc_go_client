/*
 * Spectra Logic Corporation Software
 *
 * StorCycle® REST API Copyright © 2020 Spectra Logic Corporation
 *
 * API version: 0.1
 * Contact: developer@spectralogic.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi
// ApiProject struct for ApiProject
type ApiProject struct {
	// Description of the project.
	Description *string `json:"description,omitempty" xml:"description"`
	// The source share with which the project will work. This is not required when updating the project.
	Share *string `json:"share" xml:"share"`
	// For restore projects, the location within the share where the files are restored. For archive projects, the starting subdirectory on the share to archive. For scan projects, the directory on the share to scan.  Default is the share root.
	WorkingDirectory *string `json:"workingDirectory,omitempty" xml:"workingDirectory"`
	// Tags to attach to jobs executed by this project.
	Tags *[]string `json:"tags,omitempty" xml:"tags"`
	// Used to re-activate a deactivated project.
	Active *bool `json:"active,omitempty" xml:"active"`
	// Used to pause/resume a project.
	Enabled *bool `json:"enabled,omitempty" xml:"enabled"`
	// The type of project.
	ProjectType *string `json:"projectType,omitempty" xml:"projectType"`
	// For archive projects, the action that is performed on the original source file after successful archive. For restore breadcrumb projects, this specifies the type of bread crumb to create.
	BreadCrumbAction *string `json:"breadCrumbAction,omitempty" xml:"breadCrumbAction"`
	// A list of users from which email addresses may be found for sending notifications.
	EmailOnComplete *[]string `json:"emailOnComplete,omitempty" xml:"emailOnComplete"`
	// The name of the manifest containing the files to be restored. This is synonymous with the original archive job that originally archived the files.
	RestoreManifest *string `json:"restoreManifest,omitempty" xml:"restoreManifest"`
	// The list of versions from the manifest to restore.
	RestoreVersions *[]string `json:"restoreVersions,omitempty" xml:"restoreVersions"`
	// The storage endpoints where files will be archived.
	Targets *[]string `json:"targets,omitempty" xml:"targets"`
	Filter ApiProjectFilter `json:"filter,omitempty" xml:"filter"`
	Schedule ApiProjectSchedule `json:"schedule" xml:"schedule"`
	Status ApiProjectStatus `json:"status,omitempty" xml:"status"`
}
