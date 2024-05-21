package types

const SingleByte = 1
const Kilobyte = 1024 * SingleByte
const Megabyte = 1024 * Kilobyte
const Gigabyte = 1024 * Megabyte

// Follows limit set in the official implementation
// https://github.com/BeamMP/BeamMP-Server/blob/minor/src/TNetwork.cpp#L476
const MaxHeaderSize = 100 * Megabyte

// Follows limit set in the official implementation
// https://github.com/BeamMP/BeamMP-Server/blob/5f9726f10fe3d9a353108a680b63856e1db9b9bc/src/TNetwork.cpp#L287
const MaxAuthKeyLength = 50
