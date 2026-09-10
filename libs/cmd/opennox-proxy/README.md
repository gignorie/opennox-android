# OpenNox net proxy

This tool helps debug networking problems by sitting between Nox/OpenNox clients and the server.

![Diagram](./diagram.png)

## Flow from client
1. Packets from clients are accepted on the server proxy port.
2. Proxy then allocates a unique proxy port for each client (since Nox uses ip+port for client id).
3. Client packets are then sent from client proxy port to the real server.

## Flow from server
1. Packets are received on unique client proxy port.
2. They are then sent from server proxy port to the client.

## How to run

```shell
go run ./cmd/opennox-proxy --server=<server-ip>:18590 --host=0.0.0.0:18600 --file=network.jsonl
```

This will run a proxy on port `18600`, which is not standard for Nox, thus server discovery will not find it.

You must add the proxy address to `game_ip.txt` file in Nox directory (as `127.0.0.1:18600`).

After this, you should see a new server in the server list, you will recognize it by "Proxy:" prefix.

While you are connected to this server, all network messages are logged to `network.jsonl`.
Once you're done with testing, disconnect from the proxy, and close it.

Now you can run `go run ./cmd/opennox-packet-decode` that will decode known network messages
in `network.jsonl` and will write them to `network-dec.jsonl`, which can then be inspected.