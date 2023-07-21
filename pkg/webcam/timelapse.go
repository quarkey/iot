package webcam

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"time"
)

// StreamTimelapse is a struct that holds the configuration for a timelapse project.
type StreamTimelapse struct {
	StorageDir  string
	Hostname    string
	ProjectPath string
	ProjectName string
	JsonConfig  string
	projectfh   *os.File
	retryCount  int
}

// NewTimelapse either opens or creates a new timelapse project.
func NewTimelase(storageDir string, hostname string, projectName string) (StreamTimelapse, error) {
	projectPath := path.Join(storageDir, projectName)
	cfgPath := path.Join(projectPath, fmt.Sprintf("%s.json", projectName))

	stl := StreamTimelapse{
		StorageDir:  storageDir,
		Hostname:    hostname,
		ProjectPath: projectPath,
		ProjectName: projectName,
		JsonConfig:  cfgPath,
		retryCount:  5,
	}
	if stl.projectExist() {
		err := stl.openProject()
		if err != nil {
			return StreamTimelapse{}, err
		}
	} else {
		stl.createAndOpenProject()
	}
	return stl, nil
}

// CaptureTimelapseImage captures a single image from stream and saves it to the project folder.
// The image is saved as a jpeg file with the name: <unix-timestamp>-image.jpg to prevent overwriting.

func (s *StreamTimelapse) CaptureTimelapseImage() error {
	// we are not using fh at this moment in time, so closing right away
	s.closeProjectfh() // close file handler

	// TODO: add retry logic

	output, err := s.readStreamResp()
	if err != nil {
		return fmt.Errorf("unable to read stream: %v", err)
	}

	imgdata, err := extractImage(output.Bytes())
	if err != nil {
		if err == ErrNoEndOfImage {
			fmt.Println("no end of image detected")
		}
		return fmt.Errorf("unable to extract image: %v", err)
	}

	err = saveImg(imgdata, s.ProjectPath)
	if err != nil {
		return fmt.Errorf("unable to save image: %v", err)
	}

	return nil
}

// openProject opens the project file handler.
func (s *StreamTimelapse) openProject() error {
	fh, err := os.Open(s.JsonConfig)
	if err != nil {
		return fmt.Errorf("unable to open project file: %v", err)
	}
	s.projectfh = fh
	return nil
}

// closeProjectfh closes the project file handler.
func (s *StreamTimelapse) closeProjectfh() {
	s.projectfh.Close()
}

// createAndOpenProject creates a project folder and a json configuration file
func (s *StreamTimelapse) createAndOpenProject() error {
	if err := os.Mkdir(s.ProjectPath, os.ModePerm); err != nil {
		log.Print(err)
	}
	file, err := os.Create(s.JsonConfig)
	if err != nil {
		return fmt.Errorf("unable to create project file: %v", err)
	}
	defer file.Close()
	// write config
	b, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("unable to marshal project config: %v", err)
	}
	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("unable to write project file: %v", err)
	}
	err = s.openProject()
	if err != nil {
		return err
	}
	return nil
}

// projectExist checks if the project folder exists.
func (s *StreamTimelapse) projectExist() bool {
	if _, err := os.Stat(s.ProjectPath); os.IsNotExist(err) {
		return false
	}
	return true
}

// readStreamResp reads the stream response and returns a buffer.
func (s *StreamTimelapse) readStreamResp() (bytes.Buffer, error) {
	var output bytes.Buffer
	firstRead := true
	clen := 0
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := client.Get(s.Hostname)
	if err != nil {
		return output, fmt.Errorf("unable to open stream: %v", err)
	}
	defer resp.Body.Close()
	for {
		buffer := make([]byte, 1024)
		n, _ := resp.Body.Read(buffer)
		if firstRead {
			clen = contentLenght(buffer)
			firstRead = false
		}
		output.Write(buffer[:n])
		// reading until content-lenght is reached
		if output.Len() > clen {
			break
		}
	}
	return output, nil
}

var (
	ErrNoStartOfImage = fmt.Errorf("image/jpeg start mark not detected")
	ErrNoEndOfImage   = fmt.Errorf("image/jpeg end mark not detected")
)

// extractImage scans the data to find the beginning and end markers of an image in the JPEG format.
// It looks for markers 0xFF 0xD8 (start marker) and 0xFF 0xD9 (end marker) to determine the image's boundaries.
func extractImage(data []byte) ([]byte, error) {
	var (
		n           = len(data)
		start       = 0
		end         = 0
		found_end   = false
		found_start = false
	)
	for i := 0; i < n-1; i++ {
		if data[i] == 0xFF && data[i+1] == 0xD8 {
			// fmt.Printf("Image/jpeg start detected at pos %d\n", i)
			start = i
			found_start = true
		}
		if data[i] == 0xFF && data[i+1] == 0xD9 {
			// fmt.Printf("Image/jpeg end detected at pos %d\n", i)
			end = i
			found_end = true
			break
		}
	}
	if !found_start {
		return nil, ErrNoStartOfImage
	}
	if !found_end {
		return nil, ErrNoEndOfImage
	}
	return data[start : end+2], nil
}

// contentLenght scans bytes for the content-lenght header.
func contentLenght(b []byte) int {
	re := regexp.MustCompile(`Content-Length:\s+(\d+)`)
	match := re.FindStringSubmatch(string(b))
	if len(match) >= 2 {
		numberStr := match[1]
		number, err := strconv.Atoi(numberStr)
		if err == nil {
			return number
		} else {
			return 0
		}
	} else {
		return 0
	}
}

// saveImg saves the image to the project folder.
func saveImg(data []byte, outpath string) error {
	// Save the image data to a file
	outputName := fmt.Sprintf("%d-image.jpg", time.Now().Unix())
	file, err := os.Create(path.Join(outpath, outputName))
	if err != nil {
		return fmt.Errorf("unable to create image file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("unable to write image file: %v", err)
	}
	return nil
}
