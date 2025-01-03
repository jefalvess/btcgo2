package utils

import "fmt"

func Title(version string) {
	fmt.Println("\x1b[38;2;250;128;114m" + "╔═══════════════════════════════════════╗")
	fmt.Println("║\x1b[0m\x1b[36m" + "   ____ _______ _____    _____  ____   " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("║\x1b[0m\x1b[36m" + "  |  _ \\__   __/ ____|  / ____|/ __ \\  " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("║\x1b[0m\x1b[36m" + "  | |_) | | | | |      | |  __| |  | | " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("║\x1b[0m\x1b[36m" + "  |  _ <  | | | |      | | |_ | |  | | " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("║\x1b[0m\x1b[36m" + "  | |_) | | | | |____  | |__| | |__| | " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("║\x1b[0m\x1b[36m" + "  |____/  |_|  \\_____|  \\_____|\\____/  " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("║\x1b[0m\x1b[36m" + "                                       " + "\x1b[0m\x1b[38;2;250;128;114m" + "║")
	fmt.Println("╚════\x1b[32m" + "Investidor Internacional - " + version + "\x1b[0m\x1b[38;2;250;128;114m══╝" + "\x1b[0m")
}
