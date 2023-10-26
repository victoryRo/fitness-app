-- name: ListUsers :many
SELECT *
FROM app.users
ORDER BY user_name;

-- name: ListImages :many
SELECT *
FROM app.images
ORDER BY image_id;

-- name: ListExercises :many
SELECT *
FROM app.exercises
ORDER BY exercise_name;

-- name: ListSets :many
SELECT *
FROM app.sets
ORDER BY weight;

-- name: ListWorkouts :many
SELECT *
FROM app.workouts
ORDER BY workout_id;

-- name: GetUser :one
SELECT *
FROM app.users
WHERE user_id = $1;

-- name: GetUserByName :one
SELECT *
FROM app.users
WHERE user_name = $1;

-- name: GetUserWorkout :many
SELECT u.user_id, w.workout_id, w.start_date, w.set_id
FROM app.users u,
     app.workouts w
WHERE u.user_id = w.user_id
  AND u.user_id = $1;

-- name: GetUserSets :many
SELECT u.user_id, w.workout_id, w.start_date, s.set_id, s.weight
FROM app.users u,
     app.workouts w,
     app.sets s
WHERE u.user_id = w.user_id
  AND w.set_id = s.set_id
  AND u.user_id = $1;

-- name: GetUserImage :one
SELECT u.name, u.user_id, i.image_data
FROM app.users u,
     app.images i
WHERE u.user_id = i.user_id
  AND u.user_id = $1;

-- name: DeleteUsers :exec
DELETE
FROM app.users
WHERE user_id = $1;

-- name: DeleteUserImage :exec
DELETE
FROM app.images i
WHERE i.user_id = $1;

-- name: DeleteUserWorkouts :exec
DELETE
FROM app.workouts w
WHERE w.user_id = $1;

-- name: DeleteExercise :exec
DELETE
FROM app.exercises e
WHERE e.exercise_id = $1;

-- name: DeleteSets :exec
DELETE
FROM app.sets s
WHERE s.set_id = $1;

-- name: CreateExercise :one
INSERT INTO app.exercises (Exercise_Name)
values ($1) RETURNING Exercise_ID;

-- name: UpsertExercise :one
INSERT INTO app.exercises (Exercise_Name)
    VALUES ($1) ON CONFLICT (Exercise_ID) DO
    UPDATE
        SET Exercise_Name = EXCLUDED.Exercise_Name
        RETURNING Exercise_ID;

-- name: CreateUserImage :one
INSERT INTO app.images (User_ID, Content_Type, Image_Data)
    values ($1,
            $2,
            $3) RETURNING *;

-- name: UpsertUserImage :one
INSERT INTO app.images (Image_Data)
VALUES ($1) ON CONFLICT (Image_ID) DO
    UPDATE
        SET Image_Data = EXCLUDED.Image_Data
        RETURNING Image_ID;


-- name: CreateSet :one
INSERT INTO app.sets (Exercise_Id, Weight)
values ($1,
        $2) RETURNING *;

-- name: UpsertSet :one
INSERT INTO app.sets (Exercise_Id, Weight)
values ($1,
        $2) ON CONFLICT (Set_ID) DO
    UPDATE
        SET Exercise_Id = EXCLUDED.Exercise_Id, Weight = EXCLUDED.Weight
        RETURNING Set_ID;

-- name: CreateWorkout :one
INSERT INTO app.workouts (User_ID, Set_ID, Start_Date)
values ($1,
        $2,
        $3) RETURNING *;

-- name: UpsertWorkout :one
INSERT INTO app.workouts (User_ID, Set_ID, Start_Date)
    values ($1,
            $2,
            $3) ON CONFLICT (Workout_ID) DO
    UPDATE
        SET User_ID = EXCLUDED.User_ID,
        Set_ID = EXCLUDED.Set_ID,
        Start_Date = EXCLUDED.Start_Date
        RETURNING Workout_ID;

-- name: CreateUsers :one
INSERT INTO app.users (User_Name, Pass_Word_Hash, name)
VALUES ($1,
        $2,
        $3) RETURNING *;
