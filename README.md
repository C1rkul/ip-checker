# ip-checker
Simple ip-checker / Простой ip-checker



[ENGLISH/АНГЛИЙСКИЙ]

# IP-Checker 🌐

> A tiny desktop tool that draws a glowing neon arc across a 3D globe — from your real location to any IP you choose (including your VPN exit node).

![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)
![Wails](https://img.shields.io/badge/UI-Wails%20v2-2D2D2D?logo=linux&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Vibe-coded](https://img.shields.io/badge/built%20with-Claude.ai-purple)



## Features

- **3D neon globe** rendered with `globe.gl` / three.js — a dotted Earth in your signature cyan-on-black palette.
- **Animated arc** that slowly flies from the start point to the target, with a glowing travelling dot and a pulsing marker on arrival.
- **Configurable start point** — pin your real-world location, then aim the arc at your own IP *or any other IP / host*. Perfect for visualising a VPN: the arc clearly shows traffic leaving from home and landing on the VPN exit node.
- **Geolocation card** for the target — city, region and country resolved via the free `ipwho.is` API.
- **Glassmorphism + neon dark theme** throughout the interface.
- **No API key required** — uses `ipwho.is` and its endpoints (current IP, plus lookup by arbitrary IP).

## Stack

- **Backend:** Go
- **UI framework:** [Wails v2](https://wails.io) (native window + system WebView)
- **Frontend:** JavaScript / TypeScript, `globe.gl` (three.js)
- **Geo API:** [ipwho.is](https://ipwho.is)

## Installation & run

Prerequisites: **Go 1.23+**, **Node.js 18+**, and on Windows the **WebView2** runtime (preinstalled on Windows 10/11).

```bash
# 1. Install the Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor

# 2. Clone and enter the project
git clone https://github.com/C1RKUL/ip-checker.git
cd ip-checker

# 3. Development mode (hot reload)
wails dev

# 4. Build a native binary
wails build
```

The compiled executable lands in `build/bin/`.

## How the VPN view works

Set your **start point** to your actual location, then hit *Check IP*. The arc animates from home to the exit node — so you can literally see where your traffic is being routed. Leave the target empty to trace your current public IP instead.

## API & limits

`ipwho.is` is free and keyless, but it is still rate-limited — the app caches your start point to keep requests to a minimum. Be kind to the service if you fork it.

## Roadmap (ideas, not implemented yet)

- Multi-hop traceroute with chained arcs
- Lookup history / favourites
- Fallback geo-providers

## About this project

This app was **vibe-coded** — written in Go end-to-end with [Claude.ai](https://claude.ai) as a pair-programmer. The whole UI layer is Wails v2, so there is no Electron bloat and no C compiler headaches: just Go + a webview.

## License

Released under the **MIT License**. See [LICENSE](./LICENSE).

```
MIT License

=======

[RUSSIAN/РУССКИЙ]

# IP-Checker 🌐

> Маленькая десктоп-утилита, которая рисует светящуюся неоновую дугу на 3D-глобусе — от вашей реальной точки до любого IP, который вы укажете (включая выходной узел VPN).

![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8?logo=go&logoColor=white)
![Wails](https://img.shields.io/badge/UI-Wails%20v2-2D2D2D?logo=linux&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)
![Vibe-coded](https://img.shields.io/badge/built%20with-Claude.ai-purple)


## Возможности

- **3D-глобус в неоне** на `globe.gl` / three.js — точечная Земля в фирменной палитре «циан на чёрном».
- **Анимированная дуга**, которая медленно летит от точки старта к цели, со светящейся движущейся точкой и пульсирующим маркером в конце.
- **Настраиваемая точка старта** — задаёте своё реальное местоположение, а дугой целитесь в свой IP *или в любой другой IP / хост*. Отлично показывает VPN: видно, как трафик уходит из дома и приземляется на выходной узел.
- **Карточка геолокации** цели — город, регион и страна через бесплатный API `ipwho.is`.
- **Тёмная тема glassmorphism + неон** во всём интерфейсе.
- **API-ключ не нужен** — используется `ipwho.is` и его эндпоинты (текущий IP + поиск по произвольному IP).

## Стек

- **Бэкенд:** Go
- **UI-фреймворк:** [Wails v2](https://wails.io) (нативное окно + системный WebView)
- **Фронтенд:** JavaScript / TypeScript, `globe.gl` (three.js)
- **Гео-API:** [ipwho.is](https://ipwho.is)

## Установка и запуск

Требуется: **Go 1.23+**, **Node.js 18+**, а на Windows — среда выполнения **WebView2** (уже стоит в Windows 10/11).

```bash
# 1. Ставим CLI Wails
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor

# 2. Клонируем проект
git clone https://github.com/C1RKUL/ip-checker.git
cd ip-checker

# 3. Режим разработки (hot reload)
wails dev

# 4. Сборка нативного бинарника
wails build
```

Готовый исполняемый файл появится в `build/bin/`.

## Как работает VPN-вид

Задаёте **точку старта** = ваше реальное местоположение, жмёте *Check IP*. Дуга анимируется от дома до выходного узла — буквально видно, куда маршрутизируется трафик. Оставьте цель пустой, чтобы провести трассировку до своего текущего публичного IP.

## API и лимиты

`ipwho.is` бесплатен и не требует ключа, но у него есть rate-limit — приложение кэширует вашу точку старта, чтобы свести запросы к минимуму. Если форкаете — будьте к сервису вежливы.

## Планы (идеи, пока не реализовано)

- Многоскачковый traceroute с цепочкой дуг
- История запросов / избранные точки
- Резервные гео-провайдеры

## Об этом проекте

Приложение написано методом **вайбкодинга** — целиком на Go вместе с [Claude.ai](https://claude.ai) в роли пары-программиста. Весь слой UI — на Wails v2, поэтому никаких тяжестей Electron и возни с C-компилятором: только Go + вебвью.

## Лицензия

Проект распространяется под свободной лицензией **MIT**. См. [LICENSE](./LICENSE).

```
MIT License

Copyright (c) 2026 C1RKUL
```

Copyright (c) 2026 C1RKUL
```
