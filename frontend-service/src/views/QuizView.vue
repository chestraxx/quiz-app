<template>
  <div>
    Participant ID: {{ quizStore.participantInfo.participant_id }}

    <div class="card">
      <h2>Score Leaderboard</h2>
        <ul id="score-list">
          <li v-for="(participant, pIndex) in leaderboard" :key="participant.id">
            {{ participant.id }}: {{ participant.score }} points
          </li>
        </ul>
    </div>

    <div class="card" v-if="!isSubmitted">
      <h1>Quiz Time!</h1>
      <p>Answer the following questions to test your knowledge.</p>
      <form id="quiz-form">
        <div class="question-container" v-for="(question, qIndex) in questions" :key="question.ID">
          <h4>{{ question.Text }}</h4>

          <div v-for="(option, oIndex) in question.Options">
            <div>
              <input type="radio" :id="`answer-${qIndex}-${oIndex}`" :name="`answer-${qIndex}`" :value="option" v-model="answer[`${question.ID}`]">
              <label :for="`answer-${qIndex}-${oIndex}`">{{ option }}</label>
            </div>
          </div>

          <hr>
        </div>

        <button type="submit" @click.prevent.stop="submitAnswer()">Submit Answer</button>
      </form>
    </div>
  </div>
</template>

<script>
import axios from 'axios';
import { useQuizStore } from '../stores/quiz';
const TYPE_LEADERBOARD = "scoreLeaderboard"

export default {
  name: 'QuizView',

  data() {
    return {
      quizStore: useQuizStore(),
      questions: [],
      answer: {},
      isSubmitted: false,
      connection: null,
      leaderboard: [],
    }
  },

  mounted() {
    this.initPage()
  },

  async beforeRouteLeave (to, from, next) {
    let answer = window.confirm("You will lose your progress if you leave. Are you sure?");
    if (answer) {
      try{
          // Resolved
          next() 
      } catch (err) { // Rejected
          next(false)
          log(err)
      }
    } else {
      this.initPage()
      next(false)
    }
  },

  methods: {
    initPage() {
      if (!this.quizStore.participantExists) {
        this.$router.push({ name: 'join'})
        return
      }

      this.fetchQuestions()
      this.connectWebsocket()
    },

    buildAnswer() {
      for (let index = 0; index < this.questions.length; index++) {
        this.answer[`${this.questions[index].ID}`] = ''
      }
    },

    connectWebsocket() {
      this.connection = new WebSocket('ws://localhost:9090/ws');
      this.connection.onopen = () => {
        console.log('Connected to WebSocket server');
      };
      this.connection.onmessage = (event) => {
        const {type, data} = JSON.parse(event.data);

        if (type == TYPE_LEADERBOARD) {
          this.leaderboard = data
        }
      }
    },

    fetchQuestions() {
      axios.get('http://localhost:9090/quizzes/questions', {params: {
        id: this.quizStore.participantInfo.participant_id,
        quiz_id: this.quizStore.participantInfo.quiz_id,
      }}).then(({data}) => {
        this.questions = data
        this.buildAnswer()
      }).catch(error => {
        console.error('Error fetching questions:', error);
      })
    },

    submitAnswer() {
      axios.post('http://localhost:9090/quizzes/answers', {
        id: this.quizStore.participantInfo.participant_id,
        quiz_id: this.quizStore.participantInfo.quiz_id,
        answer: this.answer,
      })
      .then(({data}) => {
        alert('Submit answer successful!');

        this.isSubmitted = true
      })
      .catch(error => {
        console.error('Error submitting answer:', error);
      });
    },
  },
};
</script>