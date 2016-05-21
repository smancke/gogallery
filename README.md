
Golang based web image gallery
================================

Gogallery is an image gallery written in golang and jquery.
[It can be seen here in action](https://bilderbuch-stoff.de/gallery/ui/pub/index_integration.html)

The images are stored in the filesystem and indexed within a sqlite database file.

Functionality
------------------
* Image upload of logged in users
* Gallery view of the images
* The user can change it's data
* A user can see and delete it's own images
* Image scaling, and auto rotation
* Facebook integration for the images (like and share)
* Chunk-Uploading

Building the application
---------------------------
To download and build the gogallery using golang, simply execute

```shell
go get github.com/smancke/gogallery
```

The application used the libvips tools for scaling of the images.
To install them on e.g. on ubuntu:

```shell
sudo apt-get install libvips-tools --no-install-recommends
```

Running with docker
----------------------
We provide a docker container based on the master of this repository
```shell
docker run -e testOverwriteUsername=demo \
      -v /opt/gallery/data:/var/lib/gallery \
      -v /tmp:/tmp \
      -p 5005:5005 \
      smancke/gogallery
```

Tip: You can start the container with `--restart=always` on your server, to have it running after restarts of the machine.

Configuration
------------------
The gogallery can be configured by environment variables with the following defaults:

```shell
# The address to listen (host:port)
address=:5005

# The directory for storing data
galleryDir=/tmp/gallery

# The directory for the static html templates
htmlDir=./html

# The name of the cookie to verify logins
cookieName=okmsdc

# The session secret to decrypt cookies with
session_secret=XYcretsecretsecretsecretsecretse

# The maximum lifetime for sessions
sessionLifetimeMinutes=180

# Disable the authentication and
# set the username of the logged in user for testing
testOverwriteUsername=
```

Authentication
--------------------
For access control, the gogallery looks for a cookie with the name configured by `cookieName`.
This cookie is interpreted as a base64encoded, pkcs5Padded and AES encrypted by the key configured in `session_secret`.
The payload of the cookie is a json object with the following structure:

```json
{
   "groups" : [
      "admin"
   ],
   "lastSeen" : 1463818233690,
   "userId" : "07323464-76da-44f9-8eec-123451fb72e4",
   "displayName" : "Sebastian Mancke",
   "userName" : "s.mancke@example.com"
}
```

If someone is interested in that, it would be easy to integrate further access controll methods for simple integration.
E.g. forwarding a cookie from the user to an endpoint configured REST endpoint, which returns the above json.
Please file a github issue, if you are interested in something like that.
