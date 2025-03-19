package main

import (
	"log"
	"time"

	"github.com/yusufpapurcu/wmi"
)

// Estructura para almacenar información sobre trabajos de impresión
type PrintJob struct {
	//Caption        string    `wmi:"Caption"`
	//Description    string    `wmi:"Description"`
	//InstallDate    time.Time `wmi:"InstallDate"`
	//Name           string    `wmi:"Name"`
	//Status         string    `wmi:"Status"`
	//ElapsedTime    time.Time `wmi:"ElapsedTime"`
	//JobStatus      string    `wmi:"JobStatus"`
	//Notify         string    `wmi:"Notify"`
	Owner string `wmi:"Owner"`
	//Priority       uint32    `wmi:"Priority"`
	//StartTime      time.Time `wmi:"StartTime"`
	//TimeSubmitted  time.Time `wmi:"TimeSubmitted"`
	//UntilTime      time.Time `wmi:"UntilTime"`
	//Color          string    `wmi:"Color"`
	//DataType       string    `wmi:"DataType"`
	Document string `wmi:"Document"`
	//DriverName     string    `wmi:"DriverName"`
	//HostPrintQueue string    `wmi:"HostPrintQueue"`
	JobId uint32 `wmi:"JobId"`
	//PagesPrinted   uint32    `wmi:"PagesPrinted"`
	//PaperLength    uint32    `wmi:"PaperLength"`
	PaperSize  string `wmi:"PaperSize"`
	PaperWidth uint32 `wmi:"PaperWidth"`
	//Parameters     string    `wmi:"Parameters"`
	PrintProcessor string `wmi:"PrintProcessor"`
	Size           uint32 `wmi:"Size"`
	//StatusMask     uint32    `wmi:"StatusMask"`
	TotalPages uint32 `wmi:"TotalPages"`
}

func main() {
	log.Println("Iniciando monitoreo de impresiones...")

	for {
		printJobs, err := getPrintJobs()
		if err != nil {
			log.Println("Error obteniendo trabajos de impresión:", err)
		} else {
			for _, job := range printJobs {
				if job.TotalPages > 0 {
					log.Printf("JobId: %d | Usuario: %s | Impresora: %s | Documento: %s | Páginas: %d\n",
						job.JobId, job.Owner, job.PrintProcessor, job.Document, job.TotalPages)
				}
			}
		}
		time.Sleep(1 * time.Second) // Revisar cada 5 segundos
	}
}

// Función que obtiene los trabajos de impresión activos
func getPrintJobs() ([]PrintJob, error) {
	var jobs []PrintJob
	query := "SELECT JobId, Owner, PrintProcessor, Document, TotalPages, PaperSize, PaperWidth, PrintProcessor FROM Win32_PrintJob"

	err := wmi.Query(query, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
