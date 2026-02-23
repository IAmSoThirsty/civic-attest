# Canonical Encoding Specification

**Version:** 2.0
**Status:** Formal Specification
**Last Updated:** 2026-02-23
**Classification:** PUBLIC

## Abstract

This document provides a complete, formal specification for deterministic canonical encoding used in Civic Attest cryptographic operations. Canonicalization is critical for ensuring reproducible hashes and preventing ambiguity attacks.

**Invariant:** For any input X: `canonical(X) == canonical(canonical(X))`

**Requirement:** All implementations MUST produce byte-identical output for the same logical input.

## 1. Overview

### 1.1 Purpose

Canonical encoding ensures that:
1. The same logical data always produces the same byte stream
2. Hash values are reproducible across implementations
3. Signature verification is deterministic
4. No ambiguity in data representation

### 1.2 Scope

This specification covers:
- Canonical CBOR (binary encoding)
- Canonical JSON (text encoding)
- Unicode normalization
- Floating point handling
- Data type restrictions

## 2. Canonical CBOR

**Base Standard:** RFC 8949 Section 4.2 (Deterministic Encoding)

**Additional Requirements:** This specification adds stricter requirements beyond RFC 8949.

### 2.1 Integer Encoding

**Rule:** Use shortest possible encoding

**Compliance:**
- Values 0-23: Single byte (major type 0 or 1)
- Values 24-255: Additional 1 byte
- Values 256-65535: Additional 2 bytes (network byte order)
- Values 65536-4294967295: Additional 4 bytes
- Values ≥ 4294967296: Additional 8 bytes

**Prohibited:**
- Leading zeros
- Overlong encodings
- Negative zero (use positive zero)

**Examples:**
```
Correct:   10 → 0x0A (single byte)
Incorrect: 10 → 0x1801 0x000A (overlong)

Correct:   255 → 0x18 0xFF
Correct:   256 → 0x19 0x0100
```

### 2.2 Floating Point

**Rule:** PROHIBITED in cryptographic contexts

**Rationale:** Floating point representation has multiple equivalent encodings:
- Different NaN representations
- Positive vs negative zero
- Subnormal number variations

**Alternative:** Use rational numbers or fixed-point decimal strings

**If Absolutely Required (non-cryptographic metadata only):**
- Use smallest representation (half-precision if possible)
- Normalize -0.0 to +0.0
- Use canonical NaN: 0x7e00 (half), 0x7fc00000 (single), 0x7ff8000000000000 (double)
- No infinity in canonical form (use null or string "infinity")

**Examples:**
```
PROHIBITED (cryptographic): price = 19.99 (float)
CORRECT: price_cents = 1999 (integer)
CORRECT: price = "19.99" (string)
CORRECT: price = { "numerator": 1999, "denominator": 100 }
```

### 2.3 String Encoding

**Rule:** UTF-8 encoding with NFC normalization

**Requirements:**
1. Apply Unicode NFC (Canonical Decomposition + Canonical Composition)
2. Encode as UTF-8
3. No overlong UTF-8 sequences
4. No byte order mark (BOM)
5. Use major type 3 (text string) for text
6. Use major type 2 (byte string) for binary data

**Unicode Normalization:**
- Form: NFC (Canonical Decomposition followed by Canonical Composition)
- Apply BEFORE encoding to UTF-8
- Use Unicode version 15.0.0 or later

**Prohibited:**
- NFD, NFKC, NFKD forms in cryptographic contexts
- Non-normalized Unicode
- Overlong UTF-8 (e.g., 0xC0 0x80 for NULL)

**Examples:**
```
Input: "café" (U+0063 U+0061 U+0066 U+00E9)
Normalized (NFC): U+0063 U+0061 U+0066 U+00E9 (already normalized)
UTF-8: 0x63 0x61 0x66 0xC3 0xA9

Input: "café" (U+0063 U+0061 U+0066 U+0065 U+0301) [decomposed]
Normalized (NFC): U+0063 U+0061 U+0066 U+00E9 (composed)
UTF-8: 0x63 0x61 0x66 0xC3 0xA9

Both inputs MUST produce identical CBOR encoding.
```

### 2.4 Arrays

**Rule:** Definite-length encoding only

**Requirements:**
- Use definite-length encoding (major type 4)
- Encode length in shortest form
- No indefinite-length arrays (0x9F ... 0xFF)
- No duplicate elements (for sets)

