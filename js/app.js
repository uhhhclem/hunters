var huntersApp = angular.module('huntersApp', ['ngMaterial']);

huntersApp.config(function($mdThemingProvider) {
      $mdThemingProvider.theme('default')
        .primaryPalette('blue')
        .accentPalette('light-blue');
    }
);

var Ctrl = function($http) {
    this.http_ = $http;
    this.restart();
}

Ctrl.prototype.restart = function() {
    this.board = null;
    this.cfg = null;
    this.status = [];
    this.prompt = null;
    this.newGame();
}

Ctrl.prototype.newGame = function() {
    var self = this;
    this.http_.get('/api/newGame').success(function(d){
        self.cfg = {params: {ID: d.ID}};
        self.getBoard();
        self.getStatus();
        self.getPrompt();
    });
}

Ctrl.prototype.getBoard = function() {
    var self = this;
    self.http_.get('/api/board', self.cfg).success(function(result) {
        if (result.End) {
            return;
        }
        self.board = result;
        return self.getBoard();
    });
}

Ctrl.prototype.getStatus = function() {
    var self = this;
    self.http_.get('/api/status', self.cfg).success(function(result){
        if (result.End) {
            return;
        }
        self.status.push(result.Status);
        return self.getStatus();
    });
}

Ctrl.prototype.getPrompt= function() {
    var self = this;
    self.http_.get('/api/prompt', self.cfg).success(function(result){
        if (result.End) {
            return;
        }
        self.prompt = result.Prompt;
        return self.getPrompt();
    });
}

Ctrl.prototype.makeChoice = function(key) {
    var self = this;
    this.http_.post('/api/choice', {ID: self.cfg.params.ID, Key: key})
        .success(function(d){
        })
        .error(function(d){
            self.status.push('Choice failed: ' + d);
        });
}

huntersApp.controller('huntersCtrl', Ctrl);