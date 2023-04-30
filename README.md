# P2P Network Analysis

This is a Go web application that monitors connection metrics on the Bitcoin and Lightning networks, and saves the data to a SQLite database using GORM.

## Prerequisites
- Go version 1.16 or later
- btcd: a full-node implementation of the Bitcoin protocol
- lnd: a Lightning Network Daemon implementation
- SQLite3: a lightweight database engine

### Installation
1. Clone the repository to your local machine:
```bash
 git clone https://github.com/NonsoAmadi10/Heimdall

 cd Heimdall
```
2. Install the dependencies:
```bash
go mod download

```
3. Build the binary:
```bash
go build

```

### Usage 
1. Start a btcd node

```bash
btcd --testnet --rpcuser=rpcuser --rpcpass=rpcpass
```

2. Start a lightning node
```bash
lnd --bitcoin.testnet --bitcoin.rpcuser=rpcuser --bitcoin.rpcpass=rpcpass
```
3. Run the application:
```bash
./heimdall
```
The application will run in the background and collect connection metrics every minute, and save them to the metrics table in the metrics.db database file.

To view the metrics, you can open your postman on the following endpoints:

`http://localhost:1700/node-info` - Fetches node information for both bitcoin and lightning
`http://localhost:1700/conn-metrics` - fetches metrics based on connection between your node and other peers as well as network bandwidth

Alternatively, there is a frontend metrics that gives visualization on the metrics generated. To view it simply:

1. Enter the frontend directory and install the node dependencies:

```bash
cd dashboard
yarn install 
```

2. Start the dev server:
```bash
yarn dev
```
3. Open your browser on `http://localhost:3000`


PS: I would be looking to add network hop messages, transaction blocks metrics and hashrate metrics

## Contributing 
Feel free to submit issues or pull requests if you have suggestions for improvement or find any bugs in the tool.

## License 

This project is licensed under the MIT License - see the [License](https://github.com/NonsoAmadi10/Heimdall/blob/main/LICENSE) file for details.