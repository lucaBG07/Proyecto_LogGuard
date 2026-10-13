package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// Definimos un tipo personalizado para los niveles
type LogLevel string

// Definimos las opciones fijas (constantes)
const (
	LevelInfo     LogLevel = "INFO"
	LevelWarn     LogLevel = "WARN"
	LevelError    LogLevel = "ERROR"
	LevelCritical LogLevel = "CRITICAL"
)

// Layout de tiempo estándar en Go para parsear la fecha: YYYY-MM-DD HH:MM:SS
const timeLayout = "2006-01-02 15:04:05"

// Aca definimos el struct de la entrada de log
type LogEntry struct {
	Fecha   time.Time
	Level   LogLevel
	Mensaje string
}

func ParseLogLine(line string) (LogEntry, error) {
	texto := strings.TrimSpace(line)
	openBracket := strings.Index(texto, "[")  // Devuelve la posición entera de '['
	closeBracket := strings.Index(texto, "]") // Devuelve la posición entera de ']'

	if openBracket == -1 || closeBracket == -1 || closeBracket <= openBracket {
		return LogEntry{}, errors.New("formato invalido: faltan corchetes")
	}

	fecha := texto[:openBracket]
	fecha = strings.TrimSpace(fecha)
	fechaParseada, err := time.Parse(timeLayout, fecha)
	if err != nil {
		return LogEntry{}, errors.New("Error al parsear la fecha")
	}

	nivel := LogLevel(texto[openBracket+1 : closeBracket])
	mensaje := texto[closeBracket+1:]

	// return LogEntry{fechaParseada, nivel, mensaje}, nil

	return LogEntry{
		Fecha:   fechaParseada,
		Level:   nivel,
		Mensaje: mensaje,
	}, nil
}

func ProcessLogFile(filePath string) error {
	var wg sync.WaitGroup

	// Creamos el canal
	lineasChan := make(chan string)

	// Anotamos que vamos a lanzar una goroutine
	wg.Add(1)

	// Lanzamos la go routine
	go ConsumeLogs(lineasChan, &wg)

	// Ejecutamos el Productor (lee el archivo y alimenta el canal)
	err := ProduceLogs(filePath, lineasChan)
	if err != nil {
		return err
	}
	// El main se planta en la puerta a esperar que se vacíe el canal
	wg.Wait()

	return nil
}

func ProduceLogs(filePath string, lineasChan chan string) error {
	archivo, err := os.Open(filePath)
	if err != nil {
		return errors.New("Hubo un problema en la apertura del archivo")
	}
	defer archivo.Close()
	scan := bufio.NewScanner(archivo)
	defer close(lineasChan)
	for scan.Scan() {
		lineaTexto := scan.Text()
		lineasChan <- lineaTexto // Envía la línea al canal y espera a que alguien la reciba
	}
	if err := scan.Err(); err != nil {
		return fmt.Errorf("error al leer el archivo: %w", err)
	}

	return nil
}

func ConsumeLogs(lineasChan chan string, wg *sync.WaitGroup) {
	// Esto asegura que al terminar la función se ejecute Done() automáticamente
	defer wg.Done()

	conteo := make(map[LogLevel]int)

	// Bucle que lee de la cinta hasta que se cierre con close()
	for linea := range lineasChan {
		entry, err := ParseLogLine(linea)
		if err != nil {
			fmt.Println("Omitiendo linea invalida:", err)
			continue // Salta a la siguiente iteracion for
		}
		conteo[entry.Level]++
		if entry.Level == LevelError || entry.Level == LevelCritical {
			fmt.Printf("[ALERTA] Evento grave detectado a las %s: %s\n", entry.Fecha.Format("15:04:05"), entry.Mensaje)
		}
	}
	fmt.Println("\n=== RESUMEN DE LOGS ===")
	fmt.Printf("INFO     : %d\n", conteo[LevelInfo])
	fmt.Printf("WARN     : %d\n", conteo[LevelWarn])
	fmt.Printf("ERROR    : %d\n", conteo[LevelError])
	fmt.Printf("CRITICAL : %d\n", conteo[LevelCritical])
}

func main() {
	err := ProcessLogFile("app.log")
	if err != nil {
		fmt.Println("Error procesando archivo:", err)
	}
}
