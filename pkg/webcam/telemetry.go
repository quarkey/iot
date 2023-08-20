package webcam

import "net/http"

// Online checks if the camera is online and returning true by checking if status code is 200.
func (c *Camera) Online() bool {
	response, err := http.Get(c.Hostname)
	if err != nil {
		return false
	}
	defer response.Body.Close()

	return response.StatusCode == 200
}
