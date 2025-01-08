import os
import posixpath

from gigachat import GigaChat

# Ключ авторизации, полученный в личном кабинете, в проекте GigaChat API.
api_key = os.getenv("GIGACHAT_API_KEY", "").strip('\n')

# Путь к коду функции внутри рантайм-окружения
# https://yandex.cloud/ru/docs/functions/concepts/runtime/environment-variables#files
code_folder = "/function/code"
ca_bundle_file = posixpath.join(code_folder, "russian_trusted_root_ca.cer")


def handler(event, context):
    with GigaChat(
            credentials=api_key,
            ca_bundle_file="russian_trusted_root_ca.cer",
            timeout=300
    ) as giga:
        response = giga.chat("Как использовать GigaChat API?")
        content = response.choices[0].message.content

    return {
        "statusCode": 200,
        "body": content
    }
