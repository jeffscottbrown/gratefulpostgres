# Use the official PostgreSQL image
FROM postgres:latest

# Accept environment variables for PostgreSQL
ARG POSTGRES_PASSWORD
ARG POSTGRES_DB

# Set environment variables
ENV POSTGRES_PASSWORD=$POSTGRES_PASSWORD
ENV POSTGRES_DB=$POSTGRES_DB
ENV POSTGRES_USER=postgres

# Copy initialization SQL file into the container
COPY sql/data.sql /docker-entrypoint-initdb.d/data.sql

# Expose PostgreSQL port
EXPOSE 5432
