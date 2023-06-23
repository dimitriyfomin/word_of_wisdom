# Service description
It's example of TCP service to provide Anti-DDOS protection with PoW.
On connection server sends "challenge" to client. Client must answer token by the following rules: hash_sha256(challenge+token) starts with 000 bytes.
If token is correct, service will return random quote of wisdom.

# Explanation of the choice of the POW algorithm
In this implementation, the client generates random data and computes the SHA-256 hash result by appending a counter in a loop until a hash result with the required prefix (number of leading zeros) is found. This requires a certain amount of computation time which can effectively limit the rate at which requests can be made to the server.

Here PoW algorithm is simple Hashcash, it is high enough to provide sufficient protection against DDoS attacks, but not too high that it causes significant delays or inconvenience for legitimate users.

# Prepare server docker
```
docker build -t word-of-wisdom-server -f Dockerfile.server .
```

# Run server
```
docker run -p 8080:8080 word-of-wisdom-server
```

# Prepare client docker
```
docker build -t word-of-wisdom-client -f Dockerfile.client .
```

# Run client
```
docker run --network="host" word-of-wisdom-client
```
