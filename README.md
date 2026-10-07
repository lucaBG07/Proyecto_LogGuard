Proyecto_LogGuard
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

FASE 1

1. ¿Qué es un "Log"?
En la industria, cualquier aplicación (servidor web, base de datos, backend de un banco) genera un historial de texto mientras funciona. Ese historial son los logs.

Cada línea de log suele tener un formato estructurado como este:

2026-09-19 14:30:10 [INFO] Usuario id=402 inició sesión exitosamente
2026-09-19 14:30:12 [WARN] Conexión a la base de datos tardó más de 500ms
2026-09-19 14:30:15 [ERROR] No se pudo procesar el pago del usuario id=402: Conexión rechazada

Cuando un servidor corre 24/7, este archivo .log crece a miles de líneas por minuto. Ningún ser humano puede estar abriendo el archivo de texto y leyendo línea por línea para ver si la aplicación está fallando.

2. ¿Qué problema resuelve LogGuard?
LogGuard es una herramienta de monitoreo pasivo en tiempo real. Va a actuar como un "guardia de seguridad" que se queda parado al lado del archivo de logs mirando todo lo que entra.

Su trabajo se divide en 3 tareas principales:

Estar al acecho (Streaming): No abre el archivo, lee todo y se cierra. Se queda "escuchando" el archivo de forma continua. Si una aplicación escribe una nueva línea en el log, LogGuard la detecta al instante.

Filtrar y Parsear: Lee la cadena de texto cruda (string) y la transforma en una estructura con datos claros en memoria (fecha, nivel de gravedad y mensaje).

Analizar y Alertar: Lleva la cuenta de cuántos [ERROR] o [CRITICAL] ocurren en un período corto (por ejemplo, en los últimos 30 segundos). Si detecta un pico inusual de fallos (un umbral que vos definas), dispara una alerta en pantalla para avisar al equipo técnico.

3. ¿Por qué este proyecto destaca en tu GitHub?
Muchos proyectos de estudiantes son estáticos (ej. "calculadora", "lista de tareas"). LogGuard resuelve un problema técnico real usando conceptos avanzados del backend:

Concurrencia con Goroutines y Channels: En lugar de congelar el programa mientras lee el archivo, Go nos permite tener un hilo de ejecución leyendo el disco y otro procesando datos al mismo tiempo sin bloquearse.

Eficiencia de Memoria: No Podés cargar un archivo de logs de 10 GB entero en la memoria RAM. Vas a aprender a procesar flujos de datos (streams) línea por línea consumiendo una cantidad mínima de memoria.

Modelado de Datos: Vas a convertir texto sin formato en estructuras de datos (structs) bien organizadas.

____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

LOS NIVELES DE LOG

los niveles de log (o log levels) son una forma de clasificar que tan grave o importante es el mensaje que la aplicacion acaba de resgistrar.

Por ejemplo: estoy manejando un auto. No es lo mismo que se prenda la luz de "puerta abierta" a que se prenda la luz de "motor sobrecalentado".
Los niveles en software son exactamente esa escala de prioridad.

LA ESCALA ESTANDAR EN LA INDUSTRIA

Van de lo mas cotidiano hasta lo mas destructivo:

1. INFO (Informacion): Todo funciona normal. Son eventos de rutina.
- Ejemplo: "El usuario 'Juan' inicio sesion", "El servidor se inicio en el puerto 8080".

2. WARN (Warning/Advertencia): Algo raro o sospechoso paso, pero el programa pudo seguir funcionando. Es una llamada de atencion para revisar mas tarde.
- Ejemplo: "La base de datos tardo 3 segundos en responder (esta lenta)", "Intento de login con contraseña incorrecta".

3. ERROR (Error): Una funcionalidad concreta fallo y no se pudo completar, pero el resto de la aplicacion sigue viva.
- Ejemplo: "No se pudo procesar el pago con tarjeta", "No se encontro el archivo que el usuario queria descargar".

4. CRITICAL o FATAL (Critico): Fallo algo tan grave que el sistema entero se cayo o se volvio inusable.
- Ejemplo: "Se perdio la conexion total con la base de datos", "No hay espacio disponible en el disco duro".

¿PARA QUE NOS SIRVE EN LogGuard?

