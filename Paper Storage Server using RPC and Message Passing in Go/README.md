
# Paper Storage Server

A scalable and efficient distributed system tailored for storing, retrieving, and managing academic papers. This system utilizes RPC (Remote Procedure Calls) for seamless client-server interactions and RabbitMQ to handle real-time message-passing notifications.

## Features

- **Upload Papers**: Easily add academic papers along with relevant metadata like authors and titles.
- **Retrieve Papers**: Access detailed metadata or the full content of stored papers using unique identifiers.
- **Browse Papers**: View a comprehensive list of all papers available in the system.
- **Real-Time Notifications**: Receive instant alerts when new papers are added, powered by RabbitMQ.
- **Distributed Architecture**: Built on a robust client-server model using gRPC and message-passing mechanisms to ensure high performance and scalability.

---

## Installation

### Prerequisites

Before setting up the Paper Storage Server, ensure the following are installed on your system:

1. **Go** (version 1.16 or higher)
2. **RabbitMQ** (installed and actively running)
3. **Protocol Buffers Compiler** (`protoc`)

### Setup Steps

1. **Clone the Repository**
   ```bash
   git clone https://github.com/your-username/paper-storage-server.git
   cd paper-storage-server
   ```

2. **Install Dependencies**
   ```bash
   go mod tidy
   ```

3. **Generate gRPC Code** (if necessary)
   ```bash
   protoc --go_out=. --go-grpc_out=. paperpb/paper.proto
   ```

---

## Usage

### Starting the Server

1. Navigate to the server directory:
   ```bash
   cd cmd/paperserver
   ```

2. Launch the server:
   ```bash
   go run main.go
   ```

### Interacting with the Client

1. Move to the client directory:
   ```bash
   cd cmd/paperclient
   ```

2. Run the client to connect with the server:
   
   **Initialize the Client**
   ```bash
   go run main.go <server-address>
   ```

   **Add a New Paper**
   ```bash
   add <server-address> "<author>" "<title>" <file-path>
   ```
   *Example:*
   ```bash
   add localhost:50051 "Jane Smith" "Distributed Computing" research_paper.pdf
   ```

   **List All Papers**
   ```bash
   list <server-address>
   ```
   *Example:*
   ```bash
   list localhost:50051
   ```

   **Retrieve Paper Details**
   ```bash
   detail <server-address> <paper-id>
   ```
   *Example:*
   ```bash
   detail localhost:50051 1
   ```

   **Fetch Paper Content**
   ```bash
   fetch <server-address> <paper-id>
   ```
   *Example:*
   ```bash
   fetch localhost:50051 1
   ```

---

## Notification System

The Paper Storage Server employs RabbitMQ to deliver real-time notifications to subscribed clients whenever a new paper is added. This asynchronous messaging ensures that users are immediately informed of updates without the need for continuous polling.

---

## Author

- **Fuad Mohammed**  
  **ID:** UGR/6052/14

```