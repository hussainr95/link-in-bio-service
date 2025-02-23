# link-in-bio-service
A take home assignment where an API should allow users to create, update, delete, and retrieve bio links.

## Quick Start

1. **Install Dependencies**

   Run the following command from the project root:
   ```bash
   go mod tidy

2. **Generate Swagger Documentation**

   This command generates the Swagger docs (creates the docs/ folder):
    ```bash
    make swag
   ```
   Note: Ensure you have the Swagger tool installed via:
   ```bash
    make swag-install
   ```

3. **Run Unit Tests**

   To execute all unit tests:
    ```bash
    make test
   
4. **Manage Docker Containers**

   Build and start the application (Go API + MongoDB) with:
    ```bash
   make docker-up
   ```
   Your API will be available at http://localhost:8080.

   To stop the service:
    ```bash
   make docker-down
   ```

5. **Access Swagger UI**
   Once the containers are running, open your browser and navigate to:
    ```bash
   http://localhost:8080/swagger/index.html
   ```
   Use the Authorize button to enter the Bearer token value Bearer test
