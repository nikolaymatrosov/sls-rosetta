request "req" {
  method  = "POST"
  uri     = "https://functions.yandexcloud.net/d4eoe7ndffqrgitograp"
  headers = {
    Content-Type = "application/json"
    Useragent    = "Loadtest"
    X-Request-ID = "{{ uuid }}"
  }
  body             = <<EOF
{"name": "{{ randString }}"}
EOF
}


scenario "test_scenario" {
  requests         = [
    "req",
  ]
}
