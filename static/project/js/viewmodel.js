function CardGameViewModel()  {
    var self = this;
    self.stock = ko.observableArray();
    self.foundations = ko.observableArray();
    self.tableaus = ko.observableArray();
    self.waste = ko.observableArray();
    self.score = ko.observable();
    self.gameStateID = ko.observable();

    // set our query param, if we have it, before we make any json requests
    var pageUrl = new URL(window.location.href);
    if (pageUrl.searchParams.get("gameStateID")) {
        self.gameStateID(pageUrl.searchParams.get("gameStateID"));
    }

    self.imageFilename = function(card) {
        return 'static/project/cards-png/' + card.Suit + '-' + card.Face + '.png';
    };

    self.tableauPosition = function(index) {
        return "top: " + index * 50 + 'px';
    };

    // Send a newgame post on button click. Update cards with updateGamestate.
    self.newgamePost = function() {
        $.post("/newgame", '{ }', self.updateGamestate, "json");
    };

    // add our query param to a full URL
    function addGameStateIdToURL(urlString) {
        var url = new URL(urlString);
        url.searchParams.set("gameStateID", self.gameStateID());
        return url.href;
    }

    // add our query param to a partial (path only) url
    function addGameStateIdToPath(path) {
        if (path.includes("?")) {
            console.log("warning: request %s already contains query param", path);
        }
        path += "?gameStateID=" + self.gameStateID();
        return path;
    }

    // Send a move request. Use the response to update cards
    self.movePost = function(fromLocation, fromIndex, toLocation, toIndex) {
        $.post(
            addGameStateIdToPath("/move"),
            { "FromPile": fromLocation,
              "FromIndex": fromIndex,
              "ToPile": toLocation,
              "ToIndex": toIndex},
            self.updateGamestate, "json");
    };

    // Send a flip stock request. Use the response to update cards
    self.flipStockPost = function() {
        $.post(addGameStateIdToPath("/flipstock"), {}, self.updateGamestate, "json");
    };

    self.foundationCardPost = function() {
        $.post(addGameStateIdToPath("/foundationcard"), {}, self.updateGamestate, "json");
    };

    self.updateGamestate = function(gamestate) {
        self.stock(gamestate.Stock.Cards);
        self.foundations(gamestate.Foundations);
        self.tableaus(gamestate.Tableaus);
        self.waste(gamestate.Waste.Cards);
        self.score(gamestate.Score);
        self.gameStateID(gamestate.GameStateID);

        // set the url to match our gamestate
        var url = addGameStateIdToURL(window.location.href);
        history.pushState(null, '', url);
    };

    // Update gamestate on load
    $.getJSON(addGameStateIdToPath("/state"), self.updateGamestate);
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

function flipStock() {
    viewmodel.flipStockPost();
}
