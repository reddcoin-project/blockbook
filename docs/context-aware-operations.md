# Context-Aware Database Operations

This document explains Blockbook's context-aware database operations that provide granular timeout control and resource protection.

## Overview

Context-aware operations allow database queries to be cancelled mid-operation when request timeouts occur, preventing long-running operations from consuming server resources indefinitely.

## How Context Cancellation Works

### 1. Request Flow with Context

```
HTTP Request → Timeout Wrapper → API Worker → Database Operation
     ↓              ↓              ↓              ↓
  r.Context()   WithTimeout()  ctx.Context()  Iterator Loop
                  (45s)                        ctx.Err() checks
```

### 2. Context Timeout Mechanism

```go
// When request starts
ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
defer cancel()

// Go runtime creates background timer
go func() {
    timer := time.NewTimer(45*time.Second)
    <-timer.C
    close(ctx.Done())  // Signal timeout occurred
    ctx.err = context.DeadlineExceeded
}()
```

### 3. Database Cancellation Points

```go
// In RocksDB iterator - checks every iteration
for it.Seek(startKey); it.Valid(); it.Next() {
    if ctx.Err() != nil {  // ← Instant check (nanoseconds)
        return ctx.Err()   // Returns context.DeadlineExceeded
    }

    // Process database record (potentially slow)
    processRecord(it.Key(), it.Value())
}
```

## Context-Aware Methods

### Database Layer (db/rocksdb.go)

#### GetAddrDescTransactionsContext
```go
func (d *RocksDB) GetAddrDescTransactionsContext(
    ctx context.Context,
    addrDesc bchain.AddressDescriptor,
    lower, higher uint32,
    fn GetTransactionsCallback
) error
```

**What it does:**
- Iterates through all transactions for an address
- Checks `ctx.Err()` every iteration through the database
- Can be cancelled while processing thousands of transactions

**Use case:** Addresses with many transactions (exchanges, popular contracts)

#### GetTxAddressesContext
```go
func (d *RocksDB) GetTxAddressesContext(
    ctx context.Context,
    txid string
) (*TxAddresses, error)
```

**What it does:**
- Retrieves transaction address data
- Checks context before expensive database lookups
- Used by GetSpendingTxid operations

### API Worker Layer (api/worker.go)

#### GetSpendingTxidContext
```go
func (w *Worker) GetSpendingTxidContext(
    ctx context.Context,
    txid string,
    n int
) (string, error)
```

**What it does:**
- Finds which transaction spent a given output
- Can involve expensive database queries
- Previously caused 2+ minute response times

**Use case:** `/spending/` endpoint operations

#### GetAddressContext
```go
func (w *Worker) GetAddressContext(
    ctx context.Context,
    address string,
    page, txsOnPage int,
    option AccountDetails,
    filter *AddressFilter,
    secondaryCoin string
) (*Address, error)
```

**What it does:**
- Retrieves address balance and transaction history
- Calls GetAddrDescTransactionsContext internally
- Can process thousands of transactions per address

**Use case:** Address explorer pages and API endpoints

## Context Integration Points

### 1. HTTP Handlers → API Methods

```go
// In explorerSpendingTx
spendingTx, err := s.api.GetSpendingTxidContext(r.Context(), tx, n)

// In apiAddress
address, err := s.api.GetAddressContext(r.Context(), addressParam, page, ...)

// In explorerAddress
address, err := s.api.GetAddressContext(r.Context(), addressParam, page, ...)
```

### 2. API Methods → Database Operations

```go
// In GetSpendingTxidContext
tsp, err := w.db.GetTxAddressesContext(ctx, txid)

// In GetAddressContext
err = w.db.GetAddrDescTransactionsContext(ctx, addrDesc, lower, higher, fn)
```

## Performance Characteristics

### Cancellation Speed

```go
// Context check performance
if ctx.Err() != nil {  // ~1-10 nanoseconds
    return ctx.Err()
}
```

**Why it's fast:**
- No system calls or network operations
- Simple channel state check
- Go runtime optimization

### Cancellation Granularity

| Operation | Cancellation Points | Typical Frequency |
|-----------|-------------------|------------------|
| GetAddrDescTransactions | Every database record | 100-10,000+ per second |
| GetSpendingTxid | Before expensive operations | 1-10 per operation |
| Iterator loops | Every iteration | 1,000-100,000+ per second |