Si un servidor escribe 10.000 lineas de log por minuto, la mayoria van a ser de tipo INFO. A nuestro programa no le importa cuando todo anda bien.

LogGuard va a usar estos niveles para filtrar: va a ignorar los mensajes INFO y va a prestarle atencion solo a los ERROR y CRITICAL. Si detecta que
en los ultimos 10 segundos aparecieron 5 niveles ERROR seguidos, sabra que tiene que disparar una alerta inmediatamente.

____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

EXPLICACION STRUCT LogEntry

LogEntry es el struct que va a representar una linea individual de log dentro de la memoria de nuestro programa.

Es comno el "formulario" o la "ficha" donde vas a guardar los datos limpios y ordenados de cada mensaje que lea LogGuard.

DE TEXTO SUCIO A UN LogEntry ORDENADO

En el archivo de texto en disco, el log es solo una cadena de caracteres plana (un string sin formato para la computadora):

"2026-09-20 11:30:00 [ERROR] fallo la conexion a la base de datos"

Para que tu programa en Go pueda entender y trabajar con ese log, necesitamos separar ese texto en partes y guardarlo en una 
estructura limpia. Ahi entra LogEntry.

type LogEntry struct {
    Timestamp time.Time // Para saber EXACTAMENTE cuándo ocurrió (fecha y hora)
    Level     LogLevel  // Para saber QUÉ TAN GRAVE es (INFO, WARN, ERROR, CRITICAL)
    Message   string    // El DETALLE de lo que pasó ("fallo la conexion...")
}

¿POR QUE NECESITAMOS EL LogEntry?

Porque una vez que transformas el texto plano en un LogEntry, podes hacer preguntas logicas en Go super facil:
- Filtrar: if miLog.Level == LevelError (¿Es un error?).
- Calcular tiempo: time.Since(miLog.Timestamp) (¿Paso hace menos de 10 segundos?)
- Analizar texto: strings.Contains(miLog.Menssage, "base de datos") (¿Falla la base de datos?).

En resumen: LogEntry es el molde en Go que convierte una línea de texto suelta en un objeto estructurado con fecha, nivel y mensaje.

____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

EXPLICACION FUNCION ParseLogLine

Ahora que tenemos nuestro molde listo, necesitamos una funcion que tome una linea de texto cruda y nos devuelva un LogEntry 
cargado con sus datos.

func ParseLogLine(line string) (LogEntry, error) {}

PENSANDO LA LOGICA DEL TROCEADO (Slicing)

Supongamos que la funcion recibe este parametro line:
"2026-09-20 11:30:00 [ERROR] fallo la conexion a la base de datos"

Para desarmar este texto en las 3 partes del struct, necesitamos identificar los puntos de corte:

1. Los crochetes [ y ] delimitan el nivel del log.

Para encontrar las posiciones de esos corchetes en el texto, el paquete estandar strings nos da la funcion strings.Index(s, substr).

openBracket := strings.Index(line, "[")  // Devuelve la posición entera de '['
closeBracket := strings.Index(line, "]") // Devuelve la posición entera de ']'

¿COMO FUNCIONA LOS RECORTES (Slicing) CON LOS INDICES?

Si ya encontramos dónde están los corchetes usando strings.Index:
- openBracket es la posición de [
- closeBracket es la posición de ]

Podemos extraer cada segmento de texto directamente así:
1. Fecha/Hora: Está desde el principio (0) hasta justo antes del corchete que abre (openBracket).
-En Go se escribe: line[:openBracket]

2. Nivel: Está desde la posición siguiente al corchete que abre (openBracket + 1) hasta justo antes del corchete que cierra (closeBracket).
- En Go se escribe: line[openBracket+1 : closeBracket]

3. Mensaje: Está desde la posición siguiente al corchete que cierra (closeBracket + 1) hasta el final del texto.
-En Go se escribe: line[closeBracket+1:]

Un detalle útil: Limpiar espacios extras
A veces al cortar nos quedan espacios al principio o al final de las palabras (por ejemplo, " 2026-09-20 11:30:00 ").
Para remover esos espacios sobrantes de cualquier string, Go nos da la función:

strings.TrimSpace(texto)

