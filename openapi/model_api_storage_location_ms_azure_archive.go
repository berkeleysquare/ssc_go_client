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
// ApiStorageLocationMsAzureArchive Used to create MSAzure Archive storage locations.
type ApiStorageLocationMsAzureArchive struct {
	// For use only with MSAzure Archive.  The access key ID of the credentials for MSAzure Archive.
	AccessKeyId *string `json:"accessKeyId" xml:"accessKeyId"`
	// Used to re-activate a 'deleted' storage location.
	Active *bool `json:"active,omitempty" xml:"active"`
	// Used for MSAzure Archive endpoints.  The container to use for storage.
	Bucket *string `json:"bucket" xml:"bucket"`
	// Cost per TB of storage in units of department currency.  Only valid when Department is set.
	CostPerTB *float32 `json:"costPerTB,omitempty" xml:"costPerTB"`
	// ID of department responsible for location storage cost.
	Department *string `json:"department,omitempty" xml:"department"`
	// Description of the storage location.
	Description *string `json:"description,omitempty" xml:"description"`
	// Used to enable encryption that has been configured for the product.
	EnableEncryption *bool `json:"enableEncryption,omitempty" xml:"enableEncryption"`
	// Prevent location from use as a source. Set to true for MSAzure Archive.
	IsTarget *bool `json:"isTarget" xml:"isTarget"`
	// The limit during off-peak hours, expressed in total number of bytes per second.  Specify 0 for no limit.
	OffPeakBytesPerSecondLimit *float64 `json:"offPeakBytesPerSecondLimit,omitempty" xml:"offPeakBytesPerSecondLimit"`
	// The file scan limit during off-peak hours, expressed in total number of files per second. Specify 0 for no limit.
	OffPeakFilesPerSecondScanLimit *float64 `json:"offPeakFilesPerSecondScanLimit,omitempty" xml:"offPeakFilesPerSecondScanLimit"`
	// The limit during peak hours, expressed in total number of bytes per second.  Specify 0 for no limit.
	PeakBytesPerSecondLimit *float64 `json:"peakBytesPerSecondLimit,omitempty" xml:"peakBytesPerSecondLimit"`
	// The file scan limit during peak hours, expressed in total number of files per second. Specify 0 for no limit.
	PeakFilesPerSecondScanLimit *float64 `json:"peakFilesPerSecondScanLimit,omitempty" xml:"peakFilesPerSecondScanLimit"`
	// A list of the days of the week where rate limits are enforced during peak and off-peak times.  If a day of week is not in the list, it is assumed to be unlimited.
	PeakTimesSchedule *[]ApiPeakTimeSchedule `json:"peakTimesSchedule,omitempty" xml:"peakTimesSchedule"`
	// The number of days to retain files on this target before scheduling files for deletion. This field is only applicable to targets.
	RetentionDays *int32 `json:"retentionDays,omitempty" xml:"retentionDays"`
	// For use only with MSAzure Archive.  The secret access key of the credentials for MSAzure Archive. This is write-only.
	SecretAccessKey *string `json:"secretAccessKey" xml:"secretAccessKey"`
	Type *string `json:"type" xml:"type"`
}
