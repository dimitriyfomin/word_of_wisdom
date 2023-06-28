# Service description
It's example of TCP service to provide Anti-DDOS protection with PoW.
Client requests challange from server and calculate token by the following rules: hash_sha256(challenge+token) starts with 000 bytes.
If passed token is correct, service will process request to get random quote of wisdom.

## Explanation of the choice of the POW algorithm
In this implementation, the client generates random data and computes the SHA-256 hash result by appending a counter in a loop until a hash result with the required prefix (number of leading zeros) is found. This requires a certain amount of computation time which can effectively limit the rate at which requests can be made to the server.

Here PoW algorithm is simple Hashcash, it is high enough to provide sufficient protection against DDoS attacks, but not too high that it causes significant delays or inconvenience for legitimate users.

## Build&Run
### Prepare server docker
```
docker build -t word-of-wisdom-server -f Dockerfile.server .
```

### Prepare client docker
```
docker build -t word-of-wisdom-client -f Dockerfile.client .
```

### Run server
```
docker run -p 8080:8080 word-of-wisdom-server
```

### Run client
```
docker run --network="host" word-of-wisdom-client
```

## Contract
Client sends message in format:
```
<fullCommand> <token:optional>
```

When command supports versions:
```
<fullCommand> = <version>.<command>
```
Versions are implemented in hub-spike approach: each version must call hub methods. So, versions implement mapping under hub.

### Client commands
#### CHG
Client requests "challenge" from server.
```
v1.CHG
```
Response:
```
OK <challenge>
```

#### CHGT
Client passes "challenge" token to server.
```
v1.CHGT <token>
```
Response:
```
OK
```

#### QTR
Client requests random qoute from server.
```
v1.QTR
```
```
OK <quote>
```

### Server commands
#### OK
Server responses successful message.

```
OK
OK <payload>
```

#### ERR
Server responses failed error and reason.

```
ERR <error-message>
```
