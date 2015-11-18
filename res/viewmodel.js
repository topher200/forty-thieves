function CardGameViewModel(data)  {
    var self = this;

    self.cards = ko.observableArray(data.cards);
    self.player1Cards = ko.observableArray(data.player1Cards);
    self.player2Cards = ko.observableArray(data.player2Cards);
    self.stock = ko.observableArray();
    // self.stock([
    //     {"Face":11,"Suit":0},{"Face":7,"Suit":3},{"Face":7,"Suit":1},{"Face":2,"Suit":3},{"Face":10,"Suit":0},{"Face":6,"Suit":1},{"Face":8,"Suit":2},{"Face":8,"Suit":0},{"Face":7,"Suit":0},{"Face":2,"Suit":0},{"Face":5,"Suit":0},{"Face":3,"Suit":3},{"Face":5,"Suit":2},{"Face":4,"Suit":1},{"Face":7,"Suit":2},{"Face":2,"Suit":0},{"Face":4,"Suit":0},{"Face":4,"Suit":1},{"Face":5,"Suit":3},{"Face":8,"Suit":0},{"Face":11,"Suit":2},{"Face":0,"Suit":1},{"Face":0,"Suit":2},{"Face":12,"Suit":3},{"Face":5,"Suit":1},{"Face":12,"Suit":0},{"Face":1,"Suit":1},{"Face":4,"Suit":3},{"Face":3,"Suit":0},{"Face":9,"Suit":2},{"Face":12,"Suit":0},{"Face":6,"Suit":0},{"Face":9,"Suit":0},{"Face":10,"Suit":2},{"Face":4,"Suit":2},{"Face":3,"Suit":2},{"Face":1,"Suit":2},{"Face":10,"Suit":2},{"Face":6,"Suit":3},{"Face":8,"Suit":1},{"Face":12,"Suit":2},{"Face":8,"Suit":3},{"Face":12,"Suit":1},{"Face":12,"Suit":2},{"Face":6,"Suit":0},{"Face":8,"Suit":1},{"Face":3,"Suit":0},{"Face":6,"Suit":2},{"Face":5,"Suit":2},{"Face":5,"Suit":1},{"Face":9,"Suit":1},{"Face":2,"Suit":1},{"Face":1,"Suit":0},{"Face":5,"Suit":0},{"Face":4,"Suit":0},{"Face":7,"Suit":3},{"Face":10,"Suit":3},{"Face":0,"Suit":3},{"Face":1,"Suit":2},{"Face":0,"Suit":3},{"Face":2,"Suit":2},{"Face":11,"Suit":3},{"Face":9,"Suit":0},{"Face":6,"Suit":2}
    // ]);

    $.getJSON("/state", function(state) {
        self.stock([
            {"Face":11,"Suit":0},{"Face":7,"Suit":3},{"Face":7,"Suit":1},{"Face":2,"Suit":3},{"Face":10,"Suit":0},{"Face":6,"Suit":1},{"Face":8,"Suit":2},{"Face":8,"Suit":0},{"Face":7,"Suit":0},{"Face":2,"Suit":0},{"Face":5,"Suit":0},{"Face":3,"Suit":3},{"Face":5,"Suit":2},{"Face":4,"Suit":1},{"Face":7,"Suit":2},{"Face":2,"Suit":0},{"Face":4,"Suit":0},{"Face":4,"Suit":1},{"Face":5,"Suit":3},{"Face":8,"Suit":0},{"Face":11,"Suit":2},{"Face":0,"Suit":1},{"Face":0,"Suit":2},{"Face":12,"Suit":3},{"Face":5,"Suit":1},{"Face":12,"Suit":0},{"Face":1,"Suit":1},{"Face":4,"Suit":3},{"Face":3,"Suit":0},{"Face":9,"Suit":2},{"Face":12,"Suit":0},{"Face":6,"Suit":0},{"Face":9,"Suit":0},{"Face":10,"Suit":2},{"Face":4,"Suit":2},{"Face":3,"Suit":2},{"Face":1,"Suit":2},{"Face":10,"Suit":2},{"Face":6,"Suit":3},{"Face":8,"Suit":1},{"Face":12,"Suit":2},{"Face":8,"Suit":3},{"Face":12,"Suit":1},{"Face":12,"Suit":2},{"Face":6,"Suit":0},{"Face":8,"Suit":1},{"Face":3,"Suit":0},{"Face":6,"Suit":2},{"Face":5,"Suit":2},{"Face":5,"Suit":1},{"Face":9,"Suit":1},{"Face":2,"Suit":1},{"Face":1,"Suit":0},{"Face":5,"Suit":0},{"Face":4,"Suit":0},{"Face":7,"Suit":3},{"Face":10,"Suit":3},{"Face":0,"Suit":3},{"Face":1,"Suit":2},{"Face":0,"Suit":3},{"Face":2,"Suit":2},{"Face":11,"Suit":3},{"Face":9,"Suit":0},{"Face":6,"Suit":2}
        ]);
    });

    self.player1Points = ko.computed(function() {
        var points = 0;
        self.player1Cards().forEach(function(card) {
            points = points + app.cards[card.suit][card.rank];
        });
        return points;
    });

    self.player2Points = ko.computed(function() {
        var points = 0;
        self.player2Cards().forEach(function(card) {
            points = points + app.cards[card.suit][card.rank];
        });
        return points;
    });

    self.remainingInDeck = ko.computed(function() {
        return self.cards().length;
    });

    self.deal = function(hand, num) {
        if(num > self.cards().length) {
            num = self.cards().length;
        }
        for(var i = 0; i < num; i++) {
            var card = self.cards.pop();
            hand.push(card);
        }
    };

    self.newGame = function() {
        while(self.player1Cards().length > 0) {
            self.cards.push(self.player1Cards.pop());
        }
        while(self.player2Cards().length > 0) {
			      self.cards.push(self.player2Cards.pop());
		    }
		    app.shuffle(self.cards());
	  };
}

var gameData = {
    cards: app.shuffle(app.makeDeck()),
    player1Cards: [],
    player2Cards: []
};

var model = new CardGameViewModel(gameData);
ko.applyBindings(model);
