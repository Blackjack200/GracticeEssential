package main

import "github.com/Blackjack200/GracticeEssential/bootstrap"

func main() {
	log := bootstrap.NewLogger()
	bootstrap.Default(log, nil, nil)()
}
