#!/bin/sh

set -e

BINARY_NAME="docky"
INSTALL_DIR="/usr/local/bin"

echo "Установка $BINARY_NAME..."
TMP_FILE="$(mktemp)"
curl -sSL https://raw.githubusercontent.com/BkycHblu-6oPwuK/docky/main/bin/docky -o "$TMP_FILE" || {
    echo "Ошибка загрузки docky"
    exit 1
}

chmod +x "$TMP_FILE"
sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
echo "Установка прошла успешно"

# === Установка автодополнения ===
if command -v docky >/dev/null; then
    echo "Настройка автодополнения..."
    sudo mkdir -p /etc/bash_completion.d
    docky completion bash | sudo tee /etc/bash_completion.d/docky >/dev/null
    echo "Автодополнение установлено в /etc/bash_completion.d/docky"
fi

echo "Выполните: exec bash (или перезапустите терминал)"
echo "Команда для очистки кеша: docky clean-cache"
