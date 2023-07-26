package blackjack

import(
	"fmt"
	deck "github.com/gallachaitanya/Deck"
)

type AI interface{
	Bet(shuffled bool)int
	Play(hand []deck.Card,dealer deck.Card)Move
	Results(hands [][]deck.Card,dealer []deck.Card)
}

type dealerAI struct{}

func(ai dealerAI)Bet(shuffled bool)int{
	return 1
}

func(ai dealerAI)Play(hand []deck.Card,dealer deck.Card)Move{
		dscore := Score(hand...)
		if (dscore <= 16 || (dscore == 17 && Soft(hand...))){
			return MoveHit
		}
		return MoveStand
}

func(ai dealerAI)Results(hands [][]deck.Card,dealer []deck.Card){}

func HumanAI() AI{
	return humanAI{}
}

type humanAI struct{}

func(ai humanAI) Bet(shuffled bool)int{
	if shuffled{
		fmt.Println("The deck was just shuffled.")
	}
	fmt.Println("Would you like to bet?")
	var bet int
	fmt.Scanf("%d\n",&bet)
	return bet
}

func(ai humanAI) Play(hand []deck.Card,dealer deck.Card)Move{
	for{
	fmt.Println("Player: ",hand)
		fmt.Println("Dealer: ",dealer)
		fmt.Println("What will you do? (h)it, (s)tand, (d)ouble, or s(p)lit")
		var input string
		fmt.Scanf("%s\n",&input)
		switch input{
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("Invalid input option", input)
		}
	}
}



func(ai humanAI)Results(hands [][]deck.Card,dealer []deck.Card){
	fmt.Println("##### FINAL HANDS #####")
	for _,h := range hands{
		fmt.Println(" ",h)
	}
	fmt.Println("Dealer: ",dealer)
}

