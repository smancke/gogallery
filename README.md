
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
To download and build the gogallery using go, simply execute

```shell
go get github.com/smancke/gogallery
```

The application used the libvips tools for scaling of the images.
To install them on e.g. on ubuntu:

```shell
sudo apt-get install libvips-tools 
```

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
session_secret=secretsecretsecretsecretsecretse

# The maximum lifetime for sessions
sessionLifetimeMinutes=180
```
