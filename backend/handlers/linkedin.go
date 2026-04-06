package handlers

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strings"

	"cv-go/models"

	"github.com/gin-gonic/gin"
)

// ImportLinkedIn parses a LinkedIn data export ZIP file and returns CVData
func ImportLinkedIn(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large (max 10MB)"})
		return
	}

	// Read file into memory for zip processing
	buf, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read file"})
		return
	}

	reader, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ZIP file"})
		return
	}

	cvData := models.CVData{}

	for _, f := range reader.File {
		name := strings.ToLower(f.Name)
		rc, err := f.Open()
		if err != nil {
			continue
		}

		switch {
		case strings.Contains(name, "profile") && strings.HasSuffix(name, ".csv"):
			parseProfile(rc, &cvData)
		case strings.Contains(name, "position") && strings.HasSuffix(name, ".csv"):
			parsePositions(rc, &cvData)
		case strings.Contains(name, "education") && strings.HasSuffix(name, ".csv"):
			parseEducation(rc, &cvData)
		case strings.Contains(name, "skill") && strings.HasSuffix(name, ".csv"):
			parseSkills(rc, &cvData)
		case strings.Contains(name, "certification") && strings.HasSuffix(name, ".csv"):
			parseCertifications(rc, &cvData)
		}
		rc.Close()
	}

	c.JSON(http.StatusOK, gin.H{"data": cvData})
}

func parseProfile(r io.Reader, data *models.CVData) {
	records := readCSV(r)
	if len(records) < 2 {
		return
	}
	headers := records[0]
	values := records[1]

	fieldMap := mapFields(headers, values)
	data.PersonalInfo.FullName = fmt.Sprintf("%s %s", fieldMap["First Name"], fieldMap["Last Name"])
	data.PersonalInfo.Email = fieldMap["Email Address"]
	data.PersonalInfo.Address = fieldMap["Geo Location"]
	data.PersonalInfo.JobTitle = fieldMap["Headline"]
	data.Objective = fieldMap["Summary"]
}

func parsePositions(r io.Reader, data *models.CVData) {
	records := readCSV(r)
	if len(records) < 2 {
		return
	}
	headers := records[0]
	for _, row := range records[1:] {
		fieldMap := mapFields(headers, row)
		exp := models.Experience{
			Title:       fieldMap["Title"],
			Company:     fieldMap["Company Name"],
			Location:    fieldMap["Location"],
			StartDate:   fieldMap["Started On"],
			EndDate:     fieldMap["Finished On"],
			Description: fieldMap["Description"],
		}
		if exp.EndDate == "" {
			exp.Current = true
		}
		if exp.Title != "" || exp.Company != "" {
			data.Experiences = append(data.Experiences, exp)
		}
	}
}

func parseEducation(r io.Reader, data *models.CVData) {
	records := readCSV(r)
	if len(records) < 2 {
		return
	}
	headers := records[0]
	for _, row := range records[1:] {
		fieldMap := mapFields(headers, row)
		edu := models.Education{
			Degree:      fieldMap["Degree Name"],
			Institution: fieldMap["School Name"],
			StartDate:   fieldMap["Start Date"],
			EndDate:     fieldMap["End Date"],
			Description: fieldMap["Notes"],
		}
		if edu.Institution != "" {
			data.Education = append(data.Education, edu)
		}
	}
}

func parseSkills(r io.Reader, data *models.CVData) {
	records := readCSV(r)
	if len(records) < 2 {
		return
	}
	headers := records[0]
	for _, row := range records[1:] {
		fieldMap := mapFields(headers, row)
		name := fieldMap["Name"]
		if name == "" {
			// Try alternative column name
			for _, h := range headers {
				if strings.Contains(strings.ToLower(h), "skill") {
					name = fieldMap[h]
					break
				}
			}
		}
		if name != "" {
			data.Skills = append(data.Skills, models.Skill{
				Name:  name,
				Level: "intermediate",
			})
		}
	}
}

func parseCertifications(r io.Reader, data *models.CVData) {
	records := readCSV(r)
	if len(records) < 2 {
		return
	}
	headers := records[0]
	for _, row := range records[1:] {
		fieldMap := mapFields(headers, row)
		cert := models.Certificate{
			Name:   fieldMap["Name"],
			Issuer: fieldMap["Authority"],
			Date:   fieldMap["Started On"],
			URL:    fieldMap["Url"],
		}
		if cert.Name != "" {
			data.Certificates = append(data.Certificates, cert)
		}
	}
}

func readCSV(r io.Reader) [][]string {
	reader := csv.NewReader(r)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		return nil
	}
	return records
}

func mapFields(headers, values []string) map[string]string {
	m := make(map[string]string)
	for i, h := range headers {
		if i < len(values) {
			m[strings.TrimSpace(h)] = strings.TrimSpace(values[i])
		}
	}
	return m
}
