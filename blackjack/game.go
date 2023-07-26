package blackjack

import (
	"errors"

	deck "github.com/gallachaitanya/Deck"
)

type state int8

const ( 
	statePlayerTurn state = iota
	stateDealerTurn
	stateHandOver
)

type Options struct{
	Decks int
	Hands int
	BlackjackPayout float64
}

func New(opts Options)Game{
	g := Game{
		state: statePlayerTurn,
		dealerAI: dealerAI{},
		balance: 0,
	}

	if opts.Decks == 0{
		opts.Decks = 3
	}
	if opts.Hands == 0{
		opts.Hands = 100
	}
	if opts.BlackjackPayout == 0.0{
		opts.BlackjackPayout = 1.5
	}

	g.nDeck = opts.Decks
	g.nHand = opts.Hands
	g.blackjackPayout = opts.BlackjackPayout

	return g
}


type Game struct {
	nDeck int
	nHand int
	deck   []deck.Card
	state  state
	player []hand
	handIdx int
	playerBet int
	dealer []deck.Card
	dealerAI AI
	balance int
	blackjackPayout float64
}

func(g *Game)currentHand() *[]deck.Card{
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handIdx].cards
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("Currently it is no player state")
	}
}

type hand struct{
	cards []deck.Card
	bet int
}

func deal(g *Game){
	playerHand := make([]deck.Card, 0, 5)
	g.handIdx = 0
	g.dealer = make([]deck.Card, 0, 5)
	var card deck.Card
	for i:=0; i<2; i++{
		card, g.deck = draw(g.deck)
		playerHand = append(playerHand, card)
		card, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, card)
	}
	g.player = []hand{
	{
		cards: playerHand,
		bet: g.playerBet,
	},
	}
	g.state = statePlayerTurn
}

func bet(g *Game, ai AI,shuffled bool){
	bet := ai.Bet(shuffled)
	if bet < 100{
		panic("bet must be atleast 100")
	}
	g.playerBet = bet
}

func (g *Game)Play(ai AI)int{
	g.deck = nil
	min := (52 * g.nDeck)/3
	for i:= 0; i< g.nHand; i++{
		shuffled := false
		if len(g.deck) < min{
			g.deck = deck.New(deck.Deck(g.nDeck),deck.Shuffle)
			shuffled = true
		}
		bet(g,ai,shuffled)
		deal(g)
		if BlackJack(g.dealer...){
			endHand(g,ai)
			continue
		}
	for g.state == statePlayerTurn{
		hand := make([]deck.Card,len(*g.currentHand()))
		copy(hand,*g.currentHand())
		move := ai.Play(hand,g.dealer[0])
		err := move(g)
		switch err{
		case errBurst:
			MoveStand(g)
		case nil:
		default:
			panic(err)
		}
	}

	for g.state == stateDealerTurn{
		hand := make([]deck.Card,len(g.dealer))
		copy(hand,g.dealer)
		move:= g.dealerAI.Play(hand,g.dealer[0])
		move(g)
	}
	 endHand(g,ai)
}
return g.balance
}

type Move func(*Game) error

var errBurst = errors.New("hand score exceeded 21")

func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card deck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand,card)
	if Score(*hand...) > 21{
		return errBurst
	}
	return nil
}

func MoveSplit(g *Game) error{
	cards := g.currentHand()
	if len(*cards) != 2 {
		return errors.New("you can only split with two cards in your hand")
	}
	if (*cards)[0].Rank != (*cards)[1].Rank{
		return errors.New("both cards must have same rank to split")
	}
	g.player = append(g.player, hand{
		cards: []deck.Card{(*cards)[1]},
		bet: g.player[g.handIdx].bet,
	})
	g.player[g.handIdx].cards = (*cards)[:1]
	return nil
}

func MoveDouble(g *Game) error{
	if len(*g.currentHand()) != 2{ 
		return errors.New("can only double on a hand with 2 cards")
	}
	return MoveStand(g)
}

func MoveStand(g *Game) error{
	if g.state == stateDealerTurn{
		g.state++
		return nil
	}
	if g.state == statePlayerTurn{
		g.handIdx++
		if g.handIdx >= len(g.player){
			g.state++
		}
		return nil
	}
	return errors.New("invalid game state")
}

func endHand(g *Game, ai AI) {
	dscore := Score(g.dealer...)
	dBlackJack := BlackJack(g.dealer...)
	allHands := make([][]deck.Card,len(g.player))
	for hi, hand := range g.player{
		cards := hand.cards
		allHands[hi] = cards
		pscore, pBlackJack := Score(cards...), BlackJack(cards...)
		
		winnings := hand.bet
		switch{
		case pBlackJack && dBlackJack:
			winnings = 0
		case dBlackJack:
			winnings = -winnings
		case pBlackJack:
			winnings = int(float64(winnings)*g.blackjackPayout)
		case pscore > 21:
			winnings = - winnings
		case dscore > pscore:
			winnings = -winnings
		}
	
		g.balance += winnings
	}
	
	ai.Results(allHands,g.dealer)
	g.player = nil
	g.dealer = nil
}

func draw(cards []deck.Card)(deck.Card,[]deck.Card){
	return cards[0],cards[1:]
}

func Score(hand ...deck.Card)int{
	score := minScore(hand...)
	if score > 10{
		return score
	}
	for _,card := range hand{
		if card.Rank == deck.Ace{
			return score + 10
		}
	}
	return score
}

func Soft(hand ...deck.Card)bool{
	minscore := minScore(hand...)
	actscore := Score(hand...)
	return minscore != actscore
}

func BlackJack(hand ...deck.Card)bool{
	return len(hand) == 2 && Score(hand...) == 21
}

func minScore(hand ...deck.Card)int{
	score := 0
	for _,card := range hand{
		score += min(int(card.Rank),10)
	}
	return score
}

func min(a,b int)int{
	if a < b{
		return a
	}
	return b
}