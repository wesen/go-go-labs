# Changelog

## Bank Statement Parser Implementation

Added a Bank of America statement parser that can process CSV files containing bank statement data. The parser:

- Detects accounts and their details
- Processes transactions (deposits, withdrawals, service fees)
- Handles multi-line transactions
- Maintains running totals in account summaries
- Includes comprehensive test coverage

Changes:
- Created pkg/textract/boa/processor.go with core parsing logic
- Created pkg/textract/boa/processor_test.go with test suite
- Implemented according to boa-format-spec.md specification 

## Writer Update for Unified Transaction Storage

Updated the writer implementation to use a unified transaction storage approach:

- Modified writer.go to work with the single transactions map
- Improved transaction categorization during output generation
- Optimized summary calculations to use the unified storage
- Maintained backward compatibility with existing output formats 

## Pattern Matching Improvements

Fixed transaction parsing issues with more precise pattern matching:

- Made date pattern require match at start of line to avoid false positives
- Made amount pattern require exact line match
- Improved multi-line transaction handling
- Fixed test failures in TestMultiLineTransaction 

## Debug Logging Enhancement

Added comprehensive debug logging to track transaction processing:

- Added context-rich logging for each processing step
- Included detailed transaction state information
- Added summary updates tracking
- Improved debugging capabilities for multi-line transactions 

## Security and Testing Improvements

Enhanced security and test coverage:

- Added account number anonymization in logs
- Preserved first/last 4 digits for traceability
- Added test for multiple section transitions
- Added test for account number anonymization
- Improved test coverage for edge cases 

## Test Data Security

Improved security of test data:

- Replaced real account numbers with test ones
- Added test account number constants
- Updated all test cases to use anonymized accounts
- Maintained test coverage while improving security 

## API Simplification for Line-Based Processing

Simplified the Bank of America statement processing API:
- Removed table processing to focus on line-based parsing
- Updated boa.go to use the new streamlined API
- Removed unused table-related code and interfaces
- Simplified the textractPage implementation 

## Input Format Support Enhancement

Added support for multiple input formats in the BoA command:
- Added CSV file processing capability
- Maintained existing JSON file support
- Automatic format detection based on file extension
- Robust error handling for malformed input
- Skip invalid rows in CSV files gracefully 