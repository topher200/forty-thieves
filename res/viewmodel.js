function CardGameViewModel(data)  {
    var self = this;

    self.cards = ko.observableArray(data.cards);
    self.player1Cards = ko.observableArray(data.player1Cards);
    self.player2Cards = ko.observableArray(data.player2Cards);
    self.stock = ko.observableArray();
    self.foundations = ko.observableArray();
    self.tableaus = ko.observableArray();

    $.getJSON("/state", function(state) {
        self.stock(state.Stock.Cards);
        self.foundations(state.Foundations);
        self.tableaus(state.Tableaus);
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
