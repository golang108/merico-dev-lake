module.exports = {
  // Configuration of lake's own services
  lake: {
    // Enable basic authentication to the lake API
    // token: 'mytoken'
    // Set how often does lake fetch new data from data sources, defaults to every hour
  },
  // Configuration of MongoDB
  mongo: {
    connectionString: 'mongodb://lake:lakeIScoming@localhost:27017/lake?authSource=admin'
  },
  // Configuration of rabbitMQ
  rabbitMQ: {
    connectionString: 'amqp://guest:guestWhat@localhost:5672/rabbitmq'
  },
  // Configuration of PostgreSQL
  postgres: {
    username: 'postgres',
    password: 'postgresWhat',
    host: 'localhost',
    database: 'lake',
    port: 5432,
    dialect: 'postgres'
  },
  cron: {
    // uncomment and update following configuration to enable the cron job
    /*
    job: {
      jira: {
        // boardId: 123
      },
      gitlab: {
        projectId: 123
      }
    },
    */
    interval: 5000,
    loopIntervalInMinutes: 60
  }
}
