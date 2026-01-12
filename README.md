# CVWO Winter Assignment NUS Lifters Web Forum 

## Project Overview
NUS Lifters Web Forum is a full-stack web forum application that allows users to register accounts, create fitness and gym-related topics, posts and comments as well as interact through likes. 

---

## Features
- User registration and login
- JWT-based authentication
- Topic creation, editing, and deletion
- Post creation, editing, deletion, and liking
- Comment creation, editing, deletion, and liking
- Search and sorting for topics and posts
- Responsive UI using Material UI

## Tech Stack
**Frontend**
- React (TypeScript)
- Material UI
- Axios

**Backend**
- Go
- SQLite database
- JWT authentication

**Deployment**
- Frontend: Netlify
- Backend: Render

---

## Setup Instructions

### Backend Setup
1. Clone the backend repository:  
   https://github.com/Sengernest/CVWO-NUS-Lifters-Club-Web-Forum-Backend.git
2. Navigate into the cloned directory
3. Before running the server, a JWT secret key must be set as an environment variable in the terminal:
   - **Windows (PowerShell)**  
     ```powershell
     $env:JWT_KEY="cvwo-secret-key"
     ```
   - **Windows (Command Prompt)**  
     ```cmd
     set JWT_KEY=cvwo-secret-key
     ```
   - **macOS / Linux**  
     ```bash
     export JWT_KEY="cvwo-secret-key"
     ```
4. Run the backend server:
   ```bash
   go run main.go
5. Backend will run on
   ```bash
   https://localhost8080

 ### Frontend Setup
 1. Clone the frontend repo: https://github.com/Sengernest/cvwo-nus-lifters-club-web-forum-frontend.git
 2. Go into the directory which the repo was cloned
 3. Install dependencies:
    ```bash
    npm install
 4. Start the frontend using this command:
    ```bash
    npm start
 5. Frontend will run on:
    ```bash
    http://localhost:3000

---

## Deployed Application 
**Frontend (Netlify):**
https://nus-lifters-web-forum.netlify.app/forum

**Backend (Render):**
https://cvwo-nus-lifters-club-web-forum-backend.onrender.com

---

## AI Usage Declaration
AI tool (ChatGPT) was used strictly as a learning aid and code review
assistant, in accordance with assignment guidelines.

ChatGPT was used for:
- Understanding deployment concepts such as environment variables,
CORS configuration, and frontend–backend integration.
- Interpreting error messages encountered during development
(React, Axios, and Material UI).
- Reviewing code structure and suggesting refactoring improvements.
- Learning how to use Material UI components effectively.

---
