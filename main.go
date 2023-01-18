package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type Player struct {
	id             int
	currentStamina int
	fullStamina    int // max capacity of stamina player has
	team           int // which team
}

var score = 0
var turn = 0
var maxTurns = 10

func main() {
	rand.Seed(time.Now().UnixNano())

	players := [10]Player{}

	for i := 0; i < len(players); i++ {
		players[i].id = i + 1
		players[i].fullStamina = rand.Intn(5) + 1
		players[i].currentStamina = players[i].fullStamina
		players[i].team = (i % 2) + 1
	}
	log.Print(players)

	var wg sync.WaitGroup
	wg.Add(1) // only wait for play game, not counting

	go func() {
		playGame(players)
		wg.Done()
	}()

	go func() {
		for {
			log.Println("Turn:", turn, "Current score:", score)
			turn += 1
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()

	if score > 0 {
		log.Println("Team 1 won")
	} else if score < 0 {
		log.Println("Team 2 won")
	} else {
		log.Println("It's a tie!")
	}
}

func playGame(players [10]Player) {
	for turns := 0; turns < maxTurns; turns++ {
		var wg sync.WaitGroup
		wg.Add(len(players)) // Increment the WaitGroup counter by the number of players

		for i := 0; i < len(players); i++ {
			go func(player *Player) {
				if player.currentStamina == 0 {
					rest(player)
				} else if player.currentStamina == player.fullStamina {
					pull(player)
				} else {
					if rand.Intn(2) == 0 {
						rest(player)
					} else {
						pull(player)
					}
				}
				wg.Done() // Decrement the WaitGroup counter by 1
			}(&players[i])
		}
		wg.Wait() // Block until the WaitGroup counter is 0
		time.Sleep(time.Second)
		log.Println("Players Status:", players)
	}
}

func pull(player *Player) {
	moved := rand.Intn(player.currentStamina) + 1
	player.currentStamina -= moved
	if player.team == 1 {
		moved = 0 + moved
	} else {
		moved = 0 - moved
	}
	log.Println("Player", player.id, "added", moved, "for team", player.team)
	score += moved
}

func rest(player *Player) {
	min := 1
	max := player.fullStamina - player.currentStamina
	rests := rand.Intn(max-min+1) + min

	log.Println("Player", player.id, "is resting for", rests, "seconds")
	time.Sleep(time.Duration(rests))
	player.currentStamina = player.fullStamina
}
