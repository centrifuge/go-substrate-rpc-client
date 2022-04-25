# Types Test

The purpose of this test is to ensure that any types can be decoded successfully by using live blockchain
data (i.e. metadata and storage). In an effort to remove dependencies on external service providers such as `subscan` for this test,
the test data can be generated on demand by using the [test-gen](test-gen/main.go) tool, which is retrieving the data 
based on the `test-gen-*` tags set on the `types.EventRecords` struct.

The generation of test data and test execution can be triggered via the
`generate-test-data` and `test-types-decode` targets, see [Makefile](../../Makefile). 