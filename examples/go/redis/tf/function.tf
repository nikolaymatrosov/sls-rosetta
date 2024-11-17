data "archive_file" "function_files" {
  output_path = "./function.zip"
  source_dir  = "../function"
  type        = "zip"
}

locals {
  redis_hosts = [for host in yandex_mdb_redis_cluster.demo.host : host.fqdn]
  funcs = {
    "redis" = {
      name    = "plain"
      handler = "index.PlainHandler"
    }
    "pooled" = {
      name    = "pooled"
      handler = "index.PoolHandler"
    }
    "az-detect" = {
      name    = "az-detect"
      handler = "index.AzDetectHandler"
    }
    "check" = {
      name    = "check"
      handler = "index.Handler"
    }
  }
}


resource "yandex_function" "redis_function" {
  for_each          = local.funcs
  name              = each.value.name
  user_hash         = data.archive_file.function_files.output_sha256
  runtime           = "golang121"
  entrypoint        = each.value.handler
  memory            = "128"
  execution_timeout = "10"

  service_account_id = yandex_iam_service_account.handler.id

  content {
    zip_filename = data.archive_file.function_files.output_path
  }

  environment = {
    REDIS_ADDRS = join(",", local.redis_hosts)
    # REDIS_ADDRS = local.redis_hosts[0]
    REDIS_MASTER = yandex_mdb_redis_cluster.demo.name
  }

  secrets {
    environment_variable = "REDIS_PASSWORD"
    id                   = yandex_lockbox_secret.redis.id
    key                  = "password"
    version_id           = yandex_lockbox_secret_version_hashed.password_version.id
  }

  connectivity {
    network_id = data.yandex_vpc_network.default.id
  }
  concurrency = 16

  depends_on = [
    yandex_mdb_redis_cluster.demo,
    yandex_lockbox_secret_version_hashed.password_version,
    yandex_resourcemanager_folder_iam_binding.handler,
  ]
}


// IAM binding for making function public
resource "yandex_function_iam_binding" "test_function_binding" {
  for_each = yandex_function.redis_function
  function_id = each.value.id
  #  role        = "functions.functionInvoker"
  role     = "serverless.functions.invoker"
  members = ["system:allUsers"]
}



