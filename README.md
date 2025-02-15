# Distributed scheduler

Welcome to the Distributed Scheduler repository!
This system, consisting of a Management API 🛠️ and a Runner service 🏃‍♀️, allows you to easily schedule and manage jobs
that execute at specified times in the future.

## Features

- **Job Scheduling**: Schedule jobs to run at specific times in the future.
    - **One-Time and Recurring Jobs**: Schedule jobs to run once or on a recurring basis.
    - **Cron Syntax**: Use cron syntax to schedule recurring jobs.
    - **HTTP or AMQP Jobs**: Send messages to an HTTP endpoint or an AMQP queue.
- **Job Management**: View, update, and delete jobs.

## Roadmap

- [ ] **Limit number of job executions**: Limit the number of times a job can be executed.
- [ ] **Job Dependencies**: Allow jobs to depend on other jobs.
- [ ] **Job Priorities**: Allow jobs to be assigned priorities.
- [x] **Job Retries**: Allow jobs to be retried if they fail.
- [ ]  **Job callbacks**: Allow jobs to call a specified endpoint after completion.

## Quickstart

You must have Docker and Docker Compose installed on your machine to run the system locally.

1. Clone the repository:

    ```bash
    git clone https://github.com/TimeSnap/distributed-scheduler.git
    
    cd distributed-scheduler
    ```

2. Run the system using Docker Compose:

    ```bash
    docker compose -f docker-compose.yml up
    ```

## Configuration

Check out our detailed [Local Development Guide](./documentation/development.md)
and [Configuration Guide](./documentation/configuration.md).

## Architecture

Check out our detailed [Architecture Overview](./documentation/architecture.md) in the Documentation directory to learn
more about the system's design.

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) to learn more about how
to get involved.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.