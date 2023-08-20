package webcam

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmoiron/sqlx"
	"github.com/quarkey/iot/pkg/helper"
)

type Camera struct {
	ID              int    `db:"id" json:"id"`
	SensorID        int    `db:"sensor_id" json:"sensor_id"`
	Title           string `db:"title" json:"title"`
	Description     string `db:"description" json:"description"`
	DescriptionLong string `db:"description_long" json:"description_long"`

	Hostname        string           `db:"hostname" json:"hostname"`
	ProjectName     string           `db:"project_name" json:"project_name"`
	StorageLocation string           `db:"storage_location" json:"storage_location"`
	Interval        int              `db:"interval" json:"interval"`
	NextCapTureTime pgtype.Timestamp `db:"next_capture_time" json:"next_capture_time"`
	Status          string           `db:"status" json:"status"`
	Alert           bool             `db:"alert" json:"alert"`
	Active          bool             `db:"active" json:"active"`

	CreatedAt string            `db:"created_at" json:"created_at"`
	UpdatedAt helper.NullString `db:"updated_at" json:"updated_at"`

	DB        *sqlx.DB `json:"-"`
	lastframe []byte
}

// NewCamera instantiates a new camera struct with db connection.
func NewCameraWithDB(db *sqlx.DB) Camera {
	return Camera{
		DB: db,
	}
}

func (c *Camera) InsertNewCamera() (Camera, error) {
	var returningID int
	err := c.DB.QueryRow(`insert into 
	camera(
		sensor_id,
		title,
		description,
		description_long,
		hostname,
		project_name,
		storage_location,
		interval,
		next_capture_time,
		alert,
		active
	) values($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id`,
		c.SensorID,
		c.Title,
		c.Description,
		c.DescriptionLong,
		c.Hostname,
		c.ProjectName,
		c.StorageLocation,
		c.Interval,
		c.NextCapTureTime,
		c.Alert,
		c.Active).Scan(&returningID)
	if err != nil && err != sql.ErrNoRows {
		return Camera{}, err
	}
	newCameraItem, err := GetCameraByID(returningID, c.DB)
	if err != nil {
		return Camera{}, err
	}
	return newCameraItem, nil
}

func (c *Camera) UpdateCamera() (Camera, error) {
	_, err := c.DB.Exec(`update camera set
		sensor_id=$1,
		title=$2,
		description=$3,
		description_long=$4,
		hostname=$5,
		project_name=$6,
		storage_location=$7,
		interval=$8,
		next_capture_time=$9,
		alert=$10,
		active=$11,
		updated_at=now()
		where id=$12`,
		c.SensorID,
		c.Title,
		c.Description,
		c.DescriptionLong,
		c.Hostname,
		c.ProjectName,
		c.StorageLocation,
		c.Interval,
		c.NextCapTureTime,
		c.Alert,
		c.Active,
		c.ID)
	if err != nil {
		return Camera{}, err
	}
	newCameraItem, err := GetCameraByID(c.ID, c.DB)
	if err != nil {
		return Camera{}, err
	}
	return newCameraItem, nil
}

func (c *Camera) DeleteCamera() (int64, error) {
	res, err := c.DB.Exec(`delete from camera where id=$1`, c.ID)
	if err != nil {
		return 0, fmt.Errorf("DeleteCamera(): %v", err)
	}
	nrows, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("DeleteCamera(): %v", err)
	}
	return nrows, nil
}

func GetCameraList(db *sqlx.DB) ([]*Camera, error) {
	var cameras []*Camera
	err := db.Select(&cameras, `
	select 
	a.id, 
	a.sensor_id,
	a.title,
	a.description, 
	a.description_long,
	a.hostname,
	a.project_name,
	a.storage_location,
	a.interval,
	a.next_capture_time,
	a.status,
	a.alert,
	a.active,
	a.created_at,
	a.updated_at
from camera a`)
	if err != nil {
		return nil, fmt.Errorf("GetCameras(): %v", err)
	}
	return cameras, nil
}

func GetCameraByID(id int, db *sqlx.DB) (Camera, error) {
	var c Camera
	err := db.Get(&c, `
	select 
		a.id, 
		a.sensor_id,
		a.title,
		a.description, 
		a.description_long,
		a.hostname,
		a.project_name,
		a.storage_location,
		a.interval,
		a.next_capture_time,
		a.status,
		a.alert,
		a.active,
		a.created_at,
		a.updated_at
	from camera a
	where id=$1
		`, id)
	if err != nil {
		return c, fmt.Errorf("getCameraByID(); %v", err)
	}
	return c, nil
}

// SetNextCaptureTime sets the next capture time for a given camera based on interval,
// and replaces the value in the struct.
// Interval is in seconds.
func (c *Camera) SetNextCaptureTime() error {
	next := time.Now().Add(time.Duration(c.Interval) * time.Second)
	_, err := c.DB.Exec(`update camera set next_capture_time=$1 where id=$2`, next, c.ID)
	if err != nil {
		return fmt.Errorf("SetNextCapTureTime(): unable to set next_capture_time %w", err)
	}
	out := pgtype.Timestamp{
		Time:  next,
		Valid: true,
	}
	c.NextCapTureTime = out
	return nil
}

func ResetConnectivity(db *sqlx.DB) error {
	_, err := db.Exec(`update camera set next_capture_time = null;`)
	if err != nil {
		return fmt.Errorf("unable to reset next_capture_time: %w", err)
	}
	return nil
}

func (c *Camera) CaptureTimelapseImage(timelapse StreamTimelapse) {
	err := timelapse.CaptureTimelapseImage()
	if err != nil {
		log.Printf("[ERROR] Problems with capturing timelapse image: %v", err)
		c.Status = "error"
		return
	}
	log.Printf("[INFO] Timelapse image captured for %s (%s)", timelapse.ProjectName, timelapse.Hostname)
	c.lastframe = timelapse.data
	c.Status = fmt.Sprintf("image captured %d", time.Now().Unix())
	fmt.Println("status", c.Status)
}
