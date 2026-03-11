package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cv-go/models"
	"cv-go/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AIHandler struct {
	DB     *gorm.DB
	AESKey string
}

func (h *AIHandler) callAI(prompt string) (string, error) {
	var setting models.AISetting
	if err := h.DB.Where("is_active = ?", true).First(&setting).Error; err != nil {
		return "", fmt.Errorf("no active AI provider configured")
	}

	apiKey, err := utils.DecryptAES(setting.APIKeyEnc, h.AESKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt API key")
	}

	maxTokens := setting.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2000
	}

	switch setting.Provider {
	case "openai":
		return callOpenAI(apiKey, setting.Model, prompt, maxTokens)
	case "gemini":
		return callGemini(apiKey, setting.Model, prompt, maxTokens)
	case "deepseek":
		return callDeepSeek(apiKey, setting.Model, prompt, maxTokens)
	default:
		return "", fmt.Errorf("unsupported provider: %s", setting.Provider)
	}
}

func callOpenAI(apiKey, model, prompt string, maxTokens int) (string, error) {
	if model == "" {
		model = "gpt-4o-mini"
	}
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": maxTokens,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OpenAI API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}
	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format from OpenAI")
	}
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format from OpenAI")
	}
	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content format from OpenAI")
	}
	return content, nil
}

func callGemini(apiKey, model, prompt string, maxTokens int) (string, error) {
	if model == "" {
		model = "gemini-1.5-flash"
	}
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"maxOutputTokens": maxTokens,
		},
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	candidates, ok := result["candidates"].([]interface{})
	if !ok || len(candidates) == 0 {
		return "", fmt.Errorf("no response from Gemini")
	}
	candidate, ok := candidates[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format from Gemini")
	}
	content, ok := candidate["content"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid content format from Gemini")
	}
	parts, ok := content["parts"].([]interface{})
	if !ok || len(parts) == 0 {
		return "", fmt.Errorf("invalid parts format from Gemini")
	}
	part, ok := parts[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid part format from Gemini")
	}
	text, ok := part["text"].(string)
	if !ok {
		return "", fmt.Errorf("invalid text format from Gemini")
	}
	return text, nil
}

func callDeepSeek(apiKey, model, prompt string, maxTokens int) (string, error) {
	if model == "" {
		model = "deepseek-chat"
	}
	body := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens": maxTokens,
	}
	data, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", "https://api.deepseek.com/v1/chat/completions", bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("no response from DeepSeek")
	}
	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid response format from DeepSeek")
	}
	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid message format from DeepSeek")
	}
	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("invalid content format from DeepSeek")
	}
	return content, nil
}

func (h *AIHandler) checkCredits(userID uint, required int) (*models.User, error) {
	var user models.User
	if err := h.DB.First(&user, userID).Error; err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if user.AICredits < required {
		return nil, fmt.Errorf("insufficient credits. Required: %d, Available: %d", required, user.AICredits)
	}
	return &user, nil
}

func (h *AIHandler) deductCredits(user *models.User, credits int, feature string) {
	h.DB.Model(user).Update("ai_credits", gorm.Expr("ai_credits - ?", credits))
	h.DB.Create(&models.AIUsageLog{
		UserID:  user.ID,
		Feature: feature,
		Credits: credits,
	})
}

// verifyCVOwnership checks that the CV belongs to the requesting user
func (h *AIHandler) verifyCVOwnership(userID uint, cvID uint64) (*models.CV, error) {
	var cv models.CV
	if err := h.DB.First(&cv, cvID).Error; err != nil {
		return nil, fmt.Errorf("CV not found")
	}
	if cv.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}
	return &cv, nil
}

