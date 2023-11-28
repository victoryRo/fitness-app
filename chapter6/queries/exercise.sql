-- name: CreateUserExercise :one
INSERT INTO app.exercises (
    User_ID,
    Exercise_Name
) VALUES (
    $1,
    $2
) ON CONFLICT (Exercise_Name) DO NOTHING RETURNING (
    User_ID, Exercise_Name
);

-- name: ListUserExercises :many
SELECT Exercise_Name
FROM app.exercises
WHERE User_ID = $1;

-- name: DeleteUserExercise :exec
DELETE FROM app.exercises
WHERE User_ID = $1 AND Exercise_Name = $2;


-- name: CreateUserDefaultExercise :exec
INSERT INTO app.exercises (
    User_ID,
    Exercise_Name
) VALUES (
    1,
    'Bench Press'
),(
    1,
    'Barbell Row'
);

-- name: CreateUserWorkout :one
INSERT INTO app.workouts (
    User_ID,
    Start_Date
) VALUES (
    $1,
    NOW()
) RETURNING *;

-- name: CreateDefaultSetForExercise :one
INSERT INTO app.sets (
    Workout_ID,
    Exercise_Name,
    Weight
) VALUES (
    $1,
    $2,
    $3
) RETURNING *;

-- name: CreateSetForExercise :one
INSERT INTO app.sets (
    Workout_ID,
    Exercise_Name, 
    Weight,
    Set1,
    Set2,
    Set3
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) RETURNING *;

-- name: UpdateSet :one
UPDATE app.sets SET
    Weight = $1,
    Set1 = $2,
    Set2 = $3,
    Set3 = $4
WHERE Set_ID = $5 AND Workout_ID = $6 RETURNING *;

-- name: GetWorkoutsForUserID :many
SELECT w.Workout_ID, COALESCE(s.Set_ID,-1), COALESCE(s.name,''), COALESCE(s.set1,-1), COALESCE(s.set1,-1), COALESCE(s.set2,-1), COALESCE(s.set3,-1), COALESCE(s.weight,-1), w.Start_Date AS date FROM
(
SELECT Set_ID, Workout_ID, Exercise_Name as name, set1, set2, set3, weight FROM app.sets
) AS s RIGHT JOIN app.workouts AS w USING (Workout_ID)
WHERE w.User_ID = $1
ORDER BY date DESC;


-- name: DeleteWorkoutByIDForUser :exec
DELETE FROM app.workouts

WHERE User_ID = $1 AND Workout_ID = $2;

