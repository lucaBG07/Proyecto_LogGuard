# 🛡️ LogGuard

**LogGuard** es una herramienta de análisis y parseo concurrente de archivos de registro (*log files*) desarrollada en Go. Está diseñada para procesar volúmenes masivos de datos de manera ultra eficiente, aprovechando al máximo la concurrencia nativa de Go mediante **goroutines**, **canales (`channels`)** y **`sync.WaitGroup`**.

---

## 🚀 Características Principales

* **Procesamiento Concurrente (Worker Pool)**: Distribuye la lectura y el análisis de las líneas del log entre múltiples goroutines para reducir drásticamente el tiempo de procesamiento.
* **Seguridad y Sincronización**: Utiliza canales para la comunicación segura entre hilos y `sync.WaitGroup` para garantizar que todas las tareas finalicen correctamente antes de consolidar los resultados.
* **Agregación de Métricas**: Extrae, clasifica y contabiliza patrones clave (como niveles de severidad: `ERROR`, `WARN`, `INFO`, IPs frecuentes o endpoints más consultados).
* **Alto Rendimiento en Memoria**: Minimiza el uso de recursos procesando streams de datos sin duplicar estructuras innecesarias.

---

## 🏗️ Arquitectura y Funcionamiento

El flujo de trabajo de LogGuard se divide en etapas concurrentes:
1. **Lector (Producer)**: Lee el archivo de logs línea por línea de forma eficiente y las distribuye a través de un canal de trabajo.
2. **Workers (Consumers)**: Múltiples goroutines toman líneas del canal en paralelo, aplicando expresiones regulares o lógica de parseo para extraer la información relevante.
3. **Agregador (Collector)**: Recopila los resultados parciales de cada worker a través de canales de salida y genera el reporte consolidado final.

### Ejemplo conceptual de la estructura:
```go
package main

import (
    "fmt"
    "sync"
)

// Estructura de ejemplo para almacenar métricas de los logs
type LogMetrics struct {
    Errors   int
    Warnings int
    Infos    int
    mu       sync.Mutex // Para proteger la concurrencia si se acumula en un struct compartido
}

// (Aquí se implementaría la lógica de los workers y el control con WaitGroup)
```

## 💻 Ejemplo de Uso

Para poner a prueba LogGuard, podés seguir estos simples pasos:

### 1. Crear un archivo de logs de prueba (`app.log`)
Creá un archivo llamado `app.log` en la misma carpeta del proyecto con el siguiente contenido de ejemplo:

```text
2026-09-20 11:30:00 [INFO] Sistema iniciado correctamente en el puerto 8080
2026-09-20 11:30:15 [WARN] Uso de memoria superando el 80%
2026-09-20 11:31:00 [ERROR] Falló la conexión a la base de datos principal
2026-09-20 11:31:05 [INFO] Reintentando conexión con la base de datos...
2026-09-20 11:32:00 [CRITICAL] Pérdida total de comunicación con el cluster de almacenamiento
```

### Ejecutar la herramienta
go run LogGuard.go

### Salida esperada en consola
[ALERTA] Evento grave detectado a las 11:31:00:  fallo la conexion a la base de datos principal
[ALERTA] Evento grave detectado a las 11:32:00:  Perdida total de comunicacion con el cluster de almacenamiento

=== RESUMEN DE LOGS ===
INFO     : 2
WARN     : 1
ERROR    : 1
CRITICAL : 1

