CREATE SCHEMA IF NOT EXISTS app;

CREATE TABLE app.users (
    User_ID        BIGSERIAL PRIMARY KEY,
    User_Name      text NOT NULL,
    Pass_Word_Hash text NOT NULL,
    Name           text NOT NULL,
    Config         JSONB DEFAULT '{}'::JSONB NOT NULL,
    Create_At      TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    Is_Enabled     BOOLEAN DEFAULT TRUE NOT NULL
);

-- SQLc converts snake_case to CamelCase automatically
-- when converting PostgreSQL schema definitions to Go structs and interfaces.
-- This allows the Go code to adhere to the standard Go naming conventions
-- while still using snake_case in the database schema for readability.

-- Simple image blob demo
CREATE TABLE app.images (
    Image_ID     BIGSERIAL PRIMARY KEY,
    User_ID      BIGINT NOT NULL ,
    Content_Type TEXT NOT NULL DEFAULT 'image/png',
    Image_Data   BYTEA NOT NULL
);

CREATE TABLE app.exercices(
    Exercice_ID   BIGSERIAL PRIMARY KEY,
    Exercice_Name text NOT NULL,
);

CREATE TABLE app.workouts (
    Workout_ID BIGSERIAL PRIMARY KEY,
    User_ID    BIGINT NOT NULL,
    Set_ID     BIGINT NOT NULL,
    Start_Date TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE app.sets (
    Set_ID      BIGSERIAL PRIMARY KEY,
    Exercise_ID BIGINT NOT NULL,
    Weight      INT NOT NULL DEFAULT 0
    -- this can go up in decimal amounts so we can just divide/multiply to stay easy
);