VALIDACION DE CORCHETES (EVITAR UN CRASH)
¿Qué pasa si por error llega una línea de texto que no tiene corchetes? (Por ejemplo: "2026-09-20 mensaje suelto").

La función strings.Index(texto, "[") devolverá -1 porque no encontró el caracter.
Si intentás hacer texto[:-1], Go va a tirar un panic (el programa se rompe en ejecución) porque no existen los índices negativos en los slices.

Para solucionarlo, simplemente agregamos un control rápido antes de cortar:

if openBracket == -1 || closeBracket == -1 || closeBracket <= openBracket {
    return LogEntry{}, errors.New("formato invalido: faltan corchetes")
}
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

FASE 2: Lectura de Archivos y Streaming (Tiempo Real)

Ahora que ya podemos procesar una linea, el siguiente paso es escalar para procesar un archivo de logs completo e incluso
simular que leemos lineas nuevas en tiempo real.

EL PROBLEMA EN LA VIDA REAL:

Si un archivo de logs pesa 5 GB, no podés cargarlo completo en memoria con una función tipo os.ReadFile. Si hacés eso 
en producción, el servidor se queda sin RAM y el sistema colapsa.

LA SOLUCION EN GO: bufio.Scanner

En Go vamos a usar un paquete llamado bufio (Buffered I/O). El Scanner nos permite leer un archivo linea por linea mediante
un buffer. Carga una linea en memoria, la procesamos, la descarta y pasa a la siguiente. El consumo de memoria es practicamente
cero independientemente del tamaño del archivo.

CONCEPTOS CLAVE QUE VAMOS A USAR EN LA FASE 2:

1. os.Open(path): Abre un archivo del sistema en modo lectura y nos devuelve un puntero al archivo (*os.File).

2. defer file.Close(): La sentencia defer en Go sirve para diferir la ejecucion de una funcion hasta que la funcion actual termine.
Garantiza que el archivo se cierre siempre, incluso si ocurre un error en el medio.

3. bufio.NewScanner(file): Crea el iterador sobre el archivo.

4. scanner.Scan(): Un ciclo for scanner.Scan() avanza linea por linea. Devuelve true si hay una nueva linea o false cuando llega al 
final del archivo (EOF-End of File).

5. scanner.Text(): Devuelve el string de la linea actual.
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

FASE 3: FILTRADO Y METRICAS DE ALERTA

Ahora que ya leemos archivos en streaming, vamos a darle a LogGuard la inteligencia para detectar anomalias o eventos criticos.

En un entorno real, un sysadmin no quiere leer millones de logs INFO: quiere enterarse unicamente cuando ocurren advertencias o errores.

FILTRAR POR NIVEL DE SEVERIDAD

Vamos a modificar nuestar logica para procesar o contar unicamente las entradas de un nivel especifico (por ejemplo, filtrar solo WARN,
ERROR O CRITICAL).

Como struct LogEntry tiene un campo Level del tipo LogLevel, podemos hacer comparaciones directas en Go:

if entry.Level == LevelError || entry.Level == LevelCritical {
    // Es una entrada grave
}

ACUMULAR ESTADISTICAS (CONTADORES)

Para tener un resumen general del estado del sistema al finalizar el escaneo del archivo, podemos usar un mapa(map) en Go para contar
cuantos logs de cada tipo aparecieron.

En Go, un mapa para contar niveles se declara así:

// Mapa donde la clave es el LogLevel y el valor es un entero (contador)
conteo := make(map[LogLevel]int)

// Para incrementar un contador:
conteo[entry.Level]++

LA MAGIA DEL MAP Y PORQUE NO NECESITAMOS UN SWITCH O UN IF ELSE

como conteo es de tipo map[LogLevel]int, podemos usar la variable entry.Level directamente como clave dentro de los corchetes.
Reemplazando todo el switch por una sola linea, Go se encarga automaticamente de buscar el nivel que corresponda (INFO, WARN, 
ERROR, etc) e incrementar su contador:
conteo[entry.Level]++
Si la clave no existe todavía en el mapa, Go la inicializa en 0 y le suma 1.
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

FASE 4: ALERTAS EN TIEMPO REAL Y METRICAS AVANZADAS

