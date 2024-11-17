request "req" {
  method = "GET"
  uri    = "https://functions.yandexcloud.net/${function_id}?integration=raw"
  headers = {
    Content-Type = "application/json"
    Useragent    = "Loadtest"
    X-Request-ID = "{{ uuid }}"
  }
}


scenario "test_scenario" {
  requests = [
    "req",
  ]
}
