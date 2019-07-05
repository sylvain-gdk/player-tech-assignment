# Player Developer tech assignment

## Problem description

Imagine you want to update the software of thousand of music players that are already in the field. A music player is composed of multiple component, each having its version.
We know that every 15 mins, each player query an API to see if a new version is available and then update itself.

## The assignment

You need to create a production-ready tool that will automate the update of thousand music players.

Your tool will be used by different people using different operating systems. Most commons will be Windows, MacOS and Linux

The input is a .csv file containing, at least, MAC addresses of players to update, always in the first column.

### Example of a .cvs file:
```
mac_addresses, id1, id2, id3
a1:bb:cc:dd:ee:ff, 1, 2, 3
a2:bb:cc:dd:ee:ff, 1, 2, 3
a3:bb:cc:dd:ee:ff, 1, 2, 3
a4:bb:cc:dd:ee:ff, 1, 2, 3
```

### The API to use to update the software version is:

```
PUT /profiles/clientId:{macaddress}
```

#### Request

```
Headers

Content-Type: application/json
x-client-id: required
x-authentication-token: required

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  }
}

```

#### Reponses

##### 200
```
Headers

Content-Type: application/json

Body

{
  "profile": {    
    "applications": [
      {
        "applicationId": "music_app"
        "version": "v1.4.10"
      },
      {
        "applicationId": "diagnostic_app",
        "version": "v1.2.6"
      },
      {
        "applicationId": "settings_app",
        "version": "v1.1.5"
      }
    ]
  }
}
```

##### 401

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 401,
  "error": "Unauthorized",
  "message": "invalid clientId or token supplied"
}
```

##### 404

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 404,
  "error": "Not Found",
  "message": "profile of client 823f3161ae4f4495bf0a90c00a7dfbff does not exist"
}
```

##### 409

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 409,
  "error": "Conflict",
  "message": "child \"profile\" fails because [child \"applications\" fails because [\"applications\" is required]]"
}
```

##### 500

```
Headers

Content-Type: application/json

Body

{
  "statusCode": 500,
  "error": "Internal Server Error",
  "message": "An internal server error occurred"
}
```

## What we mean by production-ready:
- developer documentation (how to build, how to run tests)
- user documentation (how to use the tool, can be "embedded" in the tool itself or in a document)
- unit tests

## Important notes

You can use the language/technology of your choice and it is highly suggested to provide a document (plain text or markdown, nothing fancy) to explain your decisions

## Submission

Submit your assignment using a public github / gitlab / bitbucket repository (please don't use "touchtunes" in repo name or description) or a zip archive.