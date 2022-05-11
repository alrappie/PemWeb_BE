package User

import (
	"PemWeb_BE/Auth"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/user")
	//halaman profil pengguna yang udah login
	r.GET("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong in Query",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful",
			"error":   nil,
			"data": gin.H{
				"id":        user.ID,
				"nama":      user.Name,
				"nik":       user.NIK,
				"telp":      user.Telepon,
				"kota_asal": user.KotaAsal,
			},
		})
	})
	//registrasi pengguna baru
	r.POST("/register", func(c *gin.Context) {
		var input Register
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Value sent was incorrect",
				"error":   err.Error(),
			})
			return
		}
		create := User{
			NIK:      input.NIK,
			Name:     input.Name,
			Password: hash(input.Password),
			Telepon:  input.Telepon,
			KotaAsal: input.KotaAsal,
		}
		if err := db.Create(&create); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong with the database",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Akun berhasil dibuat",
			"error":   nil,
			"data": gin.H{
				"id":        create.ID,
				"nama":      input.Name,
				"kota_asal": input.KotaAsal,
				"telp":      input.Telepon,
				"nik":       input.NIK,
			},
		})
	})
	//login pengguna lama
	r.POST("/login", func(c *gin.Context) {
		var input Login
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please check login info",
				"error":   err.Error(),
			})
			return
		}
		login := User{}
		if err := db.Where("telepon=?", input.Telepon).Take(&login); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Telephone number does not exist",
				"error":   err.Error.Error(),
			})
			return
		}
		fmt.Println(login.Password)
		fmt.Println(hash(input.Password))
		if err := bcrypt.CompareHashAndPassword([]byte(login.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Hayooo lupa password XD",
				"error":   "Wrong password",
			})
			return
		} else {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})
			err := godotenv.Load(".env")
			if err != nil {
				log.Fatal(err.Error())
			}
			strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_GIT")))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Token failed to be created, please contact system administrator",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login berhasil",
				"error":   nil,
				"data": gin.H{
					"telp":  login.Telepon,
					"token": strToken,
				},
			})
			return
		}
	})
	//login dan registrasi pake google oauth2
	r.GET("/google", Auth.GInit)
	r.GET("/google/callback", func(c *gin.Context) {
		var a = Auth.GCallback(c)
		var b = []byte(a)
		var goog Google
		user := User{}
		if err := json.Unmarshal(b, &goog); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to unmarshal Google data",
				"err":     err.Error(),
			})
			return
		}
		if err := db.Where("telepon=?", goog.Email).Take(&user); err.Error != nil {
			user = User{
				Name:     goog.Name,
				Telepon:  goog.Email,
				Password: hash(goog.Sub),
			}
			if err := db.Create(&user); err.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Something went wrong on the server side. please contact system administrator",
					"error":   err.Error.Error(),
				})
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": true,
				"message": "no account was found, creating new account",
				"error":   nil,
			})
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
			"id":  user.ID,
			"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
		})
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal(err.Error())
		}
		strToken, err := token.SignedString([]byte(os.Getenv("TOKEN_GIT")))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Token failed to be created, please contact system administrator",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login berhasil",
			"error":   nil,
			"data": gin.H{
				"telp":  user.Telepon,
				"token": strToken,
			},
		})
	})
	//hapus akun
	r.DELETE("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Deletion has failed, please contact system administrator",
				"error":   err.Error.Error(),
			})
			return
		}
		if err := db.Delete(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Account has failed to be deleted",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Account deleted successfully",
			"error":   nil,
			"data": gin.H{
				"nama":      user.Name,
				"nik":       user.NIK,
				"kota_asal": user.KotaAsal,
				"telp":      user.Telepon,
			},
		})
	})
	r.PATCH("/", Auth.Authorization(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var input User
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please fix your input",
				"error":   err.Error(),
			})
			return
		}
		user := User{}
		if err := db.Where("id=?", id).Take(&user); err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Account was not found",
				"error":   err.Error.Error(),
			})
			return
		}
		user = User{
			ID:       user.ID,
			Name:     input.Name,
			NIK:      input.NIK,
			Telepon:  input.Telepon,
			KotaAsal: input.KotaAsal,
			Password: hash(input.Password),
		}
		err := db.Model(&user).Updates(user)
		if err.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong on account update",
				"error":   err.Error.Error(),
			})
			return
		}
		if err.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "No data has been changed",
				"error":   err.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Data updated successfully",
			"error":   nil,
			"data": gin.H{
				"id":        user.ID,
				"nama":      user.Name,
				"nik":       user.NIK,
				"telp":      user.Telepon,
				"kota_asal": user.KotaAsal,
			},
		})
	})
}

func hash(input string) string {
	hash := []byte(input)
	hashed, _ := bcrypt.GenerateFromPassword(hash, bcrypt.DefaultCost)
	return string(hashed)
}
