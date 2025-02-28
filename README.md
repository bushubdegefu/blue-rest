# Blue Rest - CLI Tool for Golang Backend Projects

Blue Rest is a command-line tool designed to structure Golang backend projects with GORM and Fiber/Echo for SQL-based applications. It also generates OpenAPI documentation using Golang’s Swag.

## Installation

To install Blue Rest, run:

```bash
 go install github.com/bushubdegefu/blue-rest@latest
```

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

```bash
 go run blue-rest migration
```

After fixing dependencies, run:

```bash
 go mod tidy
```

For demo purposes, you can change the type of database you want to use. The supported databases are SQLite3, PostgreSQL, and MySQL. Update the configuration accordingly and run:

```bash
 go run main.go run migrate
```

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

- Testing templates for better boilerplate code coverage.

## Echo Boilerplate Generation

Repeat the same steps for Echo-based boilerplate generation by replacing the tags that have `fiber` in the bash commands with `echo`.

## License

This project is licensed under the MIT License.
