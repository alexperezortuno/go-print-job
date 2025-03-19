package main

import (
	"fmt"
	"log"
	"time"

	"github.com/yusufpapurcu/wmi"
)

// Estructura para almacenar información sobre trabajos de impresión
type PrintJob struct {
	JobID         uint32    `wmi:"JobId"`
	PrinterName   string    `wmi:"PrinterName"`
	TotalPages    uint32    `wmi:"TotalPages"`
	DocumentName  string    `wmi:"Document"`
	JobStatus     string    `wmi:"JobStatus"`
	Owner         string    `wmi:"Owner"`
	TimeSubmitted time.Time `wmi:"TimeSubmitted"`
}

func main() {
	fmt.Println("Monitoreo de impresiones iniciado...")

	for {
		printJobs, err := getPrintJobs()
		if err != nil {
			log.Println("Error obteniendo trabajos de impresión:", err)
		} else {
			for _, job := range printJobs {
				fmt.Printf("Usuario: %s | Impresora: %s | Documento: %s | Páginas: %d\n",
					job.Owner, job.PrinterName, job.DocumentName, job.TotalPages)
			}
		}
		time.Sleep(5 * time.Second) // Revisar cada 5 segundos
	}
}

// Función que obtiene los trabajos de impresión activos
func getPrintJobs() ([]PrintJob, error) {
	var jobs []PrintJob
	query := "SELECT JobId, PrinterName, TotalPages, Document, JobStatus, Owner, TimeSubmitted FROM Win32_PrintJob"

	err := wmi.Query(query, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
