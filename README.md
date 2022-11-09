# metadata #

## Overview ##

metadata listen on a TCP port (default: XX) for connection requests,
and listens for ["metadata push"](https://forum.spinitron.com/t/metadata-push-guide/144#push-channels)
requests in XML format from
[spinitron.com](https://spinitron.com/).

metadata extracts the relevant data from the XML and writes it to a MariaDB database.

## Example XML ##

```
<?xml version="1.0" encoding="utf-8" ?>
<nowplaying>
    <artist>%an%</artist>
    <title>%sn%</title>
    <album>%dn%</album>
</nowplaying>
```
