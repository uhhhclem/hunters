var huntersApp = angular.module('huntersApp', ['ngMaterial']);

huntersApp.config(function($mdThemingProvider) {
      $mdThemingProvider.theme('default')
        .primaryPalette('blue')
        .accentPalette('light-blue');
    }
);

var Ctrl = function($http) {
    this.http_ = $http;
    this.gameState = null;
    this.status = [];
    
    this.getStatus();

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
        return self.getStatus();
    });
}

huntersApp.controller('mainCtrl', Ctrl);