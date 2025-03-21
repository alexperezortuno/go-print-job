package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	username string
	password string
	urlLogin string
	urlJobs  string
	db       *gorm.DB
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

type Jobs struct {
	JobId          uint32 `gorm:"column:job_id;primaryKey"`
	Document       string `gorm:"column:document"`
	PaperSize      string `gorm:"column:paper_size"`
	PaperWidth     uint32 `gorm:"column:paper_width"`
	PrintProcessor string `gorm:"column:print_processor"`
	Size           uint32 `gorm:"column:size"`
	TotalPages     uint32 `gorm:"column:total_pages"`
	Status         bool   `gorm:"column:status"`
}

func (p PrintJob) String() string {
	return fmt.Sprintf("JobId: %d | Usuario: %s | Impresora: %s | Documento: %s | Páginas: %d | Size: %d",
		p.JobId, p.Owner, p.PrintProcessor, p.Document, p.TotalPages, p.Size)
}

func initDatabase(migrate bool) {
	var err error
	db, err = gorm.Open(sqlite.Open("jobs.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	if migrate {
		err = db.AutoMigrate(&Jobs{})
		if err != nil {
			log.Fatal("failed to migrate the database:", err)
			return
		}
	}
}

func saveJob(job PrintJob) {
	if err := db.Create(&Jobs{
		JobId:          job.JobId,
		Document:       job.Document,
		PaperSize:      job.PaperSize,
		PaperWidth:     job.PaperWidth,
		PrintProcessor: job.PrintProcessor,
		Size:           job.Size,
		TotalPages:     job.TotalPages,
		Status:         false,
	}).Error; err != nil {
		log.Println("error saving job", err)
	}
	log.Println("job saved successfully")
}

func updateJob(job PrintJob) {
	if err := db.Model(&Jobs{}).Where("job_id = ?", job.JobId).Updates(Jobs{
		Status: true,
	}).Error; err != nil {
		log.Println("error updating job", err)
	}
	log.Println("job updated successfully")
}

func checkNetwork() bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}
	_, err := client.Get("http://www.google.com")
	if err != nil {
		log.Println("No internet connection:", err)
		return false
	}
	log.Println("Internet connection is available")
	return true
}

func getCredentials() (string, string) {
	v := map[string]string{
		"username": username,
		"password": password,
	}
	data := new(bytes.Buffer)
	if err := json.NewEncoder(data).Encode(v); err != nil {
		log.Println("error encoding data:", err)
		return "", "error encoding data"
	}

	req, err := http.NewRequest("POST", urlLogin, data)

	if err != nil {
		log.Println("error creating request:", err)
		return "", "error creating request"
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Println("error sending request:", err)
		return "", "error sending request"
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("invalid status code:", resp.StatusCode)
		return "", "error invalid status code"
	}

	var token struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		log.Println("error decoding response:", err)
		return "", "error decoding response"
	}

	return token.Token, ""
}

func sendJob(job PrintJob, token string) bool {
	data, err := json.Marshal(job)
	if err != nil {
		log.Println("error marshalling job:", err)
		return false
	}

	req, err := http.NewRequest("POST", urlJobs, bytes.NewBuffer(data))
	if err != nil {
		log.Println("error creating request:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error sending request:", err)
		return false
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error closing response body:", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Println("invalid status code:", resp.StatusCode)
		return false
	}

	log.Println("job sent successfully")
	return true
}

func scheduleTask() {
	var job Jobs
	db.Where("status = ?", false).First(&job)
	if job.JobId > 0 {
		log.Println("job pending for send to server:", job.JobId)
		// Send job to server
		if checkNetwork() {
			token, err := getCredentials()
			if err != "" {
				log.Println("error getting credentials:", err)
				return
			}

			j := PrintJob{
				JobId:          job.JobId,
				Document:       job.Document,
				PaperSize:      job.PaperSize,
				PaperWidth:     job.PaperWidth,
				PrintProcessor: job.PrintProcessor,
				Size:           job.Size,
				TotalPages:     job.TotalPages,
			}

			// Send job to server
			sended := sendJob(j, token)

			if sended {
				updateJob(j)
			}
		}
	}
}

func main() {
	initDatabase(true)
	log.Println("Iniciando monitoreo de impresiones...")

	for {
		printJobs, err := getPrintJobs()
		if err != nil {
			log.Println("Error obteniendo trabajos de impresión:", err)
		} else {
			for _, job := range printJobs {
				if job.TotalPages > 0 {
					log.Printf(job.String())
					saveJob(job)
				}
			}
		}
		time.Sleep(1 * time.Second) // Revisar cada 5 segundos
	}
}

// Función que obtiene los trabajos de impresión activos
func getPrintJobs() ([]PrintJob, error) {
	var jobs []PrintJob
	query := "SELECT * FROM Win32_PrintJob"

	err := wmi.Query(query, &jobs)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}
