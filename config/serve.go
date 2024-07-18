package config

import (
	// "fmt"
	// "net/http"
	// "os"
	// "path/filepath"

	"github.com/gin-gonic/gin"
)

func ServeImage(c *gin.Context) {
	// imageId := c.Param("imageId")

	// Get image details from the database
	// var filename string
	// query := `SELECT filename FROM images WHERE id = ?`
	// if err := db.QueryRow(query, imageId).Scan(&filename); err != nil {
	// 	if err == sql.ErrNoRows {
	// 		c.String(http.StatusNotFound, "Image not found")
	// 	} else {
	// 		c.String(http.StatusInternalServerError, "Database query failed: %s", err.Error())
	// 	}
	// 	return
	// }

	// // Determine the local file path using the imageId as the folder name
	// imageFolder := filepath.Join(localImageRootFolder, fmt.Sprintf("image%s", imageId))
	// localFilePath := filepath.Join(imageFolder, filename)

	// // Check if the file exists locally
	// if _, err := os.Stat(localFilePath); os.IsNotExist(err) {
	// 	// File does not exist, download it
	// 	remoteFilePath := remoteImageURL + filename
	// 	if err := downloadFile(localFilePath, remoteFilePath); err != nil {
	// 		c.String(http.StatusInternalServerError, "Failed to download image: %s", err.Error())
	// 		return
	// 	}
	// }

	// Serve the image
	// c.File(localFilePath)
}
