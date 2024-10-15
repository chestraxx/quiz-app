// Package quiz provides functionality for managing quiz sessions and participants.
package quiz

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

var quizSessions = make(map[string]*QuizSession)
var mutex = &sync.RWMutex{}

// QuizSession represents a quiz session.
type QuizSession struct {
	ID           string
	Questions    []Question
	Participants map[string]Participant
	mutex        sync.RWMutex
}

// Question represents a quiz question.
type Question struct {
	ID      string
	Text    string
	Options []string
	Correct string
}

// Participant represents a quiz participant.
type Participant struct {
	ID      string
	Score   int
	Answers map[string]string
}

// GetQuizSession returns a quiz session by ID.
func GetQuizSession(id string) (*QuizSession, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	quizSession, ok := quizSessions[id]
	if !ok {
		return nil, errors.New("quiz session not found")
	}
	return quizSession, nil
}

// AddQuizSession adds a quiz session to the map.
func AddQuizSession(quizSession *QuizSession) {
	mutex.Lock()
	defer mutex.Unlock()
	quizSessions[quizSession.ID] = quizSession
}

// NewQuizSession returns a new quiz session.
func NewQuizSession(id string, questions []Question) *QuizSession {
	quizSession := &QuizSession{
		ID:           id,
		Questions:    questions,
		Participants: make(map[string]Participant),
	}
	AddQuizSession(quizSession)
	return quizSession
}

// GetQuestionID returns the ID of the question that matches the given answer.
func (qs *QuizSession) GetQuestionID(answer string) string {
	for _, q := range qs.Questions {
		for _, option := range q.Options {
			if option == answer {
				return q.ID
			}
		}
	}
	return ""
}

// GetQuestion returns a question by ID.
func (qs *QuizSession) GetQuestion(id string) (*Question, error) {
	for _, q := range qs.Questions {
		if q.ID == id {
			return &q, nil
		}
	}

	return nil, errors.New("question not found")
}

// AddParticipant adds a participant to the quiz session.
func (qs *QuizSession) AddParticipant(id string) error {
	qs.mutex.Lock()
	defer qs.mutex.Unlock()
	if _, ok := qs.Participants[id]; ok {
		return errors.New("participant already exists")
	}
	qs.Participants[id] = Participant{
		ID:      id,
		Score:   0,
		Answers: make(map[string]string),
	}
	return nil
}

// GetParticipantByID returns a participant by ID.
func (qs *QuizSession) GetParticipantByID(id string) (*Participant, error) {
	qs.mutex.RLock()
	defer qs.mutex.RUnlock()
	participant, ok := qs.Participants[id]
	if !ok {
		return nil, fmt.Errorf("participant %q not found", id)
	}
	return &participant, nil
}

// SubmitAnswer submits an answer for a participant.
func (qs *QuizSession) SubmitAnswer(participantID, questionID, answer string) error {
	qs.mutex.Lock()
	defer qs.mutex.Unlock()

	participant, ok := qs.Participants[participantID]
	if !ok {
		return errors.New("participant not found")
	}

	question, err := qs.GetQuestion(questionID)
	if err != nil {
		return err
	}

	if question.Correct == answer {
		participant.Score++
	}

	participant.Answers[questionID] = answer
	qs.Participants[participantID] = participant

	return nil
}

// GetParticipants returns the list of participants in the quiz session, sorted by score in descending order.
func (qs *QuizSession) GetParticipants() []*Participant {
	qs.mutex.RLock()
	defer qs.mutex.RUnlock()

	participants := make([]*Participant, 0, len(qs.Participants))
	for _, p := range qs.Participants {
		participants = append(participants, &p)
	}

	sort.Slice(participants, func(i, j int) bool {
		return participants[i].Score > participants[j].Score
	})
	return participants
}
