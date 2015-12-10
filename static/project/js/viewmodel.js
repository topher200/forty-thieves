function CardGameViewModel()  {
    var self = this;
    self.stock = ko.observableArray();
    self.foundations = ko.observableArray();
    self.tableaus = ko.observableArray();

    self.imageFilename = function(card) {
        return 'project/cards-png/' + card.Suit + '-' + card.Face + '.png';
    };

    self.tableauPosition = function(index) {
        return "top: " + index * 50 + 'px';
    };

    // Send a newgame post on button click. Update cards with updateGamestate.
    self.newgamePost = function() {
        $.post("/newgame", '{ }', self.updateGamestate, "json");
    };

    // Send a dummy move request on button click
    self.movePost = function() {
        $.post(
            "/move",
            '{ "FromLocation": "tableau", "FromIndex": 0, "ToLocation": "tableau", "ToIndex": 1 }',
            self.updateGamestate, "json");
    };

    self.updateGamestate = function(gamestate) {
        self.stock(gamestate.Stock.Cards);
        self.foundations(gamestate.Foundations);
        self.tableaus(gamestate.Tableaus);
    };

    // Update gamestate on load
    $.getJSON("/state", self.updateGamestate);
}
