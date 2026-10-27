package lobby

import "time"

// StampCooldown is the minimum gap between one player's lobby stamps, matching
// the in-game one: long enough that a held tap cannot bury everyone's screen,
// short enough that a reply to someone else's stamp never feels blocked.
const StampCooldown = time.Second * 3
