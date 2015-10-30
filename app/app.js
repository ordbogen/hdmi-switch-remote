(function() {
    angular
        .module("HdmiSwitch", ["ngResource", "ngMaterial"])
        .factory("Mode", function($http){
            var mode = {};
            mode.address = "192.168.1.181:23";
            mode.switch = function(newMode){
                return $http.post("/switch-mode", {
                    mode: newMode,
                    address: mode.address
                });
            };
            return mode;
        })
        .factory("Connection", function($timeout){
            var socket = null;
            var connect;
            connect = function(){
                var socket;
                if (socket === null) {
                    socket = new WebSocket("ws://" + location.host + "/socket");
                }
                return socket;
            };
            return {
                connect: function(callback){
                    var socket = connect();
                    socket.onmessage = function(){
                        var args = [].slice.call(arguments);
                        return $timeout(function(){
                            return callback.apply(null, args);
                        });
                    };
                }
            };
        })
        .controller("HdmiSwitchCtrl", function(Mode, Connection){
            var vm = this;
            vm.mode = Mode;
            vm.messages = "";
            Connection.connect(function(message){
                var d = new Date();
                vm.messages = d.getHours() + ":" + d.getMinutes() + ":" + d.getSeconds() + " " + message.data + "\n" + vm.messages;
            });
            return vm;
        });
})();
