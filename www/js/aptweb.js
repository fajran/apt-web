(function(app) {

app.controller('SearchCtrl', function($scope, $http) {

  $scope.repos = [];
  $scope.dists = [];

  $scope.init = function() {
    $http.get('/api/v1/info')
      .success(function(data) {
        $scope.repos = data.repos;
        $scope.dists = data.dists;

        $scope.repo = data.repos[0];
        $scope.dist = data.dists[0];
      });
  }

  $scope.search = function() {
    var params = {
      d: $scope.dist.id,
      dn: $scope.dist.name,
      pkgs: $scope.pkgs,
    };
    $http.get('/api/v1/deps', {params: params})
      .success(function(data) {
        $scope.result = data;
      });

    _ga('send', 'pageview', '/search?' + jQuery.param(params));
  }
});


})(angular.module('aptweb', []));
