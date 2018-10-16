# IGC track viewer extended
Assignment 2 in the course IMT2681-2018 (Cloud Technologies) at NTNU Gjøvik.

This application is a RESTful API to upload and browse [IGC files.](https://www.fai.org/sites/default/files/documents/igc_fr_spec_with_al4a_2016-4-10.pdf)

## Available API calls:
### Tracks:
 * `GET:  /paragliding/api`                    - Returns information about the API.
 * `POST: /paragliding/api/track`              - Takes the URL in an json format and inserts a new track, returns the tracks ID.
 * `GET:  /paragliding/api/track`              - Returns an array of all tracks IDs.
 * `GET:  /paragliding/api/track/<id>`         - Returns metadata about a given track with the provided '<id\>'.
 * `GET:  /paragliding/api/track/<id>/<field>` - Returns single detailed metadata about a given tracks field with the provided '<id\>' and '<field\>'.

### Ticker:
 * `GET:  paragliding/api/ticker/latest`       - Returns the timestamp of the latest added track.
 * `GET:  paragliding/api/ticker/`             - Returns the JSON struct representing the ticker for the IGC tracks (array of max 5).
 * `GET:  paragliding/api/ticker/<timestamp>`  - Returns the JSON struct representing the ticker for the IGC tracks, returns only higher timestamps then the one provided.

### Webhooks:
 * `POST: paragliding/api/webhook/new_track/`  - Registration of new webhook for notifications about tracks being added to the system. Returns the details about the registration
 * `GET:  paragliding/api/webhook/new_track/<webhook_id>` - Accessing registered webhooks.
 * `DELETE: paragliding/api/webhook/new_track/<webhook_id>` - Deleting registered webhooks.

### Admin:
 * `GET:  paragliding/admin/api/tracks_count`  - Returns the current count of all tracks in the DB.
 * `DELETE: paragliding/admin/api/tracks`      - Deletes all tracks in the DB.


## Additional information:
The app runs in Heroku at https://paragliding-api.herokuapp.com/

The program will store IGC files metadata in a NoSQL Database (MongoDB).


Created by Mats
