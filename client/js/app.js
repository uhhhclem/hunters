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
    this.gameState = null;
    this.status = [];
    this.prompt = null;
}

Ctrl.prototype.getGameState = function() {
    var self = this;
    self.http_.get('/api/dump').success(function(result) {
        self.gameState = result;
    });
}

Ctrl.prototype.getStatus = function() {
    var self = this;
    self.http_.get('/api/status').success(function(result){
        self.status.push(result);
    });
}

Ctrl.prototype.getPrompt= function() {
    var self = this;
    self.http_.get('/api/prompt').success(function(result){
        self.prompt = result;
    });
}


huntersApp.controller('mainCtrl', Ctrl);