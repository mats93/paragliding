# [](#assignment-1-in-memory-igc-track-viewer)Assignment 1: in-memory IGC track viewer

## [](#about)About

Develop an online service that will allow users to browse information about IGC files. IGC is an international file format for soaring track files that are used by paragliders and gliders. The program will not store anything in a persistent storage. Ie. no information will be stored on the server side on a disk or database. Instead, it will store submitted tracks in memory. Subsequent API calls will allow the user to browse and inspect stored IGC files.

For the development of the IGC processing, you will use an open source IGC library for Go: [goigc](https://github.com/marni/goigc)

The system must be deployed on either Heroku or Google App Engine, and the Go source code must be available for inspection by the teaching staff (read-only access is sufficient).

## [](#specification)Specification

### [](#general-rules)General rules

The **igcinfo** should be the root of the URL API. The package and the project repo name can be arbitrary, yet the name must be meaningful. If it is called assignment1 or assignment_1 or variation of this name, we will not mark it.

The server should respond with 404 when asked about the root. The API should be mounted on the **api** path. All the REST verbs will be subsequently attached to the /igcinfo/api/* root.

    http://localhost:8080/igcinfo/               -> 404
    http://localhost:8080/<rubbish>              -> 404
    http://localhost:8080/igcinfo/api/<rubbish>  -> 404

**Note:** the use of `http://localhost:8080` serves only a demonstration purposes. You will have your own URL from the provider, such as Heroku. `<rubbish>` represents any sequence of letters and digits that are not described in this specification.

### [](#get-api)GET /api

*   What: meta information about the API
*   Response type: application/json
*   Response code: 200
*   Body template

```javascript
    {
      "uptime": <uptime>
      "info": "Service for IGC tracks."
      "version": "v1"
    }
```

*   where: `<uptime>` is the current uptime of the service formatted according to [Duration format as specified by ISO 8601](https://en.wikipedia.org/wiki/ISO_8601#Durations).

### [](#post-apiigc)POST /api/igc

*   What: track registration
*   Response type: application/json
*   Response code: 200 if everything is OK, appropriate error code otherwise, eg. when provided body content, is malformed or URL does not point to a proper IGC file, etc. Handle all errors gracefully.
*   Request body template

```javascript
    {
      "url": "<url>"
    }
```

*   Response body template

```javascript
    {
      "id": "<id>"
    }
```

*   where: `<url>` represents a normal URL, that would work in a browser, eg: `http://skypolaris.org/wp-content/uploads/IGS%20Files/Madrid%20to%20Jerez.igc` and `<id>` represents an ID of the track, according to your internal management system. You can choose what format should be in your system. The only restriction is that it needs to be easily used in URLs and it must be unique. It is used in subsequent API calls to uniquely identify a track, see below.

### [](#get-apiigc)GET /api/igc

*   What: returns the array of all tracks ids
*   Response type: application/json
*   Response code: 200 if everything is OK, appropriate error code otherwise.
*   Response: the array of IDs, or an empty array if no tracks have been stored yet.

    [ \<id1\>, \<id2\>, ... ]

### [](#get-apiigcid)GET /api/igc/`<id>`

*   What: returns the meta information about a given track with the provided `<id>`, or NOT FOUND response code with an empty body.
*   Response type: application/json
*   Response code: 200 if everything is OK, appropriate error code otherwise.
*   Response:

```javascript
    {
      "H_date": <date from File Header, H-record>,
      "pilot": <pilot>,
      "glider": <glider>,
      "glider_id": <glider_id>,
      "track_length": <calculated total track length>
    }
```

### [](#get-apiigcidfield)GET /api/igc/`<id>`/`<field>`

*   What: returns the single detailed meta information about a given track with the provided `<id>`, or NOT FOUND response code with an empty body. The response should always be a string, with the exception of the calculated track length, that should be a number.
*   Response type: text/plain
*   Response code: 200 if everything is OK, appropriate error code otherwise.
*   Response
    *   `<pilot>` for `pilot`
    *   `<glider>` for `glider`
    *   `<glider_id>` for `glider_id`
    *   `<calculated total track length>` for `track_length`
    *   `<H_date>` for `H_date`

## [](#resources)Resources

*   [Go IGC library](https://github.com/marni/goigc)

