import { defineStore } from 'pinia';

// Define the store
export const useQuizStore = defineStore('quiz', {
  // State properties
  state: () => ({
    participant: {
      participant_id: '',
      quiz_id: '',
    }
  }),
  
  // Actions to modify state
  actions: {
    // Set participants (e.g., loading from an API)
    setParticipant(participant) {
      this.participant.participant_id = participant.participant_id
      this.participant.quiz_id = participant.quiz_id
    }
  },

  // Optional: Getters for computed properties
  getters: {
    participantInfo: (state) => state.participant,
    participantExists: (state) => state.participant.participant_id !== '',
  }
});