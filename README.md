# Workout Tracker API

A **Workout Tracker** application built with Go, designed to manage workout routines, exercises, and member information. The API provides endpoints to handle CRUD operations for exercises, workouts, and members.

## Features
- **Exercise Management**: Add, update, delete, and fetch exercises.
- **Workout Tracking**: Manage workouts and their details.
- **Member Support**: Create and manage member profiles.
- **Database Migrations**: Predefined scripts for database setup.

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
   - Update the database connection string in env file.
   - 
   env file example:
   ```bash
   DATABASE_URL=<PostgreSQL connection string>
   ```

4. Run the application:
   ```bash
   go run cmd/api/main.go
   ```

### API Endpoints
- **Healthcheck**: `GET /healthcheck`
- **Members**:
  - `GET /members?email=<member's email>`
  - `POST /members`
  - `PUT /members/id`
  - `DELETE /memebrs/id`
- **Exercises**:
  - `GET /exercises?category=<exercise category>`
  - `POST /exercises`
  - `PUT /exercises/id`
  - `DELETE /exercises/id`
- **Workouts**:
  - `GET /member/id/workouts`

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
└── .git              # Storing Environment Variables
```

## Future Enhancements
- Authentication and Authorization for users.
- Advanced analytics for workout tracking.

## Contributing
Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License
This project is licensed under the MIT License.

## Contact
For questions or feedback, reach out to **[Your Name](mailto:youremail@example.com)**.

---

You can adjust the placeholder information, such as the repository URL, email, and license, to match your project. Let me know if you'd like further assistance!