**Examples:**
```
CORRECT: [1, 2, 3] → 0x83 0x01 0x02 0x03
INCORRECT: [1, 2, 3] → 0x9F 0x01 0x02 0x03 0xFF (indefinite)
```

### 2.5 Maps (Objects)

**Rule:** Definite-length with sorted keys

**Requirements:**
1. Use definite-length encoding (major type 5)
2. Sort keys by canonical encoding byte order (lexicographic)
3. No duplicate keys
4. Keys can be any CBOR type
5. For text string keys: sort by UTF-8 byte order AFTER NFC normalization

**Sorting Algorithm:**
```
function compare_cbor_keys(key1, key2):
  encoded1 = canonical_cbor_encode(key1)
  encoded2 = canonical_cbor_encode(key2)
  return byte_compare(encoded1, encoded2)

function byte_compare(bytes1, bytes2):
  for i in 0 to min(len(bytes1), len(bytes2)):
    if bytes1[i] < bytes2[i]: return -1
    if bytes1[i] > bytes2[i]: return +1
  if len(bytes1) < len(bytes2): return -1
  if len(bytes1) > len(bytes2): return +1
  return 0
```

**Example:**
```
Input map:
{
  "zebra": 1,
  "apple": 2,
  "Zebra": 3
}

Sorted by UTF-8 byte order:
{
  "Zebra": 3,   // 'Z' = 0x5A
  "apple": 2,   // 'a' = 0x61
  "zebra": 1    // 'z' = 0x7A
}
```

### 2.6 Disallowed Types

**Prohibited in Cryptographic Contexts:**
- Undefined (major type 7, value 23)
- Simple values except false, true, null
- Tags (major type 6) unless explicitly whitelisted
- Indefinite-length encoding (all types)
- Floating point (see section 2.2)

**Allowed Simple Values:**
- false (0xF4)
- true (0xF5)
- null (0xF6)

**Whitelisted Tags (if needed):**
- Tag 0: RFC 3339 date/time string (encoded as canonical text string)
- Tag 2: Big unsigned integer (only if value exceeds uint64)
- Tag 3: Big negative integer (only if value exceeds int64)

### 2.7 CBOR Encoding Algorithm

**Pseudocode:**
```
function canonical_cbor_encode(value):
  if is_integer(value):
    return encode_integer_shortest(value)

  if is_string(value):
    normalized = unicode_nfc_normalize(value)
    utf8 = encode_utf8(normalized)
    return cbor_text_string(utf8)

  if is_bytes(value):
    return cbor_byte_string(value)

  if is_array(value):
    encoded_items = [canonical_cbor_encode(item) for item in value]
    return cbor_array_definite(encoded_items)

  if is_map(value):
    sorted_keys = sort_by_canonical_encoding(value.keys())
    encoded_pairs = []
    for key in sorted_keys:
      encoded_key = canonical_cbor_encode(key)
      encoded_value = canonical_cbor_encode(value[key])
      encoded_pairs.append((encoded_key, encoded_value))
    return cbor_map_definite(encoded_pairs)

  if is_boolean(value) or is_null(value):
    return cbor_simple_value(value)

  else:
    error("Type not allowed in canonical CBOR")
```

## 3. Canonical JSON

**Base Standard:** RFC 8785 (JSON Canonicalization Scheme - JCS)

**Additional Requirements:** This specification adds domain-specific requirements.

### 3.1 Unicode Normalization

**Rule:** NFC normalization before encoding

**Process:**
1. Apply NFC normalization to all strings
2. Encode as UTF-8
3. Apply JCS string escaping rules

**Example:**
```
Input: { "text": "café" }  // U+0063 U+0061 U+0066 U+0065 U+0301 (decomposed)
Normalized: { "text": "café" }  // U+0063 U+0061 U+0066 U+00E9 (composed)
JSON: {"text":"café"}
```

### 3.2 Whitespace

**Rule:** Remove all unnecessary whitespace

**Requirements:**
- No spaces around ':' or ','
- No newlines
- No indentation
- No trailing whitespace

**Example:**
```
CORRECT: {"a":1,"b":2}
INCORRECT: { "a": 1, "b": 2 }
INCORRECT: {
  "a": 1,
  "b": 2
}
```

### 3.3 String Escaping

**Rule:** Minimal necessary escaping

**Requirements:**
1. Escape control characters U+0000 through U+001F
2. Escape quotation mark (U+0022) as \"
3. Escape reverse solidus (U+005C) as \\
4. Use lowercase hex in escapes: \uXXXX
5. Do NOT escape forward slash (U+002F)
6. Do NOT escape Unicode characters > U+001F (use UTF-8 directly)

