  # Go-Based Write Ahead Logging (WAL) Engine
  
  # 1. Introduction
  
  <img width="1417" height="629" alt="High Level Architecture" src="https://github.com/user-attachments/assets/e0e86f18-f866-4dd9-85d4-3d7ff6db08f2" />
  
  
  ## 1.1 Purpose
  
  The purpose of this system is to provide a durable, high-performance, crash-safe Write Ahead Log (WAL) engine implemented in Go.
  
  The system shall support:
  
  * CLI-based log ingestion
  * gRPC-based log ingestion
  * Log persistence
  * Crash recovery
  * Log replay
  * Checkpointing
  * Integrity validation
  * Replication-ready architecture
  
  The WAL acts as the system of record for downstream services such as:
  
  * Ledger engines
  * Transaction processors
  * Reconciliation systems
  * Simulation engines
  * Analytics pipelines
  
  ---
  
  ## 1.2 Scope
  
  The WAL Engine shall:
  
  * Accept log entries
  * Persist entries to disk
  * Guarantee durability
  * Recover after crashes
  * Replay logs
  * Validate data integrity
  * Support concurrent ingestion
  
  The WAL Engine shall NOT:
  
  * Execute business logic
  * Maintain account balances
  * Perform reconciliation
  * Process transactions
  
  These responsibilities belong to downstream systems.
  
  ---
  
  # 2. System Overview
  
  ## High-Level Architecture
  
  ```text
  +----------------+
  | CLI Client     |
  +----------------+
          |
          |
          v
  
  +----------------+
  | gRPC Client    |
  +----------------+
          |
          |
          v
  
  +----------------------+
  | Ingestion Layer      |
  +----------------------+
          |
          v
  
  +----------------------+
  | WAL Writer           |
  +----------------------+
          |
          v
  
  +----------------------+
  | Segment Manager      |
  +----------------------+
          |
          v
  
  +----------------------+
  | Disk Storage         |
  +----------------------+
  
          |
          v
  
  +----------------------+
  | Replay Engine        |
  +----------------------+
  ```
  
  ---
  
  # 3. Functional Requirements
  
  # FR-1 Log Ingestion
  
  ## FR-1.1 CLI Log Submission
  
  The system shall accept logs via CLI.
  
  Example:
  
  ```bash
  wal append \
  --type transfer \
  --payload '{"from":"A","to":"B","amount":100}'
  ```
  
  Expected Response:
  
  ```bash
  Entry persisted.
  Sequence: 12345
  ```
  
  ---
  
  ## FR-1.2 gRPC Log Submission
  
  The system shall expose a gRPC service.
  
  Example:
  
  ```protobuf
  service WalService {
    rpc Append(AppendRequest)
        returns (AppendResponse);
  }
  ```
  
  ---
  
  ## FR-1.3 Batch Submission
  
  The system shall support:
  
  ```protobuf
  rpc AppendBatch(...)
  ```
  
  to reduce disk sync overhead.
  
  ---
  
  # FR-2 Log Storage
  
  Every log record shall be written sequentially.
  
  No random updates allowed.
  
  Supported operations:
  
  ```text
  Append
  Read
  Replay
  Checkpoint
  ```
  
  ---
  
  # FR-3 Log Record Format
  
  Each WAL entry shall contain:
  
  ```rust
  struct WalRecord {
      magic: u32,
      version: u16,
  
      sequence: u64,
  
      timestamp: u64,
  
      payload_type: u16,
  
      payload_size: u32,
  
      crc64: u64,
  
      payload: Vec<u8>,
  }
  ```
  
  ---
  
  ## Field Description
  
  | Field        | Purpose               |
  | ------------ | --------------------- |
  | magic        | Corruption detection  |
  | version      | Upgrade compatibility |
  | sequence     | Ordering              |
  | timestamp    | Audit                 |
  | payload_type | Event classification  |
  | payload_size | Parsing               |
  | crc64        | Integrity             |
  | payload      | Actual data           |
  
  ---
  
  # FR-4 Segment Management
  
  The WAL shall be divided into segments.
  
  Example:
  
  ```text
  wal/
  ├── segment_000001.wal
  ├── segment_000002.wal
  ├── segment_000003.wal
  ```
  
  ---
  
  ## FR-4.1 Segment Rotation
  
  A new segment shall be created when:
  
  ```text
  segment_size >= configured_limit
  ```
  
  Default:
  
  ```text
  128 MB
  ```
  
  Configurable:
  
  ```text
  64 MB
  128 MB
  256 MB
  1 GB
  ```
  
  ---
  
  # FR-5 Durability
  
  Before acknowledging success:
  
  ```text
  Append
  → Write
  → Flush
  → Sync
  → ACK
  ```
  
  The system shall guarantee:
  
  ```text
  No acknowledged write is lost.
  ```
  
  ---
  
  # FR-6 Group Commit
  
  The system shall support batching.
  
  Example:
  
  ```text
  Tx1
  Tx2
  Tx3
  Tx4
  
  Single fsync()
  ```
  
  Benefits:
  
  * Reduced disk operations
  * Higher throughput
  
  ---
  
  # FR-7 Recovery
  
  Upon startup:
  
  ```text
  Scan Segments
  Validate CRC
  Replay Records
  Restore State
  ```
  
  ---
  
  ## FR-7.1 Torn Write Detection
  
  The system shall detect:
  
  ```text
  Partial Record
  Corrupt Record
  Incomplete Write
  ```
  
  using:
  
  * CRC64
  * Length validation
  * Magic number validation
  
  ---
  
  # FR-8 Replay Engine
  
  CLI:
  
  ```bash
  wal replay
  ```
  
  or
  
  ```bash
  wal replay --from-sequence 100000
  ```
  
  The engine shall stream entries in sequence order.
  
  ---
  
  # FR-9 Checkpointing
  
  The system shall support snapshots.
  
  Directory:
  
  ```text
  wal/
  ├── checkpoints/
  │   └── checkpoint_500000.chk
  ```
  
  Checkpoint contains:
  
  ```text
  Last Sequence Number
  Metadata
  State Snapshot
  ```
  
  ---
  
  # FR-10 Search & Inspection
  
  CLI:
  
  ```bash
  wal inspect 12345
  ```
  
  returns:
  
  ```json
  {
    "sequence":12345,
    "timestamp":"...",
    "payload":"..."
  }
  ```
  
  ---
  
  # FR-11 Metrics
  
  The system shall expose:
  
  ```text
  Append Rate
  Replay Rate
  Segment Count
  Disk Usage
  Recovery Time
  Sync Latency
  Queue Depth
  ```
  
  through Prometheus.
  
  ---
  
  # FR-12 Configuration
  
  Example:
  
  ```yaml
  wal:
    path: "./wal"
  
    segment_size_mb: 128
  
    sync_mode: immediate
  
    checkpoint_interval: 100000
  
    compression: false
  
    grpc_port: 50051
  ```
  
  ---
  
  # 4. Non-Functional Requirements
  
  # NFR-1 Performance
  
  Target:
  
  ```text
  Single Record Write:
  < 2 ms
  ```
  
  ---
  
  Throughput:
  
  ```text
  50,000+ appends/sec
  ```
  
  with batching enabled.
  
  ---
  
  Recovery:
  
  ```text
  10 GB WAL
  
  Recovery Time:
  < 30 seconds
  ```
  
  Target only.
  
  ---
  
  # NFR-2 Reliability
  
  System shall survive:
  
  * Process crash
  * Machine reboot
  * Power loss
  * Partial writes
  
  without corruption of committed records.
  
  ---
  
  # NFR-3 Scalability
  
  Support:
  
  ```text
  Millions of records
  Multiple segment files
  Hundreds of GB WAL size
  ```
  
  ---
  
  # NFR-4 Availability
  
  Service restart:
  
  ```text
  < 10 seconds
  ```
  
  for normal workloads.
  
  ---
  
  # NFR-5 Security
  
  gRPC endpoints shall support:
  
  * TLS
  * mTLS (optional)
  
  Authentication:
  
  ```text
  API Key
  JWT
  mTLS
  ```
  
  ---
  
  # 5. CLI Specification
  
  ## Append
  
  ```bash
  wal append \
  --type transfer \
  --payload file.json
  ```
  
  ---
  
  ## Replay
  
  ```bash
  wal replay
  ```
  
  ---
  
  ## Inspect
  
  ```bash
  wal inspect 1000
  ```
  
  ---
  
  ## Metrics
  
  ```bash
  wal metrics
  ```
  
  ---
  
  ## Checkpoint
  
  ```bash
  wal checkpoint
  ```
  
  ---
  
  ## Verify
  
  ```bash
  wal verify
  ```
  
  Verifies:
  
  * CRC
  * Segment integrity
  * Record ordering
  
  ---
  
  # 6. gRPC API
  
  ```protobuf
  service WalService {
  
      rpc Append(AppendRequest)
          returns (AppendResponse);
  
      rpc AppendBatch(AppendBatchRequest)
          returns (AppendBatchResponse);
  
      rpc Replay(ReplayRequest)
          returns (stream WalRecord);
  
      rpc Health(HealthRequest)
          returns (HealthResponse);
  
      rpc Metrics(MetricsRequest)
          returns (MetricsResponse);
  }
  ```
  
  ---
  
