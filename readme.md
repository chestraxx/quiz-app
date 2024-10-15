### **System Design Document for Real-Time Quiz Feature**

#### **Overview**
This document outlines the system architecture for a real-time quiz feature in an English learning application. The system supports multiple users joining quiz sessions simultaneously, submitting answers, updating scores in real-time, and displaying a real-time leaderboard.

---
### **1. Architecture Diagram**

The architecture diagram illustrates the components and interactions between the client, server, and database, as well as the flow of data between these components.

- **Frontend Service (Web, Mobile)**: Users interact with the quiz system via a Vue.js SPA (Single Page Application) or mobile app.
- **Backend Server**:
   - Rest API: Manages user registration, quiz data, and score submissions.
   - Websocket: Handles real-time communication for quiz participation and leaderboard updates.
- **Database (e.g., MySQL)**: Stores quiz questions, users, scores, and session data.
- **Cache (e.g., Redis)**: Caches real-time leaderboard data and quiz sessions for quick access.

#### **Architecture Diagram**
![[assets/architecture-diagram.png]]

---

### **2. Component Description**

1. **Frontend Service (Vue.js/React)**
   - **Role**: The front-end application where users interact with the quiz system. Users can join quiz sessions, submit answers, and view the real-time leaderboard.
   - **Responsibilities**: 
     - Join a quiz session by a unizue quiz ID.
     - Submit answers to quiz questions.
     - Receive real-time updates for the leaderboard.
   - **Technologies**: Vue.js or React for building a responsive and interactive UI, Axios for making HTTP requests, and WebSocket for real-time updates.

2. **Backend Service (Golang)**
   2.1 **WebSocket**
   - **Role**: The WebSocket server manages real-time communication between clients and the backend. It is responsible for broadcasting score updates and leaderboard changes.
   - **Responsibilities**:
     - Handle connections from multiple clients simultaneously.
     - Broadcast updates (e.g., new answers, leaderboard changes) to all connected clients.
   - **Technologies**: Golang with Gorilla WebSocket library for handling concurrent real-time communications.

   2.2. **REST API**
   - **Role**: The REST API server manages user data, quiz sessions, and score submissions. It interacts with the database and ensures data consistency.
   - **Responsibilities**:
     - Provide API endpoints for user registration, joining quizzes, submitting answers, and retrieving quiz questions.
     - Validate answers and calculate scores.
   - **Technologies**: Golang with Gorilla Mux for handling HTTP routes, and JSON for communication between the client and server.

4. **Database (MySQL)**
   - **Role**: The primary data store for quiz-related data.
   - **Responsibilities**:
     - Store quiz sessions, questions, user profiles, and scores.
     - Ensure the integrity of quiz session data and score calculations.
   - **Technologies**: MySQL, chosen for easy-to-use and lightweight.

5. **Cache (Redis)**
   - **Role**: A high-performance, in-memory data store to speed up access to frequently accessed data.
   - **Responsibilities**:
     - Cache real-time leaderboard data for quick updates.
     - Cache quiz session data to reduce database load during heavy traffic.
   - **Technologies**: Redis for its speed and efficiency in handling real-time data and caching frequently accessed information.

---

### **3. Data Flow**

1. **User Joins a Quiz**
   - The user inputs a unique quiz ID into Frontend Service (Vue.js).
   - The client sends a request to the REST API to join the quiz session.
   - Backend Service verifies the quiz ID, retrieves the quiz details from MySQL, and returns the data to the client.
   - Backend Service add User information to the quiz session.
   - The client establishes a WebSocket connection for real-time updates.

2. **User Submits an Answer**
   - The user submits their answer via Frontend Service (Vue.js), which is sent to Backend Service Rest API through an HTTP POST request.
   - Backend Service Rest API checks the answer's correctness and updates the user's score in PostgreSQL.
   - The updated score is cached in Redis for faster leaderboard access.

3. **Leaderboard Update**
   - Upon score updates, Backend Service WebSocket broadcasts the new scores and updated leaderboard to all connected clients.
   - Frontend Service update their leaderboards in real-time, displaying the latest scores of all participants.

---

### **4. Technologies and Tools**

| Component             | Technology                 | Justification |
|-----------------------|----------------------------|---------------|
| **Frontend Service**  | Vue.js             | Both frameworks offer a reactive UI, which is essential for real-time updates and displaying leaderboards efficiently. |
| **Backend Service Real-Time Communication** | WebSocket (Gorilla WebSocket for Golang) | WebSockets enable bidirectional communication, perfect for real-time score updates and leaderboard notifications. |
| **Backend Service API Backend**        | Golang (Gorilla Mux)       | Golang provides excellent concurrency support, making it ideal for managing multiple quiz sessions and real-time interactions. |
| **Database**           | MySQL                 | MySQL is a reliable relational database which is easy-to-use & lightweight. |
| **Cache**              | Redis                      | Redis ensures fast, in-memory access to real-time data like the leaderboard and cached quiz sessions, reducing the load on the database. |
