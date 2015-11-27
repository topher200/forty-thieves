function CardGameViewModel()  {
    var self = this;
    self.stock = ko.observableArray();
    self.foundations = ko.observableArray();
    self.tableaus = ko.observableArray();

    $.getJSON("/state", function(state) {
        self.stock(state.Stock.Cards);
        self.foundations(state.Foundations);
        self.tableaus(state.Tableaus);
    });

    self.imageFilename = function (card) {
        return 'project/cards-png/' + card.Suit + '-' + card.Face + '.png';
    };

    self.tableauPosition = function(index) {
        return "top: " + index * 50 + 'px';
    };

    // $.post("/move",
    //        '{ "FromLocation": "tableau", "FromIndex": 0, "ToLocation": "tableau", "ToIndex": 1 }',
    //        function(data) {
    //            console.log(data);
    //        }, "json");

    // $.getJSON("/state", function(state) {
    //     self.stock(state.Stock.Cards);
    //     self.foundations(state.Foundations);
    //     self.tableaus(state.Tableaus);
    // });
}
