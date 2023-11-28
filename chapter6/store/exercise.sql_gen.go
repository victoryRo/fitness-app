// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: exercise.sql

package store

import (
	"context"
	"time"
)

const createDefaultSetForExercise = `-- name: CreateDefaultSetForExercise :one
INSERT INTO app.sets (
    Workout_ID,
    Exercise_Name,
    Weight
) VALUES (
    $1,
    $2,
    $3
) RETURNING set_id, workout_id, exercise_name, weight, set1, set2, set3
`

type CreateDefaultSetForExerciseParams struct {
	WorkoutID    int64  `json:"workout_id"`
	ExerciseName string `json:"exercise_name"`
	Weight       int32  `json:"weight"`
}

func (q *Queries) CreateDefaultSetForExercise(ctx context.Context, arg CreateDefaultSetForExerciseParams) (AppSet, error) {
	row := q.db.QueryRowContext(ctx, createDefaultSetForExercise, arg.WorkoutID, arg.ExerciseName, arg.Weight)
	var i AppSet
	err := row.Scan(
		&i.SetID,
		&i.WorkoutID,
		&i.ExerciseName,
		&i.Weight,
		&i.Set1,
		&i.Set2,
		&i.Set3,
	)
	return i, err
}

const createSetForExercise = `-- name: CreateSetForExercise :one
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
) RETURNING set_id, workout_id, exercise_name, weight, set1, set2, set3
`

type CreateSetForExerciseParams struct {
	WorkoutID    int64  `json:"workout_id"`
	ExerciseName string `json:"exercise_name"`
	Weight       int32  `json:"weight"`
	Set1         int64  `json:"set1"`
	Set2         int64  `json:"set2"`
	Set3         int64  `json:"set3"`
}

func (q *Queries) CreateSetForExercise(ctx context.Context, arg CreateSetForExerciseParams) (AppSet, error) {
	row := q.db.QueryRowContext(ctx, createSetForExercise,
		arg.WorkoutID,
		arg.ExerciseName,
		arg.Weight,
		arg.Set1,
		arg.Set2,
		arg.Set3,
	)
	var i AppSet
	err := row.Scan(
		&i.SetID,
		&i.WorkoutID,
		&i.ExerciseName,
		&i.Weight,
		&i.Set1,
		&i.Set2,
		&i.Set3,
	)
	return i, err
}

const createUserDefaultExercise = `-- name: CreateUserDefaultExercise :exec
INSERT INTO app.exercises (
    User_ID,
    Exercise_Name
) VALUES (
    1,
    'Bench Press'
),(
    1,
    'Barbell Row'
)
`

func (q *Queries) CreateUserDefaultExercise(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, createUserDefaultExercise)
	return err
}

const createUserExercise = `-- name: CreateUserExercise :one
INSERT INTO app.exercises (
    User_ID,
    Exercise_Name
) VALUES (
    $1,
    $2
) ON CONFLICT (Exercise_Name) DO NOTHING RETURNING (
    User_ID, Exercise_Name
)
`

type CreateUserExerciseParams struct {
	UserID       int64  `json:"user_id"`
	ExerciseName string `json:"exercise_name"`
}

func (q *Queries) CreateUserExercise(ctx context.Context, arg CreateUserExerciseParams) (interface{}, error) {
	row := q.db.QueryRowContext(ctx, createUserExercise, arg.UserID, arg.ExerciseName)
	var column_1 interface{}
	err := row.Scan(&column_1)
	return column_1, err
}

const createUserWorkout = `-- name: CreateUserWorkout :one
INSERT INTO app.workouts (
    User_ID,
    Start_Date
) VALUES (
    $1,
    NOW()
) RETURNING workout_id, user_id, start_date
`

func (q *Queries) CreateUserWorkout(ctx context.Context, userID int64) (AppWorkout, error) {
	row := q.db.QueryRowContext(ctx, createUserWorkout, userID)
	var i AppWorkout
	err := row.Scan(&i.WorkoutID, &i.UserID, &i.StartDate)
	return i, err
}

const deleteUserExercise = `-- name: DeleteUserExercise :exec
DELETE FROM app.exercises
WHERE User_ID = $1 AND Exercise_Name = $2
`

