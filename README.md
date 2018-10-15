# In-memory IGC track viewer
Assignment 1 in the course IMT2681-2018 (Cloud Technologies) at NTNU Gjøvik.

This application is a RESTful API to upload and browse [IGC files.](https://www.fai.org/sites/default/files/documents/igc_fr_spec_with_al4a_2016-4-10.pdf)

## API calls: 
 * `GET:  /igcinfo/api`                  - Returns information about the API.
 * `POST: /igcinfo/api/igc`              - Takes the URL in an json format and inserts a new track, returns the tracks ID.
 * `GET:  /igcinfo/api/igc`              - Returns an array of all tracks IDs.
 * `GET:  /igcinfo/api/igc/<id>`         - Returns metadata about a given track with the provided '<id\>'.
 * `GET:  /igcinfo/api/igc/<id>/<field>` - Returns single detailed metadata about a given tracks field with the provided '<id\>' and '<field\>'.


## Additional information:
The app runs in Heroku at https://igcinfo-api.herokuapp.com

The app does not store data in any persistant storage, only in memory.

Created by Mats


[Task details](https://github.com/mats93/igcinfo/blob/master/OrignalTask.md)