En el mundo real, un sistema de monitoreo coomo LogGuard no solo cuenta logs al final; tiene que reaccionar inmediatamente si
ocurre algo grave.

Vamos a agregar una regla de alerta simple:

Regla de Alerta: Si un log tiene nivel ERROR o CRITICAL, queremos emitir una alerta visual inmediata en la consola con un formato
destacado y la fecha/hora en que ocurrio.

¿COMO ESTRUCTURAR LA ALERTA?

Dentro del bucle for scan.Scan(), despues de incrementar el contador en el mapa, podemos verificar si el nivel requiere atencion
inmediata:

if entry.Level == LevelError || entry.Level == LevelCritical {
    fmt.Printf("🚨 [ALERTA] Evento grave detectado a las %s: %s\n", 
        entry.Timestamp.Format("15:04:05"), 
        entry.Message,
    )
}
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

FASE 5: EL TOQUE PROFESIONAL (CONCURRENCIA CON GOROUTINES Y CHANNELS)

EL ESCENARIO REAL:

Actualmente, todo pasa en un solo hilo secuencial: el programa lee la linea, la parsea, analiza si es alerta y recien 
despues pasa a la siguiente.

En un sistema que recibe miles de logs por segundo, queremos separar el trabajo:

1. Un proceso (Productor): lee el archivo lo mas rapido posible y manda las lineas crudas por un canal.

2. Otro proceso (Consummidor/Worker): recibe esas lineas del canal, las parsea y dispara las alertas/metricas de forma
independiente.

LOS 2 CONCEPTOS CLAVE DE GO PARA ESTO:

1. goroutine (Sintaxis: go miFuncion()): Es un hilo de ejecucion liviano administrado por Go. Te permite ejecutar funciones 
en segundo plano sin bloquear el prograna principal.

2. channel (Sintaxis: chan string): Es la "tuberia" por donde las goroutines se pasan pasan datos de forma segura sin tener
problemas de memoria compartida (thread safety).
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

CONCEPTOS QUE NO SABIA

¿QUE ES UN CHANNEL?

Un canal es como una cinta transportadora de fabrica:
- En un extrema hay un productor poniendo objetos en la cinta (canal <- dato).
- En el otro extremo hay un consumidor sacando objetos de la cinta (dato := <-canal).
- Sincronización automática: Si la cinta está vacía, el consumidor se frena a esperar a que caiga algo nuevo. Si el canal 
no tiene buffer, el productor no puede poner otro dato hasta que el consumidor retire el anterior.

¿COMO SE CREA UN CANAL EN GO?
// Creamos un canal por donde solo van a viajar strings
lineasChan := make(chan string)

ESCRIBIR EN UN CANAL (PRODUCTOR)

Para meter datos dentro del canal se usa la flecha mirando hacia la variable del canal: canal <- valor.
Por ejemplo, adentro de el bucle for scan.Scan():

lineaTexto := scan.Text()
lineasChan <- lineaTexto // Envía la línea al canal y espera a que alguien la reciba

REGLA DE ORO DE LOS CANALES EN GO:

Cuando el productor termina de leer el archivo, debe cerrar el canal con close(lineasChan).
Si no cerramos el canal, el consumidor se queda esperando datos para siempre y Go lanza un 
error de colapso (deadlock)

LEER DE UN CANAL CON for ... range (CONSUMIDOR)

Para consumir los datos del canal hasta que se cierre, Go te da una sintaxis con range:

// Este bucle recibe líneas del canal automáticamente.
// Cuando el canal se cierra con close(), el bucle termina solo.
for linea := range lineasChan {
    // Acá procesamos la línea
}

COMO SON LAS GOROUTINES

Si el canal es la cinta transportadora por donde viajan las cajas de logs, la goroutine es el empleado o 
robot de la fabrica que se pone al lado de la cinta a procesar esas cajas.

- Sin goroutine (Secuencial): Sería como tener un solo empleado que tiene que ir a buscar la caja a la 
puerta, volver, leerla, procesarla y recién ahí ir a buscar la siguiente. Si el archivo tiene 100.000 
líneas, el programa se queda congelado haciendo eso de uno en uno.

