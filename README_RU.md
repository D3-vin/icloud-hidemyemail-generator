# 📧 iCloud HideMyEmail Generator

Быстрый и эффективный CLI инструмент для генерации iCloud Hide My Email адресов с TLS fingerprinting.

[![Version](https://img.shields.io/badge/version-1.0.0-blue)](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

[English](README.md) | [Русский](#русский)

---

## Возможности

- 🚀 **Быстро** - Производительность Go
- 🔒 **Безопасно** - TLS fingerprinting с профилем Chrome 146
- 📊 **Интерактивное меню** - Простой интерфейс
- 🎯 **Авто-нумерация** - Автоматическая нумерация меток (test1, test2, test3...)
- 📁 **Организованный вывод** - Созданные email → `generated/`, результаты списка → `results/`
- 🌐 **Кросс-платформа** - Windows, macOS (Intel и ARM), Linux
- 📦 **Автономный** - Один файл, без зависимостей

---

## Быстрый старт

### 1. Скачать

Скачайте последний релиз для вашей платформы:
- [Windows (64-bit)](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
- [Linux (64-bit)](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
- [macOS Intel](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)
- [macOS Apple Silicon](https://github.com/D3-vin/icloud-hidemyemail-generator/releases)

Или соберите из исходников:
```bash
git clone https://github.com/D3-vin/icloud-hidemyemail-generator.git
cd icloud-hidemyemail-generator
go build -o hidemyemail .
```

### 2. Извлечь cookies

**Используя расширение Chrome (Рекомендуется):**

1. Установите расширение из папки `cookie-extractor-extension/`
2. Откройте [icloud.com](https://www.icloud.com) и войдите
3. Кликните на иконку расширения 🍪 → Extract → Copy
4. Вставьте в `cookies.txt`

См. [cookie-extractor-extension/README.md](cookie-extractor-extension/README.md) для подробностей.

### 3. Сгенерировать email'ы

**Интерактивное меню:**
```bash
./hidemyemail
```

**CLI команды:**
```bash
# Сгенерировать 5 email с меткой "test"
./hidemyemail generate -l test -c 5

# Список всех email
./hidemyemail list

# Только активные email
./hidemyemail list --active
```

---

## Использование

### Интерактивное меню

Запустите без аргументов:

```bash
./hidemyemail
```

```
╔═══════════════════════════════════════╗
║              Menu                     ║
╠═══════════════════════════════════════╣
║  1. Generate emails                   ║
║  2. List all emails                   ║
║  3. Exit                              ║
╚═══════════════════════════════════════╝

Choose option: 1
Enter label: test
Enter count (1-100): 5
Add number to label? (y/n): y

✓ [1/5] email1@privaterelay.appleid.com (label: test1)
✓ [2/5] email2@privaterelay.appleid.com (label: test2)
...
✓ Saved to generated/emails_test_2026-06-09_01-23-45.txt

Choose option: 2
✓ Saved to results/emails_list.txt and results/emails_full.txt
```

### CLI Команды

**Генерация:**
```bash
./hidemyemail generate -l <метка> -c <количество> [опции]

Опции:
  -l, --label string         Метка для email (обязательно)
  -c, --count int            Количество email, 1-100 (обязательно)
      --cookie-file string   Путь к файлу с cookies (по умолчанию "cookies.txt")
  -o, --output string        Выходной файл (по умолчанию "emails.txt")
      --no-output-file       Не сохранять в файл
```

**Список:**
```bash
./hidemyemail list [опции]

Опции:
      --label-query string   Regex фильтр по метке
      --active               Показать только активные
      --inactive             Показать только неактивные
      --cookie-file string   Путь к файлу с cookies (по умолчанию "cookies.txt")
```

---

## Структура проекта

```
icloud-hidemyemail-generator/
├── hidemyemail              # Бинарник (Linux/macOS)
├── hidemyemail.exe          # Бинарник (Windows)
├── cookies.txt              # Ваши iCloud cookies (создайте этот файл)
├── generated/               # Сгенерированные email с timestamp
├── results/                 # Результаты list (emails_list.txt, emails_full.txt)
├── cookie-extractor-extension/  # Расширение Chrome для cookies
├── cmd/                     # CLI приложение
├── internal/                # Основная логика
│   ├── api/                 # iCloud API клиент
│   ├── config/              # Конфигурация
│   ├── generator/           # Генерация email
│   ├── lister/              # Список email
│   └── output/              # Терминальный UI
└── pkg/models/              # Модели данных
```

---

## Извлечение cookies

### Метод 1: Расширение Chrome (Рекомендуется)

Установите расширение Chrome из папки `cookie-extractor-extension/` для извлечения HttpOnly cookies.

**Почему расширение?** JavaScript скрипты в консоли не могут получить доступ к HttpOnly cookies (`X-APPLE-WEBAUTH-USER`, `X-APPLE-WEBAUTH-TOKEN`) из-за безопасности браузера.

**Установка:**
1. Откройте Chrome → `chrome://extensions/`
2. Включите "Режим разработчика" (вверху справа)
3. Нажмите "Загрузить распакованное расширение"
4. Выберите папку `cookie-extractor-extension/`
5. Закрепите расширение (иконка пазла → закрепить 📌)

**Использование:**
1. Перейдите на [icloud.com](https://www.icloud.com) и войдите
2. Нажмите на иконку расширения 🍪
3. Нажмите "Extract Cookies"
4. Скопируйте строку с cookies
5. Вставьте в `cookies.txt` в корне проекта

См. [cookie-extractor-extension/README.md](cookie-extractor-extension/README.md) для подробностей.

### Метод 2: DevTools Network Tab (Вручную)

1. Откройте [icloud.com](https://www.icloud.com) и войдите
2. Откройте DevTools (F12) → вкладка Network
3. Обновите страницу (F5)
4. Кликните на любой запрос к `icloud.com`
5. Найдите **Headers** → **Request Headers** → **Cookie:**
6. Скопируйте всю строку с cookies
7. Вставьте в `cookies.txt`

**Примечание:** Ручной метод может пропустить HttpOnly cookies. Используйте расширение для лучших результатов.

---

## Сборка

### Одна платформа

```bash
go build -o hidemyemail .
```

### Все платформы

```bash
# Windows
build.cmd

# Или вручную
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o build/hidemyemail-windows-amd64.exe .
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o build/hidemyemail-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o build/hidemyemail-macos-amd64 .
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o build/hidemyemail-macos-arm64 .
```

---

## Решение проблем

### "Cookie file not found"
Создайте `cookies.txt` в корне проекта с вашими iCloud cookies.

### "Authentication failed"
Cookies истекли. Извлеките свежие cookies с icloud.com.

### "failed to extract DSID from cookies"
Отсутствуют HttpOnly cookies. Используйте расширение Chrome вместо консольных скриптов.

### Rate Limit
iCloud ограничивает генерацию до **~5 email каждые 30 минут на члена семьи**. Подождите и попробуйте снова.

---

## Благодарности

- **TLS Client**: [bogdanfinn/tls-client](https://github.com/bogdanfinn/tls-client)
- **CLI Framework**: [spf13/cobra](https://github.com/spf13/cobra)
- **Terminal UI**: [pterm/pterm](https://github.com/pterm/pterm)

---

## Ссылки

- **GitHub**: https://github.com/D3-vin/icloud-hidemyemail-generator
- **Telegram**: [@D3_vin](https://t.me/D3_vin)
- **Автор**: [@D3vin_dev](https://t.me/D3vin_dev)

---

## Лицензия

MIT License - см. файл [LICENSE](LICENSE)

---

**⚠️ Дисклеймер**: Этот инструмент только для образовательных целей. Используйте на свой риск.
