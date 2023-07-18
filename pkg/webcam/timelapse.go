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

type StreamTimelapse struct {
	StorageDir  string
	Hostname    string
	ProjectPath string
	ProjectName string
	JsonConfig  string
	projectfh   *os.File
}

func NewTimelase(storageDir string, hostname string, projectName string, outputName string) (StreamTimelapse, error) {
	projectPath := path.Join(storageDir, projectName)
	cfgPath := path.Join(projectPath, fmt.Sprintf("%s.json", projectName))
	stl := StreamTimelapse{
		StorageDir:  storageDir,
		Hostname:    hostname,
		ProjectPath: projectPath,
		ProjectName: projectName,
		JsonConfig:  cfgPath,
	}
	if stl.projectExist() {
		// fmt.Println("opening project", stl.ProjectPath)
		err := stl.openProject()
		if err != nil {
			return StreamTimelapse{}, err
		}
	} else {
		// fmt.Println("unable to find project creating", stl.ProjectPath)
		stl.createAndOpenProject()
	}
	return stl, nil
}

// openProject opens the project file handler
func (s *StreamTimelapse) openProject() error {
	fh, err := os.Open(s.JsonConfig)
	if err != nil {
		return fmt.Errorf("unable to open project file: %v", err)
	}
	s.projectfh = fh
	return nil
}

// close closes the project file handler
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

// CaptureTimelapseImage captures a single image from the stream and saves it to the project folder.
// The image is saved as a jpeg file with the name: <unix-timestamp>-image.jpg
func (s *StreamTimelapse) CaptureTimelapseImage() error {
	resp, err := http.Get(s.Hostname)
	if err != nil {
		return fmt.Errorf("unable to open stream: %v", err)
	}
	var output bytes.Buffer
	firstRead := true
	clen := 0
	for {
		buffer := make([]byte, 1024)
		n, _ := resp.Body.Read(buffer)
		// fmt.Printf("reading %d bytes...\n", n)
		if firstRead {
			clen = contentLenght(buffer)
			// fmt.Println("content-lenght:", clen)
			firstRead = false
		}
		output.Write(buffer[:n])
		// reading until content-lenght and some extra just for good measure
		if output.Len() > clen {
			break
		}
	}
	resp.Body.Close()
	imgdata, err := extractImage(output.Bytes())
	if err != nil {
		return fmt.Errorf("unable to extract image: %v", err)
	}
	// fmt.Println("read total of", output.Len(), "bytes..")
	// fmt.Println("image lenght:", len(imgdata))
	s.closeProjectfh() // close file handler
	err = saveImg(imgdata, s.ProjectPath)
	if err != nil {
		return err
	}
	return nil
}

// projectExist checks if the project folder exists
func (s *StreamTimelapse) projectExist() bool {
	if _, err := os.Stat(s.ProjectPath); os.IsNotExist(err) {
		return false
	}
	return true
}

var ErrNoImage = fmt.Errorf("image/jpeg end not detected")

// extractImage extracts the image from the stream
func extractImage(data []byte) ([]byte, error) {
	n := len(data)
	start := 0
	end := 0
	found_end := false
	for i := 0; i < n-1; i++ {
		if data[i] == 0xFF && data[i+1] == 0xD8 {
			// fmt.Printf("Image/jpeg start detected at pos %d\n", i)
			start = i
		}
		if data[i] == 0xFF && data[i+1] == 0xD9 {
			// fmt.Printf("Image/jpeg end detected at pos %d\n", i)
			end = i
			found_end = true
			break
		}
	}
	if !found_end {
		return nil, ErrNoImage
	}
	return data[start : end+2], nil
}

// contentLenght extracts the content-lenght from the stream
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

// saveImg saves the image to the project folder
func saveImg(data []byte, outpath string) error {
	// Save the image data to a file
	outputName := fmt.Sprintf("%d-image.jpg", time.Now().Unix())
	file, err := os.Create(path.Join(outpath, outputName))
	if err != nil {
		// Handle error
		return fmt.Errorf("unable to create image file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		// Handle error
		return fmt.Errorf("unable to write image file: %v", err)
	}

	// fmt.Println("Image saved to image.jpg.")
	return nil
}
