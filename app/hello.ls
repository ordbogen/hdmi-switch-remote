app = window.angular.module "Hello", [
  "ngResource"
  "ngMaterial"
]

ng = {}

app.factory "Mode", ($http) ->
  mode = {}

  mode.address = "192.168.1.181:23"

  mode.switch = (newMode) ->
    $http.post(
      "/switch-mode"
      {
        mode: newMode
        address: mode.address
      }
    )

  mode

app.factory "Connection", ($timeout) ->
  w = null
  connect = ->
    if !w?
      w = new WebSocket("ws://#{location.host}/socket")
    w
  {
    connect: (callback) ->
      w = connect!
      w.onmessage = (...args) ->
        $timeout ->
          callback ...args

    connectJson: (callback) ->
      w = connect!
      w.onmessage = (event) ->
        callback JSON.parse event.data
  }

app.controller "HelloCtrl", ($scope, Mode, Connection) ->
  $scope.mode = Mode

  $scope.messages = ""

  Connection.connect (message) ->
    $scope.messages += message.data + "\n"
