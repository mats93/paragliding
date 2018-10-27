# IGC track viewer extended
Assignment 2 in the course IMT2681-2018 (Cloud Technologies) at NTNU Gjøvik.

This application is a RESTful API to upload and retrieve information about [IGC files.](https://www.fai.org/sites/default/files/documents/igc_fr_spec_with_al4a_2016-4-10.pdf)

***

## Available API calls:
### Tracks:
Information:
```
Allows multiple people to upload and browse IGC files. IGC is an international file format for soaring track files
that are used by paragliders and gliders. The program will store IGC files metadata in a NoSQL Database (MongoDB).

To upload a new IGC file, do a POST request to "/paragliding/api/track" with a json object:
{
    "url": {
      "type": "string"
    }
}
```
```
GET:  /paragliding/api                     - Returns information about the API.
POST: /paragliding/api/track               - Takes the URL in an json format and inserts a new track, returns the tracks ID.
GET:  /paragliding/api/track               - Returns an array of all tracks IDs.
GET:  /paragliding/api/track/<id>          - Returns metadata about a given track with the provided '<id\>'.
GET:  /paragliding/api/track/<id>/<field>  - Returns single detailed metadata about a given tracks field with the provided '<id\>' and '<field\>'.
```

### Ticker:
Information:
```
To facilitate sharing track information, the API allows people to search through the entire collection of tracks,
by track ID, to obtain the details about a given track.
The purpose of the ticker API is to notify dependant applications (such as complex IGC visualisation webapps)
about the new tracks being available.
```
```
GET:  /paragliding/api/ticker/latest        - Returns the timestamp of the latest added track.
GET:  /paragliding/api/ticker/              - Returns the JSON struct representing the ticker for the IGC tracks (array of max 5).
GET:  /paragliding/api/ticker/<timestamp>   - Returns the JSON struct representing the ticker for the IGC tracks, returns only higher timestamps then the one provided.
```

### Webhooks:
Information:
```
Allow subscribing a webhook such that it can notify subscribers about the events in the system.
Currently supports Discord Webhooks.

To register a new webhook, do a POST request to "/paragliding/api/webhook/new_track" with a json object:
{
    "webhookURL": {
      "type": "string"
    },
    "minTriggerValue": {
      "type": "number"
    }
}
Where "minTriggerValue" is how many tracks that need to be added before your webhook gets notified.
This field is optional, if not provided it will be set to 1.
```
```
POST:   /paragliding/api/webhook/new_track/              - Registration of new webhook for notifications about tracks being added to the system. Returns the details about the registration
GET:    /paragliding/api/webhook/new_track/<webhook_id>  - Accessing registered webhooks.
DELETE: /paragliding/api/webhook/new_track/<webhook_id>  - Deleting registered webhooks.
```

***

## How this app is deployed:
 * The app runs in Heroku at https://paragliding-api.herokuapp.com/
 * The database (MongoDB) is stored in [mlab.com](https://mlab.com/) MongoLabs Sandbox.
 * The ticker end-point is implemented within the Heroku deployment (as part of the API).
 * The "clock trigger" is implemented in Go as an independent executable deployed on OpenStack [(link)](https://github.com/mats93/Clock-trigger)

***
## Application testing

### Static code analysis:
Done for each folder/package:
- [x] go build     (go build .)
- [x] go fmt       (go fmt .)
- [x] golint       (golint .)
- [x] go tool vet  (go tool vet .)
- [x] gometalinter (gometalinter -- metalinter .)

Warnings & errors:
```
'go tool vet .' complained that it could not find packages, yet the program finds them when running.

'gometalinter -- metalinter .' complains about the same as 'go tool vet', but also about errors in
other packages such as the mgo package.
```

### Unit test coverage:
 * admin   - 76.0%
 * mongodb - 87.9%
 * ticker  - 91.0%
 * track   - 84.9%
 * webhook - 57.5%

***

## Additional information:

[Link to the task details](https://github.com/mats93/paragliding/blob/testing/TaskDetails.md)

Created by Mats Ove Mandt Skjærstein, 2018
