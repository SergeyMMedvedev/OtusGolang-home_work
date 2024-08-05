#!/bin/bash

# Запуск docker-compose
docker compose up -d

# Идентификатор контейнера, который нужно отслеживать
CONTAINER_NAME="test-client"

# Время ожидания в секундах (5 минут = 300 секунд)
WAIT_TIME=300
INTERVAL=10

# Проверка статуса выхода контейнера
elapsed_time=0
while [ $elapsed_time -le $WAIT_TIME ]; do
  # Получение кода выхода контейнера
  exit_code=$(docker inspect -f '{{.State.ExitCode}}' $CONTAINER_NAME 2>/dev/null)
  container_state=$(docker inspect -f '{{.State.Status}}' $CONTAINER_NAME 2>/dev/null)
  # Если контейнер завершился
  echo "test-client container status $container_state"
  if [ "$container_state" == "exited" ]; then
    if [ "$exit_code" -eq 0 ]; then
      echo "Container $CONTAINER_NAME exited with status 0"
      echo "Tests successfully completed!"
      echo "Tests logs:"
      docker compose logs test-client
      docker compose down
      exit 0
    else
      echo "Container $CONTAINER_NAME exited with status $exit_code"
      exit 1
    fi
  fi

  # Пауза перед следующей проверкой
  sleep $INTERVAL
  elapsed_time=$((elapsed_time + INTERVAL))
done

# Если контейнер не завершился в течение заданного времени
echo "Container $CONTAINER_NAME did not exit within $WAIT_TIME seconds"
exit 1