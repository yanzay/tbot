package tbot

import "fmt"

// FileURL returns file URL ready for download
func (c *Client) FileURL(file *File) string {
	return fmt.Sprintf("%s/file/bot%s/%s", c.baseURL, c.token, file.FilePath)
}
