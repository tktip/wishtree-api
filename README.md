# wishtree-api

## Introduction
Wishtree-api is the API for the wishtree service, receiving REST calls from the client service and fetching/storing data in the database. 

## How to build
**Build as a runable file**

We use Make to build our projects. You can define what system to build for by configuring the `GOOS` environment variable.

These commands will build either a runable linux or windows file in the /bin/amd64 folder:

`GOOS=windows make clean build`
`GOOS=linux make clean build`

**Build a docker container**

First you have to define the docker registry you are going to use in `envfile`. Replace the `REGISTRY` variable with your selection. This command will build the container and push it to your docker registry:

`GOOS=linux make clean build container push`

## How to run
**As an executable**

If you want to run it as a executable (e.g. a Windows service) you will need to configure the correct environment variable. When you start the application set the `CONFIG` environment to `file::<location>` for linux or run it as a argument for windows

Windows example: `set "CONFIG=file::Z://folder//cfg.yml" & wishtree-api.exe`

Linux example: `CONFIG=file::../folder/cfg.yml ./wishtree-api`

**In a Docker container**

You have to first mount the cfg file into the docker container, and then set the config variable to point to that location before running the service/container.


## Configuration file
The configuration variables for the API are located in `cfg.yml`.

| Var                   | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                 |
| --------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `adminUsers`          | A list of key-value pairs for users and passwords. A minimum of one, with no upper limit. Keys: `username`, `password`                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                      |
| `db`                  | The database connection configuration, key-value pairs.<br>  - `url`: How to connect to the database. Typically `localhost` or an IP address.<br> - `username`: The username of the database user<br> - `password`: The password of the database user <br> - `name`: The name of the database                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `tokenPrivatekeyPath` | The relative (or full, if not running as a Docker container) path to the file containing the *private* key used to verify the JWT token. Recommendation: put these key files in the same folder as `cfg.yml`.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                               |
| `tokenPublickeyPath`  | The relative/full path to the *public* key used to verify the JWT token. See ^.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| `tokenConfig`         | The configuration for the JWT authorization token, used to maintain and verify admin sessions. Key-value pairs.<br><br> - `tokenName`: The name of the token cookie, can be anything <br> - `durationHours`: The duration of an admin's login session <br> - `domain`: The domain of the cookie holding the JWT token. If the wishtree runs on its own domain, or on *domain.xx/yyy*, an empty string is sufficient. <br> - `path`: If set, the cookie will only be stored on the set path in the URL. If the wishtree exists on *domain.xx/yyy*, the path can be set to *yyy*. <br> - `secure`: If true, the cookie (and, therefore, logging in as an admin) will only work on HTTPS connections. Recommended. <br> - `httpOnly`: If true, restricts Javascript's access to the cookie in the browser. Highly recommended. <br><br> For more detailed information on these cookie settings, see https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies. |

## JWT keys
Two RS256 keys need to be generated and placed into their respective files. Such files should begin with the line `-----BEGIN <PUBLIC/PRIVATE> KEY-----` and end with the line `-----END <PUBLIC/PRIVATE> KEY-----`. The private key should always be kept secret.

The `openssl` command line tool can be used. Example that will generate two files, one for each key:
>`openssl genrsa -out jwtRS256_private.key 4096`
>`openssl rsa -in jwtRS256_private.key -pubout > jwtRS256_public.key.pub`

Alternatively, there are online tools to help generate such key pairs.

## Database
This project has been built using a Microsoft SQL Server database. With a few changes to `/internal/db/connector.go`, as well as to the generate sql scripts, the project should work with other common SQL databases like PostgreSQL or MySql. There are scripts to generate the tables in the `db_setup/generate_db.sql` file. You may choose a name for the database, which must be reflected in `cfg.yml`. The generate scripts might not work with other databases - either change the necessary data types, or create the tables manually using information gathered from the scripts

You also need to insert 1500 empty wishes into the databas. A script that does this is in the `db_setup/generate_data.sql` file. This script also sets up the categories.

There are initially 1500 empty wishes, but with predefined `x` and `y` coordinates. As wishes are taken, these will have their other database fields populated. When the number of wishes reach 1400, each new wish will cause the oldest one to be *archived*. A copy will be created with `isArchived` true, and the new wish will overwrite the old one it replaces. Archived wishes do not show up on the tree, but can still be found in the list. This way, there are always 100 free wishes, so that users may always find free leaves in the tree.

When a wish is deleted by an admin, it is reset to its initial state with only `x` and `y` values.

### Tables
While knowledge of the database tables isn't necessary to run the project, it might be helpful in case modifications are needed.

**wish**
- `id`: An ID that never changes
- `x`, `y`: The coordinates in the tree
- `text`, `author`, `zipCode`, `createdAt`: User inputs and auto-set timestamp
- `category_id`: The color chosen by the user, Foreign key
- `isArchived`: Described above

**category**
- `id`: An ID that never changes
- `name`: A descriptive but otherwise nonfunctional name. Will never show up in the front-end (eg. `purple`)
- `description`: The hex value of the color, `#` included (eg. `#ba49a5`)

**tree_status**
- `isOpen`: If 1, new wishes can be added. If 0, the tree is "closed". This is changed from the admin panel.
