# Protocol Documentation

## Message Format
Each message follows this general structure:
* Total packet length (4 bytes)
* Command name length (2 bytes)
* Command name string
* Sequence of fields, each containing:
  * Field length (4 bytes)
  * Field data

## Commands


### RegisterClient

**Command Structure:**
```
[Total Length: 4 bytes]
[Command Length: 2 bytes]["RegisterClient"]
[Field 1 Length: 4 bytes][Address (string)]
[Field 2 Length: 4 bytes][Content (string)]
```

**Fields:**
1. **Address**
   * Type: string
2. **Content**
   * Type: string


### SendMessage

**Command Structure:**
```
[Total Length: 4 bytes]
[Command Length: 2 bytes]["SendMessage"]
[Field 1 Length: 4 bytes][Address (string)]
[Field 2 Length: 4 bytes][Content (string)]
```

**Fields:**
1. **Address**
   * Type: string
2. **Content**
   * Type: string



## Notes
* All integer values are in network byte order (big-endian)
* String fields are UTF-8 encoded
* Timeout for message decoding: 10 seconds
