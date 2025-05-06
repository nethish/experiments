# Circuit Breaker

* The server toggles between online and offline every 5 seconds
* If the client experiences 3 consequetive failures, then open the circuit and don't allow any more requests through for next 3 seconds
* Allow requests after 3 seconds, and if there are two successful requests, close the circuit again
