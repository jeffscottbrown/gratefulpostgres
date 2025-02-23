# Use the official PostgreSQL image
FROM postgres:latest

# Set environment variables for default user, password, and database
# ENV POSTGRES_USER=postgres
# ENV POSTGRES_PASSWORD=mysecretpassword
# ENV POSTGRES_DB=shows

# Copy initialization SQL file into the container
COPY sql/data.sql /docker-entrypoint-initdb.d/data.sql

# Expose PostgreSQL port
EXPOSE 5432
