# Timeout Configuration and Troubleshooting Guide

This guide covers Blockbook's comprehensive timeout protection system and how to configure it for optimal performance.

## Overview

Blockbook implements multi-layered timeout protection to prevent resource exhaustion from slow requests, bot attacks, and database operations that take too long.

## Timeout Layers

### 1. HTTP Connection Timeouts
These timeouts protect against slow or malicious connections at the HTTP server level.

```bash
--httpreadtimeout=30        # Max time to read request headers/body (default: 30s)
--httpwritetimeout=60       # Max time to write response (default: 60s)
--httpidletimeout=120       # Keep-alive connection timeout (default: 120s)
--httpreadheadertimeout=10  # Max time to read headers (prevents slowloris attacks, default: 10s)
```

**What they protect against:**
- Slowloris attacks (slow header sending)
- Clients that send data very slowly
- Clients that don't read responses
- Idle connections consuming server resources

### 2. Request Operation Timeouts
These timeouts protect against slow database operations and API calls.

```bash
--requesttimeout=45         # Max time for database-heavy operations (default: 45s)
```

**What they protect against:**
- Long-running address queries (addresses with many transactions)
- Slow GetSpendingTxid operations
- Heavy RocksDB iterator operations
- XPUB queries with many addresses

### 3. Context-Aware Database Operations
The most granular timeout protection - database operations can be cancelled mid-operation.

**Protected operations:**
- `GetAddrDescTransactions` - Iterator loops check context every iteration
- `GetSpendingTxid` - Transaction spending lookups
- `GetAddress` - Address balance and transaction history
- `GetTxAddresses` - Transaction address data

## Configuration Examples

### High-Traffic Production Server
```bash
./blockbook -sync -blockchaincfg=build/blockchaincfg.json \
  --httpreadtimeout=15 \
  --httpwritetimeout=30 \
  --httpidletimeout=60 \
  --httpreadheadertimeout=5 \
  --requesttimeout=30
```
**Use case:** High-traffic servers with many bot requests

### Development/Testing Server
```bash
./blockbook -sync -blockchaincfg=build/blockchaincfg.json \
  --httpreadtimeout=60 \
  --httpwritetimeout=120 \
  --httpidletimeout=300 \
  --httpreadheadertimeout=20 \
  --requesttimeout=90
```
**Use case:** Development environments or slower hardware

### Resource-Constrained Server
```bash
./blockbook -sync -blockchaincfg=build/blockchaincfg.json \
  --httpreadtimeout=20 \
  --httpwritetimeout=45 \
  --httpidletimeout=90 \
  --requesttimeout=25
```
**Use case:** Servers with limited RAM/CPU

## Troubleshooting Common Issues

### Issue: "Request timeout - operation took too long"

**Symptoms:**
- Users see timeout error pages
- Log entries: `Request timeout after 45s for GET /address/...`

**Solutions:**
1. **Increase request timeout:**
   ```bash
   --requesttimeout=90  # Increase from default 45s
   ```

2. **Check for database performance issues:**
   ```bash
   # Monitor RocksDB performance
   curl http://localhost:9030/api/system  # Check DB stats
   ```

3. **Optimize database if needed:**
   ```bash
   # Increase RocksDB cache if you have more RAM
   --dbcache=2147483648  # 2GB cache (default: 512MB)
   ```

### Issue: "superfluous response.WriteHeader call"

**Symptoms:**
- Warning logs about duplicate WriteHeader calls
- Usually from redirect endpoints

**Solutions:**
- This is generally harmless but indicates a redirect endpoint is using timeout wrapper unnecessarily
- Check if `/search/` or `/spending/` endpoints are functioning correctly

### Issue: Slow GetSpendingTxid Operations

**Symptoms:**
- Operations taking >2 minutes
- High memory/CPU usage during spending transaction lookups

**Solutions:**
1. **Verify timeout protection is working:**
   ```bash
   # Should see timeouts in logs after requesttimeout duration
   tail -f /var/log/blockbook.log | grep "timeout"
   ```

2. **Tune timeout for your needs:**
   ```bash
   --requesttimeout=30  # More aggressive timeout
   ```

### Issue: Connection Timeouts