func (h *AIHandler) ImproveText(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		Text     string `json:"text" binding:"required"`
		Language string `json:"language"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.Text) > 10000 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Text too long (max 10000 characters)"})
		return
	}

	user, err := h.checkCredits(userID, 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lang := "English"
	if input.Language == "ar" {
		lang = "Arabic"
	}

	prompt := fmt.Sprintf("Improve the following CV text for a professional resume. Keep it in %s. Make it more impactful and professional. Only return the improved text, nothing else:\n\n%s", lang, input.Text)

	result, err := h.callAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.deductCredits(user, 1, "improve_text")
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) AnalyzeCV(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	user, err := h.checkCredits(userID, 2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cv, err := h.verifyCVOwnership(userID, cvID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	cvJSON, _ := json.Marshal(cv.Data)
	prompt := fmt.Sprintf(`Analyze this CV data and provide an ATS compatibility assessment. Include:
1. Overall ATS Score (0-100)
2. Strengths (bullet points)
3. Weaknesses (bullet points)
4. Specific improvement suggestions
5. Missing sections or information

CV Data:
%s`, string(cvJSON))

	result, err := h.callAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.deductCredits(user, 2, "analyze_cv")
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) GenerateCoverLetter(c *gin.Context) {
	userID := c.GetUint("user_id")
	var input struct {
		CVID     uint   `json:"cv_id" binding:"required"`
		JobTitle string `json:"job_title" binding:"required"`
		Company  string `json:"company" binding:"required"`
		Language string `json:"language"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.checkCredits(userID, 2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cv, err := h.verifyCVOwnership(userID, uint64(input.CVID))
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	lang := "English"
	if input.Language == "ar" {
		lang = "Arabic"
	}

	cvJSON, _ := json.Marshal(cv.Data)
	prompt := fmt.Sprintf(`Generate a professional cover letter in %s for the following:
Job Title: %s
Company: %s

Based on this CV:
%s

The cover letter should be formal, professional, and tailored to the position.`, lang, input.JobTitle, input.Company, string(cvJSON))

	result, err := h.callAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.deductCredits(user, 2, "cover_letter")
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) SuggestJobs(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	user, err := h.checkCredits(userID, 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cv, err := h.verifyCVOwnership(userID, cvID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	var skills []string
	for _, s := range cv.Data.Skills {
		skills = append(skills, s.Name)
	}

	cvJSON, _ := json.Marshal(cv.Data)
	prompt := fmt.Sprintf(`Based on this CV, suggest 10 suitable job positions. For each job:
1. Job Title
2. Why it's a good fit
3. Key skills that match
4. Estimated salary range

Skills: %s

CV Data:
%s`, strings.Join(skills, ", "), string(cvJSON))

	result, err := h.callAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.deductCredits(user, 1, "suggest_jobs")
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) EvaluateResearch(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	user, err := h.checkCredits(userID, 2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cv, err := h.verifyCVOwnership(userID, cvID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	cvJSON, _ := json.Marshal(cv.Data)
	prompt := fmt.Sprintf(`Evaluate the research capabilities shown in this CV. Assess:
1. Research Experience Score (0-100)
2. Publication potential
3. Research methodology skills
4. Academic strengths
5. Areas for improvement
6. Recommended research directions

CV Data:
%s`, string(cvJSON))

	result, err := h.callAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.deductCredits(user, 2, "evaluate_research")
	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *AIHandler) GetAITips(c *gin.Context) {
	userID := c.GetUint("user_id")
	cvID, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	user, err := h.checkCredits(userID, 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cv, err := h.verifyCVOwnership(userID, cvID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	cvJSON, _ := json.Marshal(cv.Data)
	prompt := fmt.Sprintf(`Provide personalized tips to improve this CV. Include:
1. Top 5 most impactful improvements
2. Formatting suggestions
3. Content improvements for each section
4. Keywords to add for better ATS compatibility
5. Common mistakes to fix

CV Data:
%s`, string(cvJSON))

	result, err := h.callAI(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.deductCredits(user, 1, "get_tips")
	c.JSON(http.StatusOK, gin.H{"result": result})
}
