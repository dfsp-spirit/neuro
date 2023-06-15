// Provides functions for reading and writing some structural neuroimaging file formats.
//
// The neuro package is intended to be used in neuroimaging, with a focus on structural
// brain anatomy. Currently it provides functions to access
// fileformats used by FreeSurfer and some related neuroimaging software packages.
//
// The package can read three-dimensional (3D) and 4D brain scans in MGH and MGZ format,
// typically produced from the raw DICOM files that are written by magnetic resonance imaging
// (MRI) hardware. Support for reading brain surface reconstructions (cortical meshes)
// and the related per-vertex data (like cortical thickness or sulcal depth at each point of the brain surface)
// are also included.
package neuro

// Verbosity is the verbosity level of the package. 0 = silent, 1 = info, 2 = debug.
var Verbosity int = 0 // WARNING: If you increase Verbosity here and run the unit tests, the examples included with the tests will fail, because they expect a defined output on STDOUT, and increasing verbosity will produce extra output.
