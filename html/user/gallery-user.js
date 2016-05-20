
var app = angular.module('gogallery', ['ngRoute', 'flow']);

app.config(function($routeProvider, flowFactoryProvider) {    
    $routeProvider
        .when('/gogallery', {
            templateUrl: '/blaa.html',
        }).otherwise({
            templateUrl: 'mainTemplate.html',
        });

    flowFactoryProvider.defaults = {
        target: '/gallery/api/upload',
        withCredentials: true,
        method: 'octet',
        chunkSize: 50*1024*1024,
        //chunkSize: 200*1024,
        uploadMethod: 'POST',
        simultaneousUploads: 6,
        testChunks: false,
        maxChunkRetries: 10,
        permanentErrors:[404, 500, 501]
    };
});

app.controller("GalleryCtrl", function($scope, $http) {
    $scope.myImages = [];
    $scope.myData = {};
    $scope.editMyData = false;

    $scope.redirectOnAccessError = function(response) {
        if (response.status == 403 || response.status == 401) {
            location.href= "/login";
        } else {
            console.log();
        }
    }

    $http.get("/gallery/api/myImages")
        .then(function(response) {
            $scope.myImages = response.data;
        }, $scope.redirectOnAccessError);

    $http.get("/gallery/api/myData")
        .then(function(response) {
            $scope.myData = response.data;
        },$scope.redirectOnAccessError);

    $scope.startEditMyData = function() {
        $scope.editMyData = true;
        $scope.editMyDataSuccess = undefined;
    }
    
    $scope.saveMyData = function() {
        $http.post("/gallery/api/myData", $scope.myData)
            .then(function(response) {
                $scope.editMyData = false;
                $scope.editMyDataSuccess = true;
            },function(response) {
                $scope.redirectOnAccessError(response);
                $scope.editMyData = false;
                $scope.editMyDataSucess = false;
            });
    }
    
    $scope.deleteImage = function(image) {
        $http.delete("/gallery/api/myImages/"+image.ID)
            .then(function(response) {
                var index = $scope.myImages.indexOf(image);
                if (index > -1) {
                    $scope.myImages.splice(index, 1);
                }
            },$scope.redirectOnAccessError);
    }

    $scope.success = function(file, message) {
        file.done = true;
        file.remoteObject = JSON.parse(message);
        $scope.myImages.unshift(file.remoteObject);
        file.cancel();
    };
    
    $scope.percentDone = function(file) {
        return Math.round(file._prevUploadedSize / file.size * 100).toString() + "%";
    };
    
    $scope.progress = function(file) {
        return {width: $scope.percentDone(file)};
    };
})
