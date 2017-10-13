# Asterisk API Manager

<!-- TOC -->

- [Asterisk API Manager](#asterisk-api-manager)
    - [Features (TO-DO list)](#features-to-do-list)
    - [Usage](#usage)
    - [Authentication](#authentication)
        - [POST /login](#post-login)
        - [HTTP Header using token to authorize](#http-header-using-token-to-authorize)
    - [Users](#users)
        - [POST /{PREFIX}/users/new](#post-prefixusersnew)
        - [GET /{PREFIX}/users/{userid}](#get-prefixusersuserid)
        - [PUT /{PREFIX}/users/{userid}](#put-prefixusersuserid)
        - [DELETE /{PREFIX}/users/{userid}](#delete-prefixusersuserid)

<!-- /TOC -->

## Features (TO-DO list)

- [x] JWT access tokens
- [x] CURD (Static) endpoints for PJSIP stack
- [ ] Upload and convert audio
- [x] Interface with AMI
- [x] Interface with ARI

## Usage

**Example**

```bash
\$ cd run/
\$ ./gastwrap -prefix=hapulico -listen=127.0.0.1:8080
```

## Authentication

### POST /login

Login and get token.

- Content-Type: "application/json"
- Accept: "application/json"

```json
{
    "username": "user",
    "password": "secret"
}
```

---

**Example response**

- Status: 200
- Content-Type: "application/json"

```json
{
    "token": "this*is*a*T0k3n"
}
```

---

### HTTP Header using token to authorize

```HTTP
"Authorization": "Bearer this*is*a*T0k3n"
```

## Users

### POST /{PREFIX}/users/new

Creat a new user.

- Content-Type: "application/json"
- Accept: "application/json"

```json
{
    "userid": "newuserid",
    "callerid": "newcallerid",
    "password": "newsecret",
    "context": "default",
    "phonemodel": "grandstream",
    "phonemac": "000000000000"
}
```

---

**Example response**

- Status: 200
- Content-Type: "application/json"

```json
{
    "status": "Success"
}
```

```json
{
    "status": "Fail: An error occurred"
}
```

### GET /{PREFIX}/users/{userid}

Get user infos

---

**Example response**

- Status: 200
- Content-Type: "application/json"

```json
{
    "userid": "newuserid",
    "callerid": "newcallerid",
    "password": "newsecret",
    "context": "default",
    "phonemodel": "grandstream",
    "phonemac": "000000000000"
}
```

OR

- Status: 400

### PUT /{PREFIX}/users/{userid}

Update user infos

- Content-Type: "application/json"
- Accept: "application/json"

```json
{
    "userid": "",
    "callerid": "",
    "password": "newsecret",
    "context": "",
    "phonemodel": "",
    "phonemac": ""
}
```

---

**Example response**

- Status: 200
- Content-Type: "application/json"

```json
{
    "status": "Success"
}
```

### DELETE /{PREFIX}/users/{userid}

Delete an user

---

**Example response**

- Status: 200
- Content-Type: "application/json"

```json
{
    "status": "Success"
}
```
