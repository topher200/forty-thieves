function CardGameViewModel()  {
    var self = this;
    self.stock = ko.observableArray();
    self.foundations = ko.observableArray();
    self.tableaus = ko.observableArray();
    self.waste = ko.observableArray();

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

    // Send a move request. Use the response to update cards
    self.movePost = function(fromLocation, fromIndex, toLocation, toIndex) {
        $.post(
            "/move",
            { "FromLocation": fromLocation,
              "FromIndex": fromIndex,
              "ToLocation": toLocation,
              "ToIndex": toIndex},
            self.updateGamestate, "json");
    };

    self.updateGamestate = function(gamestate) {
        self.stock(gamestate.Stock.Cards);
        self.foundations(gamestate.Foundations);
        self.tableaus(gamestate.Tableaus);
        self.waste(gamestate.Waste.Cards);
    };

    // Update gamestate on load
    $.getJSON("/state", self.updateGamestate);
}

var viewmodel = new CardGameViewModel();
ko.applyBindings(viewmodel);

function allowDrop(event) {
    event.preventDefault();
}

function cardDrag(event) {
    event.dataTransfer.setData(
        "FromLocation", $(event.target.parentElement).attr("data-location"));
    event.dataTransfer.setData(
        "FromIndex", $(event.target.parentElement).attr("data-index"));
}

function cardDrop(event) {
    event.preventDefault();
    viewmodel.movePost(
        event.dataTransfer.getData("FromLocation"),
        event.dataTransfer.getData("FromIndex"),
        $(event.target.parentElement).attr("data-location"),
        $(event.target.parentElement).attr("data-index")
    );
}
