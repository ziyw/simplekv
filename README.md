# simplekv: 

simplekv is a persist key-value store. 

Inspired by [Bitcast](https://docs.riak.com/riak/kv/2.2.3/setup/planning/backend/bitcask/index.html), it uses append-only log to persist key-value. Suitable for large writes with limited keys. 

## Storage Engine Details
Writes are persisted on **segment** file. Each **Segment** file has a in-memory hashmap. 

For each write:
1. find the current segment 
2. append key-value content directly to segment file
3. update hashmap with new (key, offset) pair. Offset is the offset of key-value in segment file 

For each read: find the key in current segment hashmap, if not exist, load previous segment, repeat until find. 

**Segment**:
- Compress: create a new segment file from a existing segment file, with only the latest key-value pairs, reduce the duplicate writes.
- Merge: create a new segment file from two existing segment files, merge two files and only keep the latest key-value pairs. 
Compress and merge happens periodically in the storage engine background.

**hashmap**:
- Snapshot: persist the in-memory hashmap to disk by creating a snapshot. Segment can then load this snapshot back to memory when loading a previous segment.

