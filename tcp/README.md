# TCP Toy
* `nc localhost 8080` or `netcat localhost 8080` to send and receive messages
* The Local Addr remains same, but remote address varies for every conn
* conn have `SetReadDeadline` that must be set before every read for timeout. Similarly for write deadline
