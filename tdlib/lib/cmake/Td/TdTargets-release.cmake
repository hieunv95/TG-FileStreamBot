#----------------------------------------------------------------
# Generated CMake target import file for configuration "Release".
#----------------------------------------------------------------

# Commands may need to know the format version.
set(CMAKE_IMPORT_FILE_VERSION 1)

# Import target "Td::tdjson" for configuration "Release"
set_property(TARGET Td::tdjson APPEND PROPERTY IMPORTED_CONFIGURATIONS RELEASE)
set_target_properties(Td::tdjson PROPERTIES
  IMPORTED_LOCATION_RELEASE "${_IMPORT_PREFIX}/lib/libtdjson.so.1.8.45"
  IMPORTED_SONAME_RELEASE "libtdjson.so.1.8.45"
  )

list(APPEND _IMPORT_CHECK_TARGETS Td::tdjson )
list(APPEND _IMPORT_CHECK_FILES_FOR_Td::tdjson "${_IMPORT_PREFIX}/lib/libtdjson.so.1.8.45" )

# Commands beyond this point should not need to know the version.
set(CMAKE_IMPORT_FILE_VERSION)
