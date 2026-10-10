# Quiz API Service
Quiz API service to manage quiz system. Provides basic quiz logic handling, user authentication and stats service.

## Architecture
Project consists of 3 main microservices:

* **Quiz Service** — Quiz management, answer check etc
* **Auth Service** — Registration, Authentification using **gRPC**
* **Stats Service** — User stats handling after completing quizzes using **Kafka**

## Quick start

### Prerequisites
* Go 1.2x or more
* Docker

### Clone
```bash
git clone https://github.com/xashaysso/quiz.git
cd quiz
```

*(change variables and ports if needed in `docker-compose.yaml` file).*

### Launch containers
```bash
docker compose up -d --build
```

### Ports
After using `docker compose up -d` services will be available at:
* **Quiz Service:** `http://localhost:8080`
* **Stats Service:** `http://localhost:8081`
* **Auth Service (gRPC):** `localhost:50051`

### DBs:
* **Quiz DB:** `localhost:5433`
* **Stats DB:** `localhost:5434`
* **Auth DB:** `localhost:5435`
* **Redis:** `localhost:6379`
* **Kafka:** `localhost:9092`

## Endpoints

All routes under `/quizzes`, `/questions`, and `/answers` are protected by `authMiddleware` and require a valid Bearer token.

### Authentication (`/auth`)

(`POST`) `/auth/register` - Register a new user  
(`POST`) `/auth/login` - Authenticate and retrieve token  
(`POST`) `/auth/logout` - Log out user  

### Quizzes (`/quizzes`)

(`GET`) `/quizzes/` - Fetch all quizzes  
(`POST`) `/quizzes/` - Create a new quiz  
(`PATCH`) `/quizzes/:quiz_id` - Update quiz details  
(`DELETE`) `/quizzes/:quiz_id` - Delete a quiz  
(`POST`) `/quizzes/:quiz_id/start` - Start a quiz attempt  
(`GET`) `/quizzes/:quiz_id/questions` - Get all questions for a quiz  
(`POST`) `/quizzes/:quiz_id/questions` - Add a question to a quiz  

### Questions (`/questions`)

(`GET`) `/questions/:question_id` - Get question by ID  
(`PATCH`) `/questions/:question_id` - Update question details  
(`DELETE`) `/questions/:question_id` - Delete a question  
(`GET`) `/questions/:question_id/answers` - Get answer options for a question  
(`POST`) `/questions/:question_id/answers` - Add an answer option  
(`POST`) `/questions/:question_id/check` - Submit and check selected answer  

### Answers (`/answers`)

(`GET`) `/answers/:answer_id` - Get answer option by ID  
(`PATCH`) `/answers/:answer_id` - Update answer option  
(`DELETE`) `/answers/:answer_id` - Delete answer option

## Useful Commands

* **View logs across all services:**
  ```bash
  docker compose logs -f
  ```
* **View logs for a specific service (e.g., main_app):**
  ```bash
  docker compose logs -f main_app
  ```
* **Stop containers:**
  ```bash
  docker compose down
  ```
* **Stop containers and reset all volumes:**
  ```bash
  docker compose down -v
  ```