- Con goroutine (Concurrente): Ponés a un empleado (el productor) a tirar líneas a la cinta, y le 
decís a otro empleado (la goroutine): "Ponete al final de la cinta a procesar todo lo que vaya 
cayendo en paralelo". Mientras el primer empleado sigue leyendo el archivo, el segundo va analizando 
los logs al mismo tiempo.

¿COMO SE CREA UNA GOROUTINE EN GO?
Es increíblemente fácil. Lo único que tenés que hacer es anteponer la palabra clave go antes de llamar a 
cualquier función.

Imaginá que tenés tu función consumidora:

func ConsumeLogs(lineasChan chan string) {
    for linea := range lineasChan {
        // Acá parseamos y analizamos cada línea que viene de la cinta
    }
}

Para poner a ese empleado a trabajar en segundo plano sin bloquear tu programa principal, simplemente 
hacés esto en el main:

// Lanzamos al consumidor como una goroutine independiente
go ConsumeLogs(lineasChan)

Al poner la palabra go, Go dice: "Listo, abro un hilo liviano en segundo plano para que corra esta 
función, y yo sigo ejecutando la línea de código que sigue inmediatamente".

EL GRAN PELIGRO DE LAS GOROUTINES (EL "PORTAZO" del main)

Acá hay un detalle vital que te va a pasar la primera vez que lo pruebes:

El programa principal (main) es el jefe de la fábrica. Si el jefe termina de leer el archivo, llega 
al final de la función main y apaga la luz (return), todo el programa se cierra al instante, sin importar 
si el empleado de la goroutine seguía procesando datos a mitad de camino.

Para evitar que el main cierre la fábrica antes de que el consumidor termine de vaciar el canal, 
necesitamos una forma de sincronizar (esperar a que el consumidor termine) usando
una herramienta nativa de Go llamada sync.WaitGroup (el estándar de la industria para esperas de goroutines).

EL CONTADOR DE EMPLEADOS (sync.WaitGroup)

Imaginate que el jefe (main) necesita irse a su casa, pero antes quiere asegurarse de que el robot 
limpiador/consumidor terminó de procesar todas las cajas de la cinta.

El sync.WaitGroup es literalmente un contador de empleados activos:

1. wg.Add(1) (Anotar un empleado): Antes de encender la goroutine, el jefe le avisa al contador: "Ojo, 
acabo de mandar a un empleado a trabajar, sumalo a la lista".

2. wg.Done() (Avisar que terminó): Cuando el empleado termina de vaciar el canal por completo, presiona 
un botón que le avisa al contador: "Listo, ya terminé mi tarea".

3. wg.Wait() (Esperar en la puerta): El jefe se planta en la puerta y dice: "De acá no me muevo hasta 
que el contador llegue a cero".
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________

ProcessLogFile (antiguo/secuencial)

 func ProcessLogFile(filePath string) error {
	archivo, err := os.Open(filePath)
	if err != nil {
		return errors.New("Hubo un problema en la apertura del archivo")
	}
	defer archivo.Close()
	scan := bufio.NewScanner(archivo)

	conteo := make(map[LogLevel]int)
	for scan.Scan() {
		lineaTexto := scan.Text()
		entry, err := ParseLogLine(lineaTexto)
		if err != nil {
			fmt.Println("Omitiendo linea invalida:", err)
			continue // Salta a la siguiente iteracion for
		}
		conteo[entry.Level]++
		if entry.Level == LevelError || entry.Level == LevelCritical {
			fmt.Printf("[ALERTA] Evento grave detectado a las %s: %s\n", entry.Fecha.Format("15:04:05"), entry.Mensaje)
		}
		// fmt.Printf("%+v\n", entry)
	}

	if err := scan.Err(); err != nil {
		return fmt.Errorf("error al leer el archivo: %w", err)
	}

	fmt.Println("\n=== RESUMEN DE LOGS ===")
	fmt.Printf("INFO     : %d\n", conteo[LevelInfo])
	fmt.Printf("WARN     : %d\n", conteo[LevelWarn])
	fmt.Printf("ERROR    : %d\n", conteo[LevelError])
	fmt.Printf("CRITICAL : %d\n", conteo[LevelCritical])
	return nil
}
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
____________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________________
