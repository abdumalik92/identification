package contreller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/abdumalik92/identification/internal/models"
	"github.com/abdumalik92/identification/internal/service"
	"github.com/abdumalik92/identification/internal/utils"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SendOrderToCFT(c *gin.Context) {
	var claims *models.Claims

	seaFile := utils.AppSettings.FileService
	Login := "b4d41e2c2cd03021279933a4e1e8fa5e"

	auth := c.GetHeader("Authorization")
	claims = utils.GetStruct(auth, c.Writer)
	if c.GetHeader("Sing") != "" {
		claims.Phone = c.GetHeader("Phone")
		hash := utils.RSHA256([]byte(claims.Phone + Login + utils.AppSettings.SecretKey.Key))
		if hash != c.GetHeader("Sing") {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Неверные данные"})
			return
		}
	}

	form, _ := c.MultipartForm()
	user := form.Value["user"][0]
	repoID := ""
	product := ""
	if user == "" {
		product = "HUMO_ONLINE"
		repoID = seaFile.HumoOnline
	} else {
		product = "MEGAFON_LIFE"
		repoID = seaFile.MegafonLife
	}
	var files []string
	t := time.Now().Format("2006_01_02_15_04")
	if len(form.File) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Недостаточное количество фотографий"})
		return
	}

	if len(form.File) > 4 {
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Вы отправили больше фотографий, чем положено (макс 4 фотографий)"})
		return
	}

	for k := 1; k <= len(form.File); k++ {
		err := c.SaveUploadedFile(form.File["files"+fmt.Sprint(k)][0], "./temp/"+t+"_"+claims.Phone+fmt.Sprint(k)+form.File["files"+fmt.Sprint(k)][0].Filename)
		if err != nil {
			log.Println("Error save ", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "Не удалось загрузить файлы..."})
			return
		}
		files = append(files, "./temp/"+t+"_"+claims.Phone+fmt.Sprint(k)+form.File["files"+fmt.Sprint(k)][0].Filename)
	}

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Authorization
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	response := models.MainModel{}

	urlAuth := seaFile.BaseURL + "/api2/auth-token/"
	methodAuth := "POST"

	payloadAuth := strings.NewReader("{\"username\":\"" + seaFile.Username + "\",\"password\":\"" + seaFile.Password + "\"}")

	clientAuth := &http.Client{}
	reqAuth, err := http.NewRequest(methodAuth, urlAuth, payloadAuth)

	if err != nil {
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't authReq Delete")
					c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
					return
				}
			}
		}
		log.Println("Can't authReq unmarshal json")
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
		return
	}
	reqAuth.Header.Add("Content-Type", "application/json")

	resAuth, err := clientAuth.Do(reqAuth)
	//defer resAuth.Body.Close()
	body, err := ioutil.ReadAll(resAuth.Body)
	if err != nil {
		log.Println("Can't authReq send req ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"reason": "Что-то пошло не так..."})
		return
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't authReq Delete")
					c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
					return
				}
			}
		}
		log.Println("Can't authReq unmarshal json")
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
		return
	}
	log.Println("Auth success")
	//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Create Directory
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	urlCreateDir := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
	methodCreateDir := "POST"

	payloadCreateDir := strings.NewReader("operation=mkdir")

	clientCreateDir := &http.Client{}
	reqCreateDir, err := http.NewRequest(methodCreateDir, urlCreateDir, payloadCreateDir)

	if err != nil {
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't func createRepository Delete File")
				}
			}
		}
		log.Println("Can't createRepository err ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
		return
	}
	reqCreateDir.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
	reqCreateDir.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resCreateDir, err := clientCreateDir.Do(reqCreateDir)
	if resCreateDir.StatusCode != 201 {
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't func createRepository Delete File")
				}
			}
		}
		log.Println("Can't createRepository err ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
		return
	}
	log.Println("Repository created")
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//Generate upload link
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	urlGenLink := seaFile.BaseURL + "/api2/repos/" + repoID + "/upload-link/?p=/" + t + "_" + claims.Phone
	methodGenLink := "GET"

	clientGenLink := &http.Client{}
	reqGenLink, err := http.NewRequest(methodGenLink, urlGenLink, nil)

	if err != nil {
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't func createRepository Delete File")
				}
			}
		}

		urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
		methodDelete := "DELETE"

		clientDelete := &http.Client{}
		reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

		if err != nil {
			log.Println("Can't deleteRepository")
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
		reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)

		resDelete, err := clientDelete.Do(reqDelete)
		//defer resDelete.Body.Close()
		_, err = ioutil.ReadAll(resDelete.Body)
		if err != nil {
			log.Println("Can't deleteRepository")
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
	}
	reqGenLink.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
	resGenLink, err := clientGenLink.Do(reqGenLink)
	//defer resGenLink.Body.Close()
	bodyGenLink, err := ioutil.ReadAll(resGenLink.Body)
	if err != nil {
		log.Println()
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't func createRepository Delete File")
				}
			}
		}
		urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
		methodDelete := "DELETE"

		clientDelete := &http.Client{}
		reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

		if err != nil {
			log.Println("Can't deleteRepository err ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
		reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
		resDelete, err := clientDelete.Do(reqDelete)
		//defer resDelete.Body.Close()
		_, err = ioutil.ReadAll(resDelete.Body)
		if err != nil {
			log.Println("Can't deleteRepository err ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
	}
	log.Println("UploadLink created")

	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//UploadFile
	///////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	for i := 0; i < len(files); i++ {
		urlUploadFile := strings.ReplaceAll(string(bodyGenLink), "\"", "")

		method := "POST"

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("parent_dir", "/"+t+"_"+claims.Phone+"/")
		file, errFile2 := os.Open(files[i])

		part2, errFile2 := writer.CreateFormFile("file", filepath.Base(files[i]))
		_, errFile2 = io.Copy(part2, file)
		if errFile2 != nil {
			_ = file.Close()
			for _, filename := range files {
				if utils.FileExists(filename) {
					if err = utils.DeleteFile(filename); err != nil {
						log.Println("Can't func UploadFile Delete File")
					}
				}
			}
			urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
			methodDelete := "DELETE"
			clientDelete := &http.Client{}
			reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

			if err != nil {
				log.Println("Can't deleteRepository")
				c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
				return
			}
			reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
			resDelete, err := clientDelete.Do(reqDelete)
			//defer resDelete.Body.Close()
			_, err = ioutil.ReadAll(resDelete.Body)
			log.Println("Can't UploadFile ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}

		_ = writer.Close()

		client := &http.Client{
			Timeout: 90 * time.Second,
		}
		req, err := http.NewRequest(method, urlUploadFile, payload)

		if err != nil {
			_ = file.Close()
			for _, filename := range files {
				if utils.FileExists(filename) {
					if err = utils.DeleteFile(filename); err != nil {
						log.Println("Can't func UploadFile Delete File")
					}
				}
			}
			urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
			methodDelete := "DELETE"
			clientDelete := &http.Client{}
			reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

			if err != nil {
				log.Println("Can't deleteRepository")
				c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
				return
			}
			reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
			resDelete, err := clientDelete.Do(reqDelete)
			//defer resDelete.Body.Close()
			_, err = ioutil.ReadAll(resDelete.Body)

			log.Println("Can't UploadFile ", err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
		req.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)

		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {
			log.Println("Upload file error ", err.Error())
			_ = file.Close()
			for _, filename := range files {
				if utils.FileExists(filename) {
					if err = utils.DeleteFile(filename); err != nil {
						log.Println("Can't func UploadFile Delete File")
					}
				}
			}
			urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
			methodDelete := "DELETE"
			clientDelete := &http.Client{}
			reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

			if err != nil {
				log.Println("Can't deleteRepository")
				c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
				return
			}
			reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
			resDelete, err := clientDelete.Do(reqDelete)
			//defer resDelete.Body.Close()
			_, _ = ioutil.ReadAll(resDelete.Body)
			log.Println("Can't UploadFile")
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
		//defer res.Body.Close()
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println("Upload file error ", err.Error())
			_ = file.Close()
			for _, filename := range files {
				if utils.FileExists(filename) {
					if err = utils.DeleteFile(filename); err != nil {
						log.Println("Can't func UploadFile Delete File")
					}
				}
			}
			urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
			methodDelete := "DELETE"
			clientDelete := &http.Client{}
			reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

			if err != nil {
				log.Println("Can't deleteRepository")
				c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
				return
			}
			reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
			resDelete, err := clientDelete.Do(reqDelete)
			//defer resDelete.Body.Close()
			_, _ = ioutil.ReadAll(resDelete.Body)
			log.Println("Can't UploadFile")
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}

		_ = file.Close()
	}

	for _, filename := range files {
		if utils.FileExists(filename) {
			if err = utils.DeleteFile(filename); err != nil {
				log.Println("Can't func UploadFile Delete File ", err.Error())
			}
		}
	}

	if err := service.SendOrderToCFT(product, claims.Phone, seaFile.BaseURL+"/#my-libs/lib/"+repoID+"/"+t+"_"+claims.Phone); err != nil {
		for _, filename := range files {
			if utils.FileExists(filename) {
				if err = utils.DeleteFile(filename); err != nil {
					log.Println("Can't ЦФТ Delete file")
				}
			}
		}
		urlDelete := seaFile.BaseURL + "/api2/repos/" + repoID + "/dir/?p=/" + t + "_" + claims.Phone
		methodDelete := "DELETE"
		fmt.Println("Part 1 2")
		clientDelete := &http.Client{}
		reqDelete, err := http.NewRequest(methodDelete, urlDelete, nil)

		if err != nil {
			log.Println("Can't deleteRepository")
			c.JSON(http.StatusBadRequest, gin.H{"reason": "Что-то пошло не так..."})
			return
		}
		reqDelete.Header.Add("Authorization", seaFile.TokenHeadName+response.Token)
		resDelete, err := clientDelete.Do(reqDelete)
		defer resDelete.Body.Close()
		_, err = ioutil.ReadAll(resDelete.Body)
		if err != nil {
			log.Println("Неудалось передать информацию в ЦФТ ", err)
			c.JSON(400, gin.H{"reason": "Что-то пошло не так..."})
			return
		}

	}
	//return response
	c.JSON(200, gin.H{"reason": "Ваша заявка принято"})
	return
}
