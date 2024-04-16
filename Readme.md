## Context

Putting in practice some blockchain concepts:

* Blockchain Structure:

The blockchain structure consists of blocks linked together through hashes. Each block contains a header and a list of transactions.
The proto.Block and proto.Header structures define the block and header formats respectively.

* Cryptographic Operations:

Cryptographic operations are implemented using the crypto package, including key generation, signing, verification, and address generation.
It utilizes the Ed25519 elliptic curve digital signature algorithm for generating key pairs and signing transactions.

* Networking:

The networking functionality is implemented using gRPC, a high-performance RPC (Remote Procedure Call) framework.
Nodes communicate with each other through gRPC calls. Each node exposes an RPC service with methods such as HandleTransaction and HandShake.
Peers exchange information about the blockchain, transactions, and network topology through these RPC calls.

* Node Management:

The node package handles the management and operation of individual nodes in the blockchain network.
Each node maintains a list of peers it is connected to, manages its own mempool for storing pending transactions, and handles incoming transactions and block validation.

* Consensus and Validation:

Consensus mechanisms are not fully implemented in this code. However, basic block validation is performed to ensure the integrity of the blockchain.
Blocks are validated to ensure that transactions are properly signed and that the previous block's hash matches the hash specified in the current block.

* Genesis Block:

The CreateGenesisBlock function initializes the blockchain with a genesis block. This is the first block in the blockchain and typically has special properties compared to regular blocks.

* Mempool:

The Mempool structure is used to store pending transactions that have been received by the node but have not yet been included in a block.
Transactions are added to the mempool upon receipt and removed when they are included in a block.

* Bootstrapping:

Nodes bootstrap the network by connecting to known peers during startup. This helps in establishing initial connections and discovering other nodes in the network.
Peers exchange version information during the handshake process, including the node's version, height, and peer list.