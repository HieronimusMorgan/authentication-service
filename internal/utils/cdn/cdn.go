package cdn

import (
	"authentication/internal/dto/out"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadImageToCDN(ipCdn string, file *multipart.FileHeader, clientID, authToken string) (out.ImageResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add clientID field
	if err := writer.WriteField("client_id", clientID); err != nil {
		return out.ImageResponse{}, err
	}

	part, err := writer.CreateFormFile("image", file.Filename) // 🔁 changed "images" → "image"
	if err != nil {
		return out.ImageResponse{}, err
	}

	src, err := file.Open()
	if err != nil {
		return out.ImageResponse{}, err
	}
	defer src.Close()

	if _, err = io.Copy(part, src); err != nil {
		return out.ImageResponse{}, err
	}

	if err := writer.Close(); err != nil {
		return out.ImageResponse{}, err
	}

	req, err := http.NewRequest("POST", ipCdn+"/v1/upload-photo-profile", body)
	if err != nil {
		return out.ImageResponse{}, err
	}
	req.Header.Set("Authorization", authToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return out.ImageResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return out.ImageResponse{}, fmt.Errorf("failed to upload image: %s", resp.Status)
	}

	var res struct {
		Data out.ImageResponse `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return out.ImageResponse{}, err
	}

	res.Data.ImageURL = ipCdn + "/v1" + res.Data.ImageURL

	return res.Data, nil
}
