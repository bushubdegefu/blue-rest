# Blue Rest - CLI Tool for Golang Backend Projects

Blue Rest is a command-line tool designed to structure Golang backend projects with GORM and Fiber/Echo for SQL-based applications. It also generates OpenAPI documentation using Golang’s Swag.

## Installation

To install Blue Rest, run:

```bash
 go install github.com/bushubdegefu/blue-rest@latest
```
## Folder Structure Overview

The Blue Rest CLI tool will generate a structured folder layout for your Golang backend project. Below is a sample folder structure:

```
fiber-sample/
├── manager/
│   ├── deveecho.go
│   ├── consumer.go
│   ├── manager.go
│   └── migrate.go
├── helper/
├── messages/
├── observe/
├── config/
├── database/
├── bluetasks/
├── testsettings/
├── tests/
└── controllers/
```

### Folder Descriptions

- **manager/**: Contains the Cobra app and commands to start the application, manage consumers, and handle database migrations.
- **helper/**: Utility functions that assist with various tasks across the project.
- **messages/**: Defines the producer and consumer connections and functions for RabbitMQ.
- **observe/**: Contains OpenTelemetry (OTel) Spanner functions for observability.
- **config/**: Functions to manage environment variables defined in `.env` files.
- **database/**: Database connection and configuration files.
- **bluetasks/**: Scheduled tasks such as clearing logs.
- **testsettings/**: Setup and configuration for testing.
- **tests/**: Contains tests for CRUD operations.
- **controllers/**: Model handler functions for CRUD operations and other functionalities. Custom handlers should also be added here.

This basic structure is intended to serve as a project starter and should be edited based on your specific needs.


## Usage

### Initialize a New Project

For a Golang Fiber project, create a new directory and initialize your project:

```bash
 mkdir fiber-rest
 cd fiber-rest
 blue-rest init -name=github.com/username/fiber-rest
```

or using shorthand:

```bash
 blue-rest init -n github.com/username/fiber-rest
```

This command creates the Go module and a sample `config.json` file that defines your project’s attributes and data models for GORM annotations.
 In the `config.json` file, the `crud` flag is used to specify which CRUD operations should include a particular field. The flag is a string with six boolean values separated by `$`, representing the following operations:

1. **Get**: Include the field in GET requests.
2. **Post**: Include the field in POST requests.
3. **Patch**: Include the field in PATCH requests.
4. **Put**: Include the field in PUT requests.
5. **OtM**: Indicates if the field is a foreign key or a one-to-many relationship.
6. **MtM**: Indicates if the field is a many-to-many relationship.

For example, if the `crud` flag for a field `id` is `true$false$false$false$true$false`, it means:
- The field should be included in GET requests.
- The field should not be included in POST, PATCH, or PUT requests.
- The field is a foreign key or a one-to-many relationship.
- The field is not a many-to-many relationship.

Additionally, in the `config.json` file, you can specify relationships for models using the `rln_model` field. This field should contain a list of relationships in the format `["ModelName$RelationshipType"]`. The supported relationship types are `otm` (one-to-many) and `mtm` (many-to-many). For each specified relationship, endpoints will be generated to add, remove, and get related entities.

For example, if a model `User` has the following relationships:
```json
"rln_model": ["Role$otm", "Group$mtm"]
```
This configuration will generate endpoints to manage:
- One-to-many relationships between `User` and `Role` (e.g., add, remove, and get roles for a user).
- Many-to-many relationships between `User` and `Group` (e.g., add, remove, and get groups for a user).
Additionally, you can specify the association table name for many-to-many relationships in the `rln_model` field. This is done by appending the table name after the relationship type, separated by a dollar sign. For example, `"Group$mtm$user_groups"` indicates a many-to-many relationship between `User` and `Group` using the `user_groups` table to join and form the SQL query.

For example, if a model `User` has the following relationships:
```json
"rln_model": ["Role$otm", "Group$mtm$user_groups"]
```
This configuration will generate endpoints to manage:
- One-to-many relationships between `User` and `Role` (e.g., add, remove, and get roles for a user).
- Many-to-many relationships between `User` and `Group` using the `user_groups` table (e.g., add, remove, and get groups for a user).
> **Note:** The `user_groups` table is an association table between `User` and `Group` models, used to manage many-to-many relationships for the above example.
### Generate Basic Components

Generate different boilerplate components:

#### Database Setup

```bash
 blue-rest basic --type db
```

or

```bash
 blue-rest basic -t db
```

#### Pagination Support

```bash
 blue-rest basic --type pagination
```

or

```bash
 blue-rest basic -t pagination
```

#### Message Producer
This command generates a basic structure for a RabbitMQ key-value store. It creates helper functions for both producers and consumers using RabbitMQ. You should adjust these functions based on your specific needs, as they are provided as a "starter structure."

```bash
 blue-rest basic --type producer
```

or

```bash
 blue-rest basic -t producer
```

#### Task Management

```bash
 blue-rest basic --type tasks
```

or

```bash
 blue-rest basic -t tasks
```

### Generate Models

```bash
 blue-rest models
```

Fix any dependencies from the generated models and format the code.

### Generate Controller Endpoints

```bash
 blue-rest crud --type fiber
```

or

```bash
 blue-rest crud -t fiber
```

### Enable Fiber Integration

```bash
 blue-rest fiber
```

### Run Migrations
Run The cobra command line migration template generation flag.

```bash
 blue-rest migration
```

After fixing dependencies, run:

```bash
 go mod tidy
```

For demo purposes, you can change the type of database you want to use. The supported databases are SQLite3, PostgreSQL, and MySQL. Update the configuration accordingly and run:

```bash
 go run main.go run migrate
```
### Generate Tests for CRUD Operations

Finally, you can generate tests for the CRUD operations of entities using the command:

```bash
 blue-rest test -f fiber
```

Note that these generated tests may still need some work after editing, but they provide a good starting point.

Then start the server:

```bash
 go run main.go run --env=dev
```

Finally, check your API on `localhost:7500` or the port specified in `.dev.env`.

### Check Available Commands

You can always check the available commands using:

```bash
 blue-rest help
```

### Available Commands

Blue Rest – command-line tool to aid in structuring your Golang backend projects with GORM and Fiber/Echo for SQL-based projects.

#### Usage:
  ```bash
  Blue [command]
  ```

#### Commands:
- **basic**       Generate a basic folder structure for a project.
- **completion**  Generate the autocompletion script for the specified shell.
- **crud**        Generate CRUD handlers based on GORM for the specified framework.
- **echo**        Generate the basic structure file to start an app using Echo.
- **fiber**       Generate basic structure files to start an app using Fiber.
- **help**        Help about any command.
- **init**        Initialize the module with a name.
- **migration**   Generate data models based on GORM using the provided spec in the `config.json` file.
- **models**      Generate data models based on GORM with annotations.
- **version**     Print out the version of the CLI app in the terminal.

#### Flags:
- **-h, --help**      Help for Blue.
- **-v, --version**   Version for Blue.

Use `Blue [command] --help` for more information about a command.

## Major Caveats

- You should always review and adjust the generated code to fit your requirements.
- Middleware configurations need to be adjusted as per your application's needs. Currently, the middleware validates everything.

## Coming Soon
- Support for MongoDB.
- Better support for tests and tests for relationships.

## Echo Boilerplate Generation

Repeat the same steps for Echo-based boilerplate generation by replacing the tags that have `fiber` in the bash commands with `echo`.

## License

This project is licensed under the MIT License.
