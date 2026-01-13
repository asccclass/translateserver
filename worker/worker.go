package worker

import (
	"log"
	"time"
	"translateserver/config"
	"translateserver/db"
	"translateserver/notify"
	"translateserver/runner"
)

func Start() {
	ticker := time.NewTicker(time.Duration(config.CheckInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			processJobs()
		}
	}
}

func processJobs() {
	rows, err := db.DB.Query("SELECT id, filename, filepath, email FROM jobs WHERE status = 'pending'")
	if err != nil {
		log.Println("Error querying jobs:", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var filename, filepath, email string
		if err := rows.Scan(&id, &filename, &filepath, &email); err != nil {
			continue
		}

		// Lock job
		_, err := db.DB.Exec("UPDATE jobs SET status = 'processing', updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
		if err != nil {
			log.Println("Error updating job status:", err)
			continue
		}

		log.Printf("Processing job %d: %s", id, filename)

		// Run Docker Command
		output, err := runner.RunWhisper(filename)
		status := "completed"
		if err != nil {
			log.Printf("Job %d failed: %v", id, err)
			status = "failed"
			output = err.Error()
		}

		// Update DB
		_, err = db.DB.Exec("UPDATE jobs SET status = ?, output_text = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", status, output, id)
		if err != nil {
			log.Println("Error updating job result:", err)
		}

		// Notify
		notify.SendEmail(email, status, output)
	}
}