## Example Scenarios

### Scenario 1: Large Address Query

```
Request: GET /address/1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa?page=1

Timeline:
0.0s  → Request starts, context timeout = 45s
0.1s  → GetAddressContext called
0.2s  → GetAddrDescTransactions iterator starts
0.3s  → Processing transaction 1 (ctx.Err() = nil)
0.4s  → Processing transaction 2 (ctx.Err() = nil)
...
44.9s → Processing transaction 50,000 (ctx.Err() = nil)
45.0s → Timer fires, context cancelled
45.0s → Processing transaction 50,001 (ctx.Err() = DeadlineExceeded)
45.0s → Iterator stops, returns timeout error
45.0s → User sees "Request timeout - operation took too long"
```

### Scenario 2: Quick Operation (No Timeout)

```
Request: GET /spending/abc123.../0

Timeline:
0.0s → Request starts, context timeout = 45s
0.1s → GetSpendingTxidContext called
0.2s → GetTxAddressesContext called (ctx.Err() = nil)
0.3s → Database lookup completes
0.3s → Returns spending transaction ID
0.3s → User redirected to spending transaction
```

## Error Handling

### Context Timeout Error
```go
if errors.Is(err, context.DeadlineExceeded) {
    // Request timed out
    return api.NewAPIError("Request timeout - operation took too long", true)
}
```

### Context Cancellation Error
```go
if errors.Is(err, context.Canceled) {
    // Request was cancelled (client disconnected)
    return api.NewAPIError("Request cancelled", true)
}
```

## Backwards Compatibility

All context-aware methods have non-context versions for backwards compatibility:

```go
// New context-aware version
func (d *RocksDB) GetTxAddressesContext(ctx context.Context, txid string) (*TxAddresses, error)

// Original version (calls context version with background context)
func (d *RocksDB) GetTxAddresses(txid string) (*TxAddresses, error) {
    return d.GetTxAddressesContext(context.Background(), txid)
}
```

**Migration strategy:**
1. Existing code continues to work unchanged
2. New timeout-sensitive code uses context versions
3. Gradual migration path for performance-critical sections

## Debugging Context Operations

### Enable Context Logging

```go
// In development, you can add context debugging
if glog.V(3) {
    glog.Infof("Context check: %v", ctx.Err())
}
```

### Monitor Context Timeouts

```bash
# Watch for context timeout errors in logs
tail -f /var/log/blockbook.log | grep "timeout"

# Monitor specific operations
tail -f /var/log/blockbook.log | grep "GetSpendingTxid"
```

### Test Context Cancellation

```bash
# Test with very short timeout
./blockbook --requesttimeout=5 ...

# Make request to slow endpoint
curl http://localhost:9130/address/[large-address]

# Should see timeout after 5 seconds
```

## Best Practices

### 1. Context Propagation
Always pass context through the call chain:
```go
// Good: Context flows through all layers
handler(r.Context()) → api(ctx) → database(ctx)

// Bad: Context lost in chain
handler(r.Context()) → api() → database()
```

### 2. Context Checking Frequency
Check context in long-running loops:
```go
// Good: Check every iteration
for i := 0; i < largeNumber; i++ {
    if ctx.Err() != nil {
        return ctx.Err()
    }
    // Process item
}

// Bad: No context checking in loop
for i := 0; i < largeNumber; i++ {
    // Process item (can't be cancelled)
}
```

### 3. Resource Cleanup
Always clean up resources when context is cancelled:
```go
defer iterator.Close()  // Always clean up
if ctx.Err() != nil {
    return ctx.Err()    // Context handles the rest
}
```

### 4. Timeout Values
Set appropriate timeouts for your use case:
- **APIs:** 15-45 seconds (balance user experience vs. completion)
- **Admin tools:** 60-180 seconds (allow for complex operations)
- **Background jobs:** No timeout or very long timeout

## Future Enhancements

Potential areas for expanding context-aware operations:

1. **More Database Methods:** Add context to remaining slow operations
2. **Blockchain RPC:** Add context to backend blockchain calls
3. **Network Operations:** Add context to external API calls
4. **Batch Operations:** Add context to bulk processing operations

Context-aware operations provide the foundation for fine-grained resource control and can be extended to any operation that might take significant time or resources.