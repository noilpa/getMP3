Установить `youtube-dl` и/или `yt-dlp` 
```sh
brew install youtube-dl
brew install yt-dlp
```

Установить `ffmpeg`
```sh
brew install ffmpeg
```

Создать новую группу в `Telegram`

Добавить туда `@RawDataBot` для получения `chatID`

Зарегистрировать нового бота у `@BotFather` и получить `Token`

Задать переменные окружения
```sh
export BOT_TOKEN=secret
export CHAT_ID=-10000001
export MP3_DIR=/path/to/mp3
```

Запустить скрипт
```sh
go run cmd/cli/main.go https://youtu.be/SOME_VIDEO
```