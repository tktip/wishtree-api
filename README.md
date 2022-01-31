#wishtree-api

####Introduction
Wishtree-api is the API for the wishtree service, receiving REST calls from the client service and fetching/storing data in the database. 

####How to build
**Build as a runable file**
We use Make to build our projects. You can define what system to build for by configuring the `GOOS` environment variable.

These commands will build either a runable linux or windows file in the /bin/amd64 folder:

`GOOS=windows make clean build`
`GOOS=linux make clean build`

**Build a docker container**
First you have to define the docker registry you are going to use in `envfile`. Replace the `REGISTRY` variable with your selection. This command will build the container and push it to your docker registry:

`GOOS=linux make clean build container push`

####How to run
**As an executable**
If you want to run it as a executable (e.g. a Windows service) you will need to configure the correct environment variable. When you start the application set the `CONFIG` environment to `file::<location>` for linux or run it as a argument for windows

Windows example: `set "CONFIG=file::Z://folder//cfg.yml" & wishtree-api.exe`

Linux example: `CONFIG=file::../folder/cfg.yml ./wishtree-api`

**In a Docker container**
You have to first mount the cfg file into the docker container, and then set the config variable to point to that location before running the service/container.


####Configuration file
The configuration variables for the API are located in `cfg.yml`.
| Var | Description |
|---|---|
| `adminUsers` | A list of key-value pairs for users and passwords. A minimum of one, with no upper limit. Keys: `username`, `password` |
| `db` | The database connection configuration, key-value pairs.<br>  - `url`: How to connect to the database. Typically `localhost` or an IP address.<br> - `username`: The username of the database user<br> - `password`: The password of the database user <br> - `name`: The name of the database|
| `tokenPrivatekeyPath` | The relative (or full, if not running as a Docker container) path to the file containing the *private* key used to verify the JWT token. Recommendation: put these key files in the same folder as `cfg.yml`. |
| `tokenPublickeyPath` | The relative/full path to the *public* key used to verify the JWT token. See ^.|
| `tokenConfig` | The configuration for the JWT authorization token, used to maintain and verify admin sessions. Key-value pairs.<br><br> - `tokenName`: The name of the token cookie, can be anything <br> - `durationHours`: The duration of an admin's login session <br> - `domain`: The domain of the cookie holding the JWT token. If the wishtree runs on its own domain, or on *domain.xx/yyy*, an empty string is sufficient. <br> - `path`: If set, the cookie will only be stored on the set path in the URL. If the wishtree exists on *domain.xx/yyy*, the path can be set to *yyy*. <br> - `secure`: If true, the cookie (and, therefore, logging in as an admin) will only work on HTTPS connections. Recommended. <br> - `httpOnly`: If true, restricts Javascript's access to the cookie in the browser. Highly recommended. <br><br> For more detailed information on these cookie settings, see https://developer.mozilla.org/en-US/docs/Web/HTTP/Cookies. |

####JWT keys
Two RS256 keys need to be generated and placed into their respective files. Such files should begin with the line `-----BEGIN <PUBLIC/PRIVATE> KEY-----` and end with the line `-----END <PUBLIC/PRIVATE> KEY-----`. The private key should always be kept secret.

The `openssl` command line tool can be used. Example that will generate two files, one for each key:
>`openssl genrsa -out jwtRS256_private.key 4096`
>`openssl rsa -in jwtRS256_private.key -pubout > jwtRS256_public.key.pub`

Alternatively, there are online tools to help generate such key pairs.
