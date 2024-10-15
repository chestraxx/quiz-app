package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"quiz-app/internal/quiz"
	"quiz-app/internal/websocket"

	"github.com/huandu/go-clone"
	"github.com/rs/cors"
)

type RegisterParticipantRequest struct {
	ID      string
	QUIZ_ID string
}

type ValidateAnswerRequest struct {
	ID      string
	QUIZ_ID string
	Answer  map[string]string
}

func RegisterParticipantForQuiz(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to register participant")

	var registerParticipantRequest RegisterParticipantRequest
	err := json.NewDecoder(r.Body).Decode(&registerParticipantRequest)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var participantData quiz.Participant
	participantData.ID = registerParticipantRequest.ID

	if participantData.ID == "" {
		log.Println("Missing participant ID")
		http.Error(w, "Participant ID is required", http.StatusBadRequest)
		return
	}

	quizSession, err := quiz.GetQuizSession(registerParticipantRequest.QUIZ_ID)
	if err != nil {
		log.Println("Error getting quiz session:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = quizSession.AddParticipant(participantData.ID)
	if err != nil {
		log.Println("Error adding participant:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Participant registered successfully:", participantData)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"message":        "Participant registered successfully",
		"participant_id": participantData.ID,
		"quiz_id":        quizSession.ID,
	}

	log.Println("Sending response:", response)
	json.NewEncoder(w).Encode(response)
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	quizID := r.URL.Query().Get("quiz_id")

	log.Println("Received request to get questions with id:", id, "and quiz_id:", quizID)

	quizSession, err := quiz.GetQuizSession(quizID)
	if err != nil {
		log.Println("Error getting quiz session:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = quizSession.GetParticipantByID(id)
	if err != nil {
		log.Println("Error getting participant:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	SendScoreLeaderboard(quizSession)

	// Get the questions for the quiz session
	questions := clone.Clone(quizSession.Questions).([]quiz.Question)

	// Delete the Correct property from the questions
	for i := range questions {
		questions[i].Correct = ""
	}

	// Return the questions as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(questions)
}

func ValidateAnswer(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request to submmiting answer")

	var validateAnswerRequest ValidateAnswerRequest
	err := json.NewDecoder(r.Body).Decode(&validateAnswerRequest)
	if err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	quizSession, err := quiz.GetQuizSession(validateAnswerRequest.QUIZ_ID)
	if err != nil {
		log.Println("Error getting quiz session:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	participant, err := quizSession.GetParticipantByID(validateAnswerRequest.ID)
	if err != nil {
		log.Println("Error getting participant:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for questionId, answer := range validateAnswerRequest.Answer {
		log.Println("Answer for question", participant.ID, questionId, answer)
		quizSession.SubmitAnswer(participant.ID, questionId, answer)
	}

	SendScoreLeaderboard(quizSession)

	// Return the questions as a JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	log.Println("Error getting participant:", quizSession.Participants)

	json.NewEncoder(w).Encode([]string{"Answer submitted successfully!"})
}

// sendScoreLeaderboard sends the score leaderboard to all connected clients.
func SendScoreLeaderboard(quizSession *quiz.QuizSession) {
	// Get the list of participants from the quiz session.
	participants := quizSession.GetParticipants()

	log.Println("SendScoreLeaderboard", participants)

	// Create a response struct to hold the participant information.
	type participantResponse struct {
		ID    string `json:"id"`
		Score int    `json:"score"`
	}

	// Create a slice to hold the participant responses.
	var participantResponses []participantResponse

	// Iterate over the participants and create a response for each one.
	for _, p := range participants {
		participantResponses = append(participantResponses, participantResponse{
			ID:    p.ID,
			Score: p.Score,
		})
	}

	// Send the participant responses to all WebSocket connections.
	response := struct {
		Type string                `json:"type"`
		Data []participantResponse `json:"data"`
	}{
		Type: "scoreLeaderboard",
		Data: participantResponses,
	}

	// Send the response to all connected clients.
	log.Println(websocket.ListConnections)
	for _, conn := range websocket.ListConnections {
		if err := conn.WriteJSON(response); err != nil {
			log.Println(err)
		}
	}
}

func main() {
	// Create a new quiz session
	questions := []quiz.Question{
		{
			ID:      "question-1",
			Text:    "What is the capital of France?",
			Options: []string{"Paris", "London", "Berlin"},
			Correct: "Paris",
		},
		{
			ID:      "question-2",
			Text:    "What is the largest planet in our solar system?",
			Options: []string{"Earth", "Saturn", "Jupiter"},
			Correct: "Jupiter",
		},
	}
	quizSession := quiz.NewQuizSession("quiz-01", questions)

	// Create a new WebSocket handler
	handler := websocket.NewHandler(quizSession)

	// Create a new HTTP server
	http.HandleFunc("POST /quizzes/participant", RegisterParticipantForQuiz)
	http.HandleFunc("GET /quizzes/questions", GetQuestions)
	http.HandleFunc("POST /quizzes/answers", ValidateAnswer)

	// Start the WebSocket server
	http.HandleFunc("/ws", handler.HandleConnection)

	// Enable CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:5173"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "application/json"},
		// AllowCredentials: true, // If using cookies or authentication
	})

	// Wrap your HTTP server with CORS middleware
	handlerWithCORS := c.Handler(http.DefaultServeMux)

	fmt.Println("Server running on port 9090")
	log.Fatal(http.ListenAndServe(":9090", handlerWithCORS))
}
