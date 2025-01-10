# Dottie MVP API

## Overview

Dottie MVP API is a backend service designed to analyze menstrual symptoms and provide actionable insights based on the ACOG guidelines. The service includes endpoints for symptom analysis, educational content retrieval, and user management.

## Goals

- **Symptom Checker**: Evaluate user inputs to determine if menstrual symptoms are normal or abnormal.
- **Abnormality Identification**: Identify menstrual abnormalities such as amenorrhea or oligomenorrhea based on user input.
- **Educational Content**: Provide educational content related to identified conditions.
- **User Management**: Handle user registration, login, update, and deletion.

## Setup Instructions

### Prerequisites

- Python 3.8+
- Neo4j Aura account
- Virtual environment tool (e.g., `venv`)

### Installation

1. **Clone the repository**:

    ```sh
    git clone https://github.com/yourusername/dottie-mvp-api.git
    cd dottie-mvp-api
    ```

2. **Create and activate a virtual environment**:

    ```sh
    python -m venv venv
    source venv/bin/activate  # On Windows use `venv\Scripts\activate`
    ```

3. **Install dependencies**:

    ```sh
    pip install -r requirements.txt
    ```

4. **Set up environment variables**:

    Create a [.env] file in the root directory and add the following variables:

    ```env
    NEO4J_URI=your_neo4j_uri
    NEO4J_USER=your_neo4j_user
    NEO4J_PASSWORD=your_neo4j_password
    JWT_SECRET_KEY=your_jwt_secret_key
    ```

5. **Initialize the Neo4j database**:

    ```sh
    python -m app.db.neo4j_connector
    ```

6. **Run the application**:

    ```sh
    uvicorn app.main:app --reload
    ```

### Usage

- **Symptom Checker**:
  - **Endpoint**: POST `/api/v1/symptoms/check`
  - **Request Body**:
    ```json
    {
        "cycle_length": 28,
        "cycle_duration": 5,
        "symptoms": ["mild pain"]
    }
    ```
  - **Response**:
    ```json
    {
        "status": "Normal",
        "recommendation": "No action needed"
    }
    ```

- **Abnormality Identification**:
  - **Endpoint**: POST `/api/v1/symptoms/identify`
  - **Request Body**:
    ```json
    {
        "cycle_length": 50,
        "cycle_duration": 7,
        "symptoms": ["heavy bleeding"],
        "missed_periods": 2
    }
    ```
  - **Response**:
    ```json
    {
        "status": "Abnormal",
        "abnormalities": ["Oligomenorrhea", "Menorrhagia"],
        "recommendation": "Consult a healthcare provider for further evaluation."
    }
    ```

- **Educational Content**: POST `/api/v1/content/get_content`
- **User Management**: POST `/api/v1/users/register`, `/api/v1/users/login`, PUT `/api/v1/users/update`, DELETE `/api/v1/users/delete`

### Testing

Use Postman or a similar tool to test the endpoints. Example requests and responses are provided in the documentation.

### Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

### License

This project is licensed under the MIT License.

## Research Resources

https://www.acog.org/clinical/clinical-guidance/committee-opinion/articles/2015/12/menstruation-in-girls-and-adolescents-using-the-menstrual-cycle-as-a-vital-sign

### Resources

https://docs.google.com/document/d/1ywyorf9YV1EZs2CjX3UxEw-BmWEca7FdP60b9jQ9VT4/edit?tab=t.0#heading=h.8ahwpe5c65y4

https://docs.google.com/document/d/1Vqu8Sph3sco7j_tOZVmnkNW_YIm0MhBEwJKFxqTO9gU/edit?tab=t.0#heading=h.spbkur2dd5h2
