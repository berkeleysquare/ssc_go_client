# ApiProjectSchedule

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DaysOfMonth** | Pointer to **[]int32** | For monthly and yearly schedules, list the calendar days of the month to execute.  If this number is past the last day of that month, the last day of the month is assumed. | [optional] 
**DaysOfWeek** | Pointer to **[]string** | For weekly and monthly schedules, list the days of the week in those weeks and months on which to execute. | [optional] 
**Interval** | Pointer to **int32** | The period should happen every n days/weeks/months/years.  This should be 1 for daily/weekly/monthly/yearly, 2 for every other day/week/month/year, 3 for every third, etc. | [optional] 
**MonthsOfYear** | Pointer to **[]string** | For monthly and yearly schedules, list the months in which to execute. | [optional] 
**Period** | Pointer to **string** | How often to execute projects associated with this schedule. | 
**StartDay** | Pointer to **int32** | The day on which to start, from 1 to 31.  Ignored if the schedule period is Now. | [optional] 
**StartHour** | Pointer to **int32** | The hour to start, from 0 to 23.  Ignored if the schedule period is Now. | [optional] 
**StartMinute** | Pointer to **int32** | The minute of the hour to start, from 0 to 59.  Ignored if the schedule period is Now. | [optional] 
**StartMonth** | Pointer to **int32** | The month in which to start, from 1 to 12.  Ignored if the schedule period is Now. | [optional] 
**StartYear** | Pointer to **int32** | The year in which to start. Ignored if the schedule period is Now. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


