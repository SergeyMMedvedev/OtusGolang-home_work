#!/bin/bash
while true; do
  # Ожидаем входящее соединение и читаем сообщение
  echo "Waiting for connection..."
  MESSAGE=$(ncat -l localhost 4242)
  echo "Received: $MESSAGE"

  # Отправляем ответное сообщение
  echo "Hello from server!" | ncat localhost 4242
done