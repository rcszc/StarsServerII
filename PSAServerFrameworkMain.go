package main

import (
	"ServerSTARS-II/ServerFrameworkCore"
)

// @pomelo_star 2023-2024 RCSZ.
// SA: PSA-STARS-II (架构: 青柚之星-星辰 II) [2023_09_07]
// Request Process Server.
// Firewall => LAN => Nginx => PSA-Server.

func main() {
	ServerFrameworkCore.MainStartServer("PSA-MAIN")
}