**Symptoms:**
- Clients getting connection reset
- High number of open connections

**Solutions:**
1. **Adjust HTTP timeouts:**
   ```bash
   --httpidletimeout=60   # Shorter idle timeout
   --httpreadtimeout=20   # Faster read timeout
   ```

2. **Configure nginx properly (if using reverse proxy):**
   ```nginx
   proxy_read_timeout 70s;     # requesttimeout + 10s buffer
   proxy_connect_timeout 10s;
   client_body_timeout 25s;    # Less than httpreadtimeout
   ```

## Monitoring and Logging

### Key Log Messages

```bash
# Successful timeout protection
"Request timeout after 45s for GET /address/... from 1.2.3.4"
"HTML request timeout after 45s for GET /blocks from 1.2.3.4"

# Performance monitoring
"GetSpendingTxid <txid> <n>, 2.5s"  # Should be < requesttimeout

# Connection issues
"recovered from panic: http: superfluous response.WriteHeader call"
```

### Monitoring Commands

```bash
# Check current server status
curl http://localhost:9030/api/system

# Monitor active connections
netstat -an | grep :9130 | wc -l

# Monitor timeout events
journalctl -u blockbook -f | grep timeout

# Check memory usage
free -h
```

## Performance Tuning Guidelines

### For Addresses with Many Transactions
If you frequently serve addresses with thousands of transactions:

```bash
--requesttimeout=120        # Longer timeout for heavy addresses
--dbcache=4294967296        # 4GB RocksDB cache
```

### For High Bot Traffic
If you're getting many bot/spam requests:

```bash
--httpreadtimeout=10        # Aggressive connection timeouts
--httpwritetimeout=20
--requesttimeout=15         # Quick timeout for slow operations
```

### For Slow Hardware
If running on slower servers or limited resources:

```bash
--requesttimeout=180        # Allow more time for operations
--httpwritetimeout=90       # More time for responses
--syncworkers=4             # Reduce sync workers to free up resources
```

## Integration with Nginx

When running behind nginx reverse proxy, coordinate timeouts:

```nginx
# Nginx should timeout slightly after Blockbook
proxy_read_timeout 70s;       # requesttimeout + 25s
proxy_connect_timeout 10s;
proxy_send_timeout 35s;       # httpwritetimeout - 25s

# Client timeouts should be less than Blockbook HTTP timeouts
client_body_timeout 25s;      # httpreadtimeout - 5s
client_header_timeout 8s;     # httpreadheadertimeout - 2s

# Rate limiting at nginx level
limit_req_zone $binary_remote_addr zone=api:10m rate=100r/m;
limit_req zone=api burst=20 nodelay;
```

## Context Cancellation Deep Dive

### How Context Timeouts Work

1. **Timer Creation:** When request starts, Go creates a timer for `requesttimeout`
2. **Context Propagation:** Request context flows: HTTP handler → API worker → Database
3. **Periodic Checks:** Database iterators check `ctx.Err()` every loop iteration
4. **Immediate Cancellation:** When timer fires, next `ctx.Err()` check returns `DeadlineExceeded`

### Cancellation Points

```go
// Database iterator - checks every record
for it.Seek(startKey); it.Valid(); it.Next() {
    if ctx.Err() != nil {  // ← Cancellation check
        return ctx.Err()
    }
    // Process database record...
}
```

### Performance Impact

Context checks are extremely fast (nanoseconds) and have negligible performance impact while providing precise timeout control.

## Best Practices

1. **Set timeouts based on your use case:**
   - Public APIs: Aggressive timeouts (15-30s)
   - Internal tools: Lenient timeouts (60-120s)
   - Development: Very lenient (180s+)

2. **Monitor timeout events:**
   - Track frequency of timeout errors
   - Adjust timeouts based on actual usage patterns

3. **Coordinate with reverse proxy:**
   - Nginx timeouts should be longer than Blockbook timeouts
   - Allow 10-30 second buffer for error handling

4. **Test timeout behavior:**
   - Verify timeouts work as expected
   - Test with slow addresses/operations
   - Monitor resource usage during timeouts

5. **Database optimization:**
   - Increase RocksDB cache for better performance
   - Monitor database statistics
   - Consider hardware upgrades for very high traffic