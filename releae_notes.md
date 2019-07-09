# v0.0.2
Still Alpha as interfaces may change, however it should be easy to convert

1. Added support for adding commands using parameters by name
1. Made command apis internal
    1. these were lower level functions that not as easy to use
1. Added higher level functions for running a script or command
    1. Added support for named arguments and unnamed arguments
1. Support for Calling into powershell from callback
    1. Known as NestedPowershell / NestedPipeline
1. Reference dependency library from checked in file via git
    1. This makes project much easier to consume
1. Updated docs & examples
    1. Increase test coverage
    1. Remove usage of examples for internal tests
1. Standardized cleanup routines to be all call Close
    1. Previously some where called Delete
1. checked in x64 binaries to make project easier to consume
    1. binaries checked in separate repo, linked via submodule

# v0.0.1
Initial drop

1. Alpha as interfaces may change, however it should be easy to convert
