cmake_minimum_required(VERSION 3.14)

# Use CPM for CMake dependency management
include(cmake/CPM.cmake)

CPMAddPackage("gh:CLIUtils/CLI11@2.3.1")

add_subdirectory(evmone)

project(evmone-runner LANGUAGES CXX C)

add_executable(runner runner.cpp)
target_link_libraries(runner PRIVATE CLI11::CLI11 evmone)

target_include_directories(runner PUBLIC ${PROJECT_BINARY_DIR} ${PROJECT_SOURCE_DIR}/evmone)