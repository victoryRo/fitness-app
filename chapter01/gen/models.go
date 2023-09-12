// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package chapter01

import (
	"encoding/json"
	"time"
)

type GowebappExercise struct {
	ExerciseID   int64  `db:"exercise_id" json:"exerciseId"`
	ExerciseName string `db:"exercise_name" json:"exerciseName"`
}

type GowebappImage struct {
	ImageID     int64  `db:"image_id" json:"imageId"`
	UserID      int64  `db:"user_id" json:"userId"`
	ContentType string `db:"content_type" json:"contentType"`
	ImageData   []byte `db:"image_data" json:"imageData"`
}

type GowebappSet struct {
	SetID      int64 `db:"set_id" json:"setId"`
	ExerciseID int64 `db:"exercise_id" json:"exerciseId"`
	Weight     int32 `db:"weight" json:"weight"`
}

type GowebappUser struct {
	UserID       int64           `db:"user_id" json:"userId"`
	UserName     string          `db:"user_name" json:"userName"`
	PassWordHash string          `db:"pass_word_hash" json:"passWordHash"`
	Name         string          `db:"name" json:"name"`
	Config       json.RawMessage `db:"config" json:"config"`
	CreatedAt    time.Time       `db:"created_at" json:"createdAt"`
	IsEnabled    bool            `db:"is_enabled" json:"isEnabled"`
}

type GowebappWorkout struct {
	WorkoutID  int64     `db:"workout_id" json:"workoutId"`
	SetID      int64     `db:"set_id" json:"setId"`
	UserID     int64     `db:"user_id" json:"userId"`
	ExerciseID int64     `db:"exercise_id" json:"exerciseId"`
	StartDate  time.Time `db:"start_date" json:"startDate"`
}
