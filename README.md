# k8s-gitlab-runner-exporter
Экспортер позволяет получать метрики gitlab runners запущенных в кластере kubernetes в формате Prometheus

Приложение может быть запущено с помощью docker образа "dru90i/k8s-gitlab-runner-exporter" или можно собрать образ из Dockerfile

Приложение опционально принимает две переменных окружения:

NAMESPACE_RUNNER (по-умолчанию default) // Неймспейс с gitlab раннером

LABEL_RUNNER (по-умолчанию app=gitlab-runner) // Лейбл по которому определяется под раннера

После запуска под с экспортером будет отдавать на порту 9191 следующие страницы:
1. /runners - список названий подов с раннерами в виде "Name": <имя пода> в формате json
2. /metrics - на вход должен получать параметр runner с именем пода. Пример. http://127.0.0.0.1/metrics?runner=gitlab-runner-123 

Для того чтобы экспортер мог получать метрики из раннеров в их конфиграции должен быть включен HTTP-сервер метрик https://docs.gitlab.com/runner/monitoring/
