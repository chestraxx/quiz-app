<template>
  <div class="card">
    <div>
      <label for="name">Quiz ID:</label>
      <input type="text" id="quiz-id" name="quiz-id" v-model="quizId" required />
    </div>
    <button type="submit" @click.prevent.stop="joinQuizSession()">Join</button>
  </div>
</template>

<script>
import axios from 'axios';
import { useQuizStore } from '../stores/quiz';

export default {
  name: 'JoinView',

  data() {
    return {
      quizId: 'quiz-01',
      quizStore: useQuizStore(),
    }
  },

  methods: {
    joinQuizSession() {
      axios.post('http://localhost:9090/quizzes/participant', {
        id: this.generateUniqueId(),
        quiz_id: this.quizId
      })
      .then(({data}) => {
        if (data && data.participant_id) {
          this.quizStore.setParticipant(data);

          alert('Join successful!');

          this.$router.push({ name: 'quiz'})
        }
      })
      .catch(error => {
        alert('Join failed! Please retry with a different ID.');

        console.error('Error registering participant:', error);
      });
    },

    generateUniqueId() {
      return Math.random().toString(36).substr(2, 9);
    },
  },
};
</script>

<style>
body {
  font-family: Arial, sans-serif;
}

label {
  display: block;
  margin-bottom: 10px;
}

input[type='text'] {
  width: -webkit-fill-available;
  height: 40px;
  margin-bottom: 20px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 5px;
}

button[type='submit'] {
  width: 100%;
  height: 40px;
  background-color: #4caf50;
  color: #fff;
  padding: 10px;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

button[type='submit']:hover {
  background-color: #3e8e41;
}

.card {
  width: 300px;
  margin: 40px auto;
  padding: 20px;
  border: 1px solid #ccc;
  border-radius: 10px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}

.d-flex {
  display: flex;
}
</style>
