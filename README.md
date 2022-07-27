# Backend Technical Test for Scalingo

## Instructions

This application needs a [github personal token](https://github.com/settings/tokens) in order to bypass the API rate limit.
You can fulfill it in a .env file respecting the .env.example template.

## Execution

```
docker-compose up
```

Application will be then running on port `5000`

# API
___
## List repositories

Get the list of the last 100 repositories and extract a sublist based on the request filters.

**URL** : `/repos`

**Method** : `GET`

**Available filters**: owner, topics, language

### Success Response

**Code** : `200 OK`

**Content examples**

GET /repos?owner=jnicklas&language=Ruby

```json
{
    "repositories": [
      {
        "github_id": 248,
        "name": "rorem",
        "owner": "jnicklas",
        "full_name": "jnicklas/rorem",
        "url": "https://github.com/jnicklas/rorem",
        "created_at": "0001-01-01T00:00:00Z",
        "topics": [],
        "languages": {
          "Ruby": 56694
        }
      }
    ]
}
```

### Bad request Response

**Code** : `400 KO`

**Content examples**


```json
{
    "code": 400,
    "message": "Invalid filters"
}
```

### Failure Response

**Code** : `500 KO`

**Content examples**

```json
{
  "code": 500,
  "message": "Failed to list repositories"
}
```
___

## Analyse repositories

Analyses and calculate statistics about the repositories extracted on the /repos route.

**URL** : `/repos`

**Method** : `GET`

**Available filters**: owner, topics, language

### Success Response

**Code** : `200 OK`

**Content examples**

GET /stat?owner=jnicklas&language=Ruby

```json
{
    "statistics": {
        "owners": {
          "jnicklas": 1
        },
        "topics": {},
        "languages": {
          "Ruby": 1
        }
    }
}
```

### Bad request Response

**Code** : `400 KO`

**Content examples**


```json
{
    "code": 400,
    "message": "Invalid filters"
}
```

### Failure Response

**Code** : `500 KO`

**Content examples**

```json
{
  "code": 500,
  "message": "Failed to list repositories"
}
```

