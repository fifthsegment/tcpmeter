# Synopsis #

Tcpmeter measures the speed of a TCP link between a client and a server. When started in server mode, it listens on an RPC port for incoming requests. In client mode, it provides a web user interface to receive instructions from the user and to display the measurements.


# Usage #

## Starting the server ##

On the far-end server, start _tcpmeter_ in server mode and specify the host's name or IP address and the port that it will listen on for RPC:

`tcpmeter -s -r "`**ipaddr:port**`"`

## Starting the client ##

On your local machine, start _tcpmeter_ in client mode:

`tcpmeter -c`

then open the link that is printed out in a modern browser (HTML5). Using the UI in the browser specify the IP and RPC port address of the far-end server, choose the direction of test (Upload/Download) and the amount of data to transfer and start the measurement.

## Screenshot ##

![http://googledrive.com/host/0B0sQhgOyZZBsZmYwOGI5ZTMtYTkwNC00NTFlLWJlZTgtNzAzOGVhYTEzNGIw/tcpmeter.png](http://googledrive.com/host/0B0sQhgOyZZBsZmYwOGI5ZTMtYTkwNC00NTFlLWJlZTgtNzAzOGVhYTEzNGIw/tcpmeter.png)

## Caveat ##

Despite what the UI shows, UDP and TCP RTT measurements aren't wired in yet.  There's a fair bit of debug output.