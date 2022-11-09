# metadata #

## Overview ##

`metadata` is an experimental hack to mess with extracting data from the 
[spinitron.com](https://spinitron.com/)
[metadata push](https://forum.spinitron.com/t/metadata-push-guide/144#push-channels)
mechanism.

`metadata` listens on a TCP port (default: 52341) for connection requests
for _metadata push_ requests in XML format from
[spinitron.com](https://spinitron.com/) (or, really anywhere).

`metadata` extracts the relevant data from the XML and writes it to a MariaDB database.

**WARNING**

This code employs **_ZERO_** encryption, authentication, authorization or security.
Use at your own peril.

## Example Spinitron Channel Configuration ##

See the [template spec](https://forum.spinitron.com/t/metadata-push-guide/144#template-specification)
for complete details.

```
tcp://metadata.example.com:52341
<?xml version="1.0" encoding="utf-8" ?>
<nowplaying>
    <time>%now%</time>
    <artist>%an%</artist>
    <title>%sn%</title>
    <album>%dn%</album>
</nowplaying>
```
## Database Schema ##

```
CREATE TABLE `m` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `time` datetime NOT NULL,
  `artist` varchar(128) DEFAULT NULL,
  `title` varchar(128) DEFAULT NULL,
  `album` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5190 DEFAULT CHARSET=utf8mb4
```
