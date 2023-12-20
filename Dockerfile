# Use the official PostgreSQL image from Docker Hub
FROM postgres:latest

# Environment variables for PostgreSQL
ENV POSTGRES_USER     ""
ENV POSTGRES_PASSWORD ""
ENV POSTGRES_DB       ""

# Expose the PostgreSQL port
EXPOSE 5432

# Start PostgreSQL service
CMD ["postgres"]