type DeleteUserExerciseParams struct {
	UserID       int64  `json:"user_id"`
	ExerciseName string `json:"exercise_name"`
}

func (q *Queries) DeleteUserExercise(ctx context.Context, arg DeleteUserExerciseParams) error {
	_, err := q.db.ExecContext(ctx, deleteUserExercise, arg.UserID, arg.ExerciseName)
	return err
}

const deleteWorkoutByIDForUser = `-- name: DeleteWorkoutByIDForUser :exec
DELETE FROM app.workouts

WHERE User_ID = $1 AND Workout_ID = $2
`

type DeleteWorkoutByIDForUserParams struct {
	UserID    int64 `json:"user_id"`
	WorkoutID int64 `json:"workout_id"`
}

func (q *Queries) DeleteWorkoutByIDForUser(ctx context.Context, arg DeleteWorkoutByIDForUserParams) error {
	_, err := q.db.ExecContext(ctx, deleteWorkoutByIDForUser, arg.UserID, arg.WorkoutID)
	return err
}

const getWorkoutsForUserID = `-- name: GetWorkoutsForUserID :many
SELECT w.Workout_ID, COALESCE(s.Set_ID,-1), COALESCE(s.name,''), COALESCE(s.set1,-1), COALESCE(s.set1,-1), COALESCE(s.set2,-1), COALESCE(s.set3,-1), COALESCE(s.weight,-1), w.Start_Date AS date FROM
(
SELECT Set_ID, Workout_ID, Exercise_Name as name, set1, set2, set3, weight FROM app.sets
) AS s RIGHT JOIN app.workouts AS w USING (Workout_ID)
WHERE w.User_ID = $1
ORDER BY date DESC
`

type GetWorkoutsForUserIDRow struct {
	WorkoutID int64     `json:"workout_id"`
	SetID     int64     `json:"set_id"`
	Name      string    `json:"name"`
	Set1      int64     `json:"set1"`
	Set1_2    int64     `json:"set1_2"`
	Set2      int64     `json:"set2"`
	Set3      int64     `json:"set3"`
	Weight    int32     `json:"weight"`
	Date      time.Time `json:"date"`
}

func (q *Queries) GetWorkoutsForUserID(ctx context.Context, userID int64) ([]GetWorkoutsForUserIDRow, error) {
	rows, err := q.db.QueryContext(ctx, getWorkoutsForUserID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetWorkoutsForUserIDRow{}
	for rows.Next() {
		var i GetWorkoutsForUserIDRow
		if err := rows.Scan(
			&i.WorkoutID,
			&i.SetID,
			&i.Name,
			&i.Set1,
			&i.Set1_2,
			&i.Set2,
			&i.Set3,
			&i.Weight,
			&i.Date,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listUserExercises = `-- name: ListUserExercises :many
SELECT Exercise_Name
FROM app.exercises
WHERE User_ID = $1
`

func (q *Queries) ListUserExercises(ctx context.Context, userID int64) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, listUserExercises, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var exercise_name string
		if err := rows.Scan(&exercise_name); err != nil {
			return nil, err
		}
		items = append(items, exercise_name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateSet = `-- name: UpdateSet :one
UPDATE app.sets SET
    Weight = $1,
    Set1 = $2,
    Set2 = $3,
    Set3 = $4
WHERE Set_ID = $5 AND Workout_ID = $6 RETURNING set_id, workout_id, exercise_name, weight, set1, set2, set3
`

type UpdateSetParams struct {
	Weight    int32 `json:"weight"`
	Set1      int64 `json:"set1"`
	Set2      int64 `json:"set2"`
	Set3      int64 `json:"set3"`
	SetID     int64 `json:"set_id"`
	WorkoutID int64 `json:"workout_id"`
}

func (q *Queries) UpdateSet(ctx context.Context, arg UpdateSetParams) (AppSet, error) {
	row := q.db.QueryRowContext(ctx, updateSet,
		arg.Weight,
		arg.Set1,
		arg.Set2,
		arg.Set3,
		arg.SetID,
		arg.WorkoutID,
	)
	var i AppSet
	err := row.Scan(
		&i.SetID,
		&i.WorkoutID,
		&i.ExerciseName,
		&i.Weight,
		&i.Set1,
		&i.Set2,
		&i.Set3,
	)
	return i, err
}
