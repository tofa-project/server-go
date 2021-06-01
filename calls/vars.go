package calls

// 60 seconds to connect to tofa client
var CALL_CONNECT_TIMEOUT uint = 60

// 45 seconds for tofa client to reply.
//
// 30 seconds to compensate for eventual network timeouts
var CALL_RESPONSE_TIMEOUT uint = 45 + 30

// Maximum response length from client.
//
// Used to mitigate potential DoS amid crafted URIs
//
// Default 50KBytes should be enough
var MAX_RESPONSE_LENGTH uint = 1024 * 50
