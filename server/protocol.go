package server

const SingleByte = 1
const Kilobyte = 1024 * SingleByte
const Megabyte = 1024 * Kilobyte
const Gigabyte = 1024 * Megabyte

// Follows limit set in the official implementation
// https://github.com/BeamMP/BeamMP-Server/blob/minor/src/TNetwork.cpp#L476
const MaxHeaderSize = 100 * Megabyte