**Examples:**
```
Input: "hello\nworld"
Canonical: "hello\u000aworld"

Input: "quote: \" end"
Canonical: "quote: \" end"

Input: "path/to/file"
Canonical: "path/to/file"  (NOT "path\/to\/file")
```

### 3.4 Number Representation

**Rule:** No floating point, use strings for precision

**Requirements:**
- Integers: No leading zeros (except "0")
- No leading '+' sign
- Negative: Use '-' prefix
- **CRITICAL:** No floating point representation in cryptographic contexts
- For monetary/precise values: Use string representation or separate numerator/denominator

**Examples:**
```
CORRECT (cryptographic): {"amount": "19.99"}
CORRECT (cryptographic): {"amount_cents": 1999}
CORRECT (cryptographic): {"amount": {"numerator": 1999, "denominator": 100}}

PROHIBITED (cryptographic): {"amount": 19.99}

CORRECT (integer): {"count": 42}
INCORRECT: {"count": 042}  // Leading zero
INCORRECT: {"count": +42}  // Leading plus
```

### 3.5 Map Key Sorting

**Rule:** Sort by UTF-16 code unit order

**Process (per RFC 8785):**
1. Normalize strings with NFC
2. Convert to UTF-16
3. Sort by code unit values (not code points)
4. Apply escape sequences after sorting

**Example:**
```
Input:
{
  "zebra": 1,
  "apple": 2,
  "🦓": 3
}

Sorted (UTF-16 code units):
{
  "apple": 2,
  "zebra": 1,
  "🦓": 3
}
```

### 3.6 Boolean and Null

**Rule:** Lowercase only

**Values:**
- true (not True or TRUE)
- false (not False or FALSE)
- null (not Null or NULL or nil)

### 3.7 Disallowed in Canonical JSON

