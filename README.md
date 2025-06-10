# Hospital System API

A microservices-based hospital management system built with Go, gRPC, and Docker.

## Services

1. **Account Service** (Port 8080)
   - User authentication and authorization
   - Role-based access control
   - JWT token management

2. **Hospital Service** (Port 50051)
   - Hospital and department management
   - Room management
   - gRPC API

3. **Timetable Service** (Port 50052)
   - Appointment scheduling
   - Doctor availability management
   - gRPC API

4. **Document Service** (Port 50053)
   - Medical document management
   - Document search using Elasticsearch
   - gRPC API

## Prerequisites

- Docker
- Docker Compose
- Go 1.21 or later
- Protocol Buffers compiler (protoc)

## Default Users

The system comes with the following pre-configured users:

| Username | Password | Role    |
|----------|----------|---------|
| admin    | admin    | Admin   |
| manager  | manager  | Manager |
| doctor   | doctor   | Doctor  |
| user     | user     | User    |

## Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/sergeimurashev/hospital-system-api.git
   cd hospital-system-api
   ```

2. Build and run the services:
   ```bash
   docker-compose up -d
   ```

3. The services will be available at:
   - Account Service: http://localhost:8080
   - Hospital Service: localhost:50051
   - Timetable Service: localhost:50052
   - Document Service: localhost:50053
   - Kibana: http://localhost:5601

## API Documentation

### Account Service

#### Authentication
- POST /api/v1/auth/login
  - Request body: `{"username": "string", "password": "string"}`
  - Response: JWT token

### Hospital Service (gRPC)

```protobuf
service HospitalService {
  rpc CreateHospital(CreateHospitalRequest) returns (Hospital);
  rpc GetHospital(GetHospitalRequest) returns (Hospital);
  rpc UpdateHospital(UpdateHospitalRequest) returns (Hospital);
  rpc DeleteHospital(DeleteHospitalRequest) returns (Empty);
  rpc ListHospitals(ListHospitalsRequest) returns (ListHospitalsResponse);
}
```

### Timetable Service (gRPC)

```protobuf
service TimetableService {
  rpc CreateAppointment(CreateAppointmentRequest) returns (Appointment);
  rpc GetAppointment(GetAppointmentRequest) returns (Appointment);
  rpc UpdateAppointment(UpdateAppointmentRequest) returns (Appointment);
  rpc DeleteAppointment(DeleteAppointmentRequest) returns (Empty);
  rpc ListAppointments(ListAppointmentsRequest) returns (ListAppointmentsResponse);
}
```

### Document Service (gRPC)

```protobuf
service DocumentService {
  rpc CreateDocument(CreateDocumentRequest) returns (Document);
  rpc GetDocument(GetDocumentRequest) returns (Document);
  rpc UpdateDocument(UpdateDocumentRequest) returns (Document);
  rpc DeleteDocument(DeleteDocumentRequest) returns (Empty);
  rpc SearchDocuments(SearchDocumentsRequest) returns (SearchDocumentsResponse);
}
```

## Development

### Project Structure

```
.
├── account-service/
├── hospital-service/
├── timetable-service/
├── document-service/
├── proto/
├── docker-compose.yml
└── README.md
```

### Adding New Features

1. Define the service interface in the proto file
2. Generate Go code from proto files
3. Implement the service
4. Update the Docker configuration if needed
5. Test the changes

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 