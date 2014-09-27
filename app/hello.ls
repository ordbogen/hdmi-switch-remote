app = window.angular.module "Hello", [
  "ngResource"
  "ngMaterial"
]

ng = {}

app.factory "Mode", ($http) ->
  mode = {}

  mode.address = "127.0.0.1:23"

  mode.switch = (newMode) ->
    $http.post(
      "/switch-mode"
      {
        mode: newMode
        address: mode.address
      }
    )

  mode

app.factory "Connection", ->
  w = null
  connect = ->
    if !w?
      w = new WebSocket("ws://#{location.host}/socket")
    w
  {
    connect: (callback) ->
      w = connect!
      w.onmessage = callback

    connectJson: (callback) ->
      w = connect!
      w.onmessage = (event) ->
        callback JSON.parse event.data
  }

app.controller "HelloCtrl", ($scope, Mode) ->
  $scope.mode = Mode
