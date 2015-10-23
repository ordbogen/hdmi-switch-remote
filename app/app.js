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
            var w = null;
            var connect;
            connect = function(){
                var w;
                if (w === null) {
                    w = new WebSocket("ws://" + location.host + "/socket");
                }
                return w;
            };
            return {
                connect: function(callback){
                    var w = connect();
                    w.onmessage = function(){
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