**Prohibited:**
- Comments (// or /* */)
- Trailing commas
- Duplicate keys
- NaN or Infinity (use null or string representation)
- Undefined (not a JSON type)
- Floating point numbers in cryptographic contexts

### 3.8 JSON Encoding Algorithm

**Pseudocode:**
```
function canonical_json_encode(value):
  if is_string(value):
    normalized = unicode_nfc_normalize(value)
    escaped = jcs_escape_string(normalized)
    return '"' + escaped + '"'

  if is_integer(value):
    return to_string_no_leading_zeros(value)

  if is_boolean(value):
    return value ? "true" : "false"

  if is_null(value):
    return "null"

  if is_array(value):
    items = [canonical_json_encode(item) for item in value]
    return "[" + join(items, ",") + "]"

  if is_object(value):
    sorted_keys = sort_by_utf16_code_units(value.keys())
    pairs = []
    for key in sorted_keys:
      encoded_key = canonical_json_encode(key)
      encoded_value = canonical_json_encode(value[key])
      pairs.append(encoded_key + ":" + encoded_value)
    return "{" + join(pairs, ",") + "}"

  else:
    error("Type not allowed in canonical JSON")
```

## 4. Test Vectors

### 4.1 CBOR Test Vectors

**Test 1: Simple Map**
```
Input (logical):
{
  "b": 2,
  "a": 1
}

Canonical CBOR (hex):
A2                    # map(2)
   61                 # text(1)
      61              # "a"
   01                 # unsigned(1)
   61                 # text(1)
      62              # "b"
   02                 # unsigned(2)
```

**Test 2: Unicode Normalization**
```
Input: { "café": 1 } with NFD encoding (U+0065 U+0301)
Normalized: { "café": 1 } with NFC encoding (U+00E9)

Canonical CBOR (hex):
A1                    # map(1)
   64                 # text(4)
      63 61 66 E9     # "café" (NFC, é = 0xC3A9 in UTF-8 → 0xE9 in text)
   01                 # unsigned(1)

Note: The actual UTF-8 for é is 0xC3 0xA9
Correct CBOR hex:
A1 64 63 61 66 C3 A9 01
```

### 4.2 JSON Test Vectors

**Test 1: Whitespace and Sorting**
```
Input (formatted):
{
  "z": 3,
  "a": 1,
  "m": 2
}

Canonical JSON:
{"a":1,"m":2,"z":3}
```

**Test 2: String Escaping**
```
Input:
{
  "text": "Line 1\nLine 2",
  "path": "C:/folder/file.txt"
}

Canonical JSON:
{"path":"C:/folder/file.txt","text":"Line 1\u000aLine 2"}
```

### 4.3 Cross-Format Consistency

**Requirement:** The same logical data MUST hash to the same value regardless of format (CBOR vs JSON) when semantically equivalent.

**Example:**
```
Logical data: { "name": "Alice", "age": 30 }

CBOR encoding: [specific bytes]
JSON encoding: {"age":30,"name":"Alice"}

Hash(CBOR) may differ from Hash(JSON), but:
- Hash(CBOR1) == Hash(CBOR2) for same logical data
- Hash(JSON1) == Hash(JSON2) for same logical data
- Applications declare which format is canonical for their context
```

## 5. Validation and Testing

### 5.1 Conformance Testing

**Requirements:**
1. Implementations MUST pass all test vectors
2. Round-trip test: canonical(X) == canonical(canonical(X))
3. Cross-implementation compatibility
4. Fuzzing for edge cases

**Test Suite:**
- Minimum 100 test vectors per format
- Edge cases for each data type
- Unicode corner cases
- Sorting edge cases
- Error cases

### 5.2 Reference Implementations

**CBOR:**
- cbor-deterministic (Go)
- cbor2 with deterministic mode (Python)
- CBOR.js with canonical mode (JavaScript)

**JSON:**
- JCS reference implementation (multiple languages)
- json-canonicalize (JavaScript)
- go-json-canonicalize (Go)

### 5.3 Validation Tools

**Canonical Encoding Validator:**
```bash
# Validate CBOR encoding
$ civic-attest-validator --format cbor --input data.cbor

# Validate JSON encoding
$ civic-attest-validator --format json --input data.json

# Compare two encodings for same logical data
$ civic-attest-validator --compare data1.cbor data2.cbor
```

## 6. Implementation Guidelines

### 6.1 Security Considerations

**CRITICAL:**
1. Never use floating point for cryptographic data
2. Always normalize Unicode before encoding
3. Always use canonical encoding for signature/hash operations
4. Validate encoding before accepting external data
5. Reject non-canonical encodings in signature verification

### 6.2 Performance

**Optimization:**
- Cache normalized strings
- Use efficient sorting algorithms (O(n log n))
- Pre-compute canonical encoding for frequently used structures
- Use streaming encoders where possible

### 6.3 Error Handling

**Reject (do not attempt to fix):**
- Non-normalized Unicode
- Floating point in cryptographic contexts
- Non-canonical encodings
- Duplicate map keys
- Disallowed types

**Never:**
- Silently convert to canonical form
- Guess intended encoding
- Accept "close enough" encodings

## 7. Versioning

### 7.1 Encoding Version Field

**All cryptographic structures include:**
```json
{
  "canonical_encoding_version": "2.0",
  "encoding_type": "CBOR_DETERMINISTIC",
  "unicode_normalization": "NFC",
  ...
}
```

### 7.2 Migration

**Version 1.x → 2.0:**
- Add Unicode normalization requirements
- Add floating point prohibition
- Add stricter sorting requirements
- Maintain backward verification compatibility

**Future Versions:**
- Must maintain ability to verify old signatures
- New restrictions only apply to new signatures
- Clear migration path documented

## 8. Compliance

### 8.1 Standards Compliance

**MUST Comply:**
- RFC 8949 Section 4.2 (CBOR Deterministic Encoding)
- RFC 8785 (JSON Canonicalization Scheme)
- Unicode Standard 15.0+ (NFC normalization)

**MUST NOT:**
- Deviate from this specification in cryptographic contexts
- Create custom encodings
- Mix canonical and non-canonical encodings

### 8.2 Certification

**Implementation Certification:**
1. Pass official test suite (100% pass rate)
2. Cross-implementation compatibility verified
3. Security audit of encoding logic
4. Performance benchmarks meet requirements

## 9. Appendices

### Appendix A: Unicode Normalization Reference

**NFC Process:**
1. Canonical Decomposition: Decompose characters into base + combining marks
2. Canonical Ordering: Reorder combining marks by combining class
3. Canonical Composition: Recompose where possible

**Example:** é (U+00E9) ↔ e (U+0065) + ́ (U+0301)

### Appendix B: Complete Test Suite

See: `tests/canonical-encoding/`

### Appendix C: Migration Guide

See: `docs/canonical-encoding-migration.md`

---

**Document Status:** Formal Specification
**Compliance Required:** Yes (all implementations)
**Review Cycle:** Annual
**Next Review:** 2027-02-23
