# Workout Tracker API

A **Workout Tracker** application built with Go, designed to manage workout routines, exercises, and member information. The API provides endpoints to handle CRUD operations for exercises, workouts, and members.

## Features

- **Member Support**: Create and manage member profiles with email and password.
- **Workout Tracking**: Manage workouts and their details.
- **Exercise Management**: Add, update, delete, and fetch exercises.
- **Database Migrations**: Predefined scripts for database setup.
- **JWT Authentication**: Secure the API using JSON Web Tokens (JWT). Members can obtain authentication tokens by sending their credentials to a designated endpoint.
- **User Account Activation**: Implemented an account activation endpoint. For now, the activation token is returned in the response when a member is created (instead of being sent via email). They can be activated by sending a PUT request to /v1/members/:id/activate

## Getting Started

### Prerequisites
- Go 1.20+
- A database (PostgreSQL recommended)

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/workout-tracker-go.git
   cd workout-tracker-go
   ```
2. Install dependencies:
   ```bash
   go mod tidy
   ```
3. Set up the database:
   - Run the migration scripts in the `migrations` folder.
   - Update the database connection string in the `.env` file.

   Example `.env` file:
   ```bash
   DATABASE_URL=<PostgreSQL connection string>
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

### API Endpoints
#### General
- **Healthcheck**: `GET /v1/healthcheck`

#### Members
- `GET /v1/members?email=<member's email>`: Get a member by email.
- `POST /v1/members`: Create a new member.
- `PUT /v1/members/:id`: Update a member's details.
- `DELETE /v1/members/:id`: Delete a member.
- `PUT /v1/members/:id/activate`: Activate a member account.

#### Authentication
- `POST /v1/tokens/authentication`: Obtain a JWT by sending user credentials.

#### Exercises
- `GET /v1/exercises?category=<exercise category>`: Get exercises by category.
- `POST /v1/exercises`: Create a new exercise.
- `PUT /v1/exercises/:id`: Update an exercise.
- `DELETE /v1/exercises/:id`: Delete an exercise.

#### Workouts
- `GET /v1/members/:id/workouts`: Get all workouts for a member.
- `POST /v1/members/:id/workouts`: Create a new workout for a member.
- `DELETE /v1/members/:id/workouts/:workout_id`: Delete a workout.

## Project Structure
```plaintext
workout-tracker-go/
│
├── cmd/
│   └── api/          # API Handlers and Routes
│
├── internal/
│   └── data/         # Data Models and Database Logic
│
├── migrations/       # SQL Migration Files
│
├── go.mod            # Go Module Configuration
├── go.sum            # Go Module Dependencies
└── .git              # Environment Variable Storage
```

## Future Enhancements
- Filtered search for workouts.
- Redis for caching.
- Improved logging.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License
This project is licensed under the MIT License.

## Contact
For questions or feedback, reach out to **[Your Name](mailto:ilijakrilovic@gmail.com)**.
