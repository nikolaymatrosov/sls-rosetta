request "req" {
  method = "POST"
  uri    = "https://functions.yandexcloud.net/${function_id}?integration=raw"
  headers = {
    Content-Type = "application/json"
    Useragent    = "Loadtest"
    X-Request-ID = "{{ uuid }}"
  }
  body = <<EOF
${body}
EOF
}


scenario "test_scenario" {
  requests = [
    "req",
  ]
}
