# mgtb_vshell

`mgtb_vshell` es un bot de Telegram en Go diseñado para funcionar como una capa tipo shell por usuario, con un espacio de trabajo aislado en el sistema de archivos. La idea principal del proyecto es permitir que cada usuario interactúe con un entorno propio desde Telegram para explorar directorios, administrar archivos y, más adelante, ejecutar tareas de procesamiento de video y transferencia de archivos.

## Objetivo del proyecto

Este bot busca convertir Telegram en una interfaz remota para operaciones de sistema y automatización. En su forma actual, ya establece la base para trabajar con un entorno por usuario. En versiones futuras se ampliará para incluir funciones como:

- codificación de video a X265
- descarga de archivos
- subida de archivos a Telegram
- automatización de tareas relacionadas con almacenamiento y procesamiento

## Estado actual

La implementación actual se enfoca en tres ideas principales:

1. Conectar el bot a Telegram usando `BOT_TOKEN`.
2. Crear un espacio persistente por usuario dentro de `roots/`.
3. Exponer comandos básicos de navegación y administración de directorios.

Actualmente existen estas funciones base:

- `ls`: lista archivos y carpetas del directorio actual del usuario.
- `mkdir`: reservado para crear directorios.
- `x265`: reservado para el flujo de compresión de video.
- `upload`: reservado para subir archivos a Telegram.

## Arquitectura general

El proyecto está organizado en capas simples:

- `bot.go` contiene el punto de entrada de la aplicación.
- `core/` concentra la lógica de usuarios, comandos y estado.
- `core/ffmpeg/` agrupa utilidades y funciones relacionadas con procesamiento multimedia.
- `roots/` actúa como espacio raíz de trabajo por usuario.
- `users/` se usa como directorio auxiliar para datos del proyecto.

### Flujo de ejecución

1. El proceso inicia en `bot.go`.
2. Se lee la variable de entorno `BOT_TOKEN`.
3. Se inicializa el cliente de Telegram.
4. Se crean los directorios necesarios para usuarios y entornos raíz.
5. El bot escucha actualizaciones de Telegram.
6. Cada mensaje se enruta a `core.OnUpdate`.
7. Si el mensaje es un comando, se ejecuta el manejador correspondiente.

## Manejo de usuarios

Cada usuario recibe un entorno propio administrado desde `core/usrenv.go`. Ese entorno se guarda en una ruta basada en su identificador de Telegram dentro de `roots/`.

Esto permite separar la información y el contexto de cada usuario, evitando que compartan el mismo espacio de trabajo.

### Capacidades actuales del entorno

- crear directorios
- cambiar entre rutas dentro del entorno propio
- listar archivos y carpetas

## Procesamiento de video

El proyecto ya incluye una base para compresión con FFmpeg en `core/ffmpeg/compress.go`. La función `Compress` prepara una ejecución de `ffmpeg` con el códec `libx265`, lo que muestra claramente la intención de usar el bot como una herramienta de codificación de video desde Telegram.

La idea es que el bot pueda recibir archivos multimedia, procesarlos y devolver el resultado al usuario sin necesidad de salir de la conversación.

## Tecnologías utilizadas

- Go
- Telegram Bot API para Go
- FFmpeg
- go-fuse como dependencia preparada para evolución futura del manejo de archivos

## Requisitos

- Go 1.26.4 o superior
- Una cuenta de Telegram con un bot creado mediante BotFather
- Variable de entorno `BOT_TOKEN`
- `ffmpeg` instalado en el sistema si se va a usar la parte de codificación de video

## Instalación

```bash
git clone https://github.com/VIRUSGAMING64/mgtb_vshell.git
cd mgtb_vshell
go mod download
```

## Configuración

Antes de ejecutar el bot, define el token de Telegram:

```bash
export BOT_TOKEN="tu_token_aqui"
```

## Ejecución

```bash
go run .
```

## Uso previsto

La intención es que el bot funcione como una consola remota en Telegram. En el futuro, el flujo de uso incluirá acciones como:

- enviar un archivo al bot
- elegir una acción sobre ese archivo
- procesarlo dentro del entorno del usuario
- devolver el resultado por Telegram

## Notas técnicas

- El proyecto mantiene un mapa de usuarios en memoria para asociar cada chat con su entorno.
- La implementación actual de comandos está pensada para crecer de forma incremental.
- Algunas funciones todavía están marcadas como pendientes, pero la estructura ya está preparada para extenderlas sin reescribir el núcleo.
