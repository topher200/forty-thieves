var ranks = {
	ACE: 1,
	TWO: 2,
	THREE: 3,
	FOUR: 4,
	FIVE: 5,
	SIX: 6,
	SEVEN: 7,
	EIGHT: 8,
	NINE: 9,
	TEN: 10,
	JACK: 10,
	QUEEN: 10,
	KING: 10
};


var suits = {
	HEARTS: 'hearts',
	CLUBS: 'clubs',
	SPADES: 'spades',
	DIAMONDS: 'diamonds'
};


var app = {

	_ranks : ranks,
	_suits : suits,

	cards : (function() {
		var cards = {};  // suit:ranks
		Object.keys(suits).forEach(function(suit) {
			cards[suit] = ranks;
		});
		return cards;
	}()),

	makeDeck : function() {
		// { suit: "SUIT", rank: "RANK" }
		var deck = [];
		Object.keys(suits).forEach(function(suit){
			Object.keys(ranks).forEach(function(rank) {
				//console.log({
				deck.push({
					 suit: suit
					,rank: rank
				});
			});
		});
		return deck;
	},

	shuffle : function(array) {
	  var currentIndex = array.length, temporaryValue, randomIndex ;

	  // While there remain elements to shuffle...
	  while (0 !== currentIndex) {

	    // Pick a remaining element...
	    randomIndex = Math.floor(Math.random() * currentIndex);
	    currentIndex -= 1;

	    // And swap it with the current element.
	    temporaryValue = array[currentIndex];
	    array[currentIndex] = array[randomIndex];
	    array[randomIndex] = temporaryValue;
	  }

	  return array;
	}
};