// Provides functions for reading some structural neuroimaging file formats.
package neuro

// Verbosity is the verbosity level of the package. 0 = silent, 1 = info, 2 = debug.
var Verbosity int = 0  // WARNING: If you increase Verbosity here and run the unit tests, the examples included with the tests will fail, because they expect a defined output on STDOUT, and increasing verbosity will produce extra output.