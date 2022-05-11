package User

import (
	"PemWeb_BE/Auth"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"os"
	"time"
)

func Routes(db *gorm.DB, q *gin.Engine) {
	r := q.Group("/user")
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
			"error":nil,
			"data": gin.H{
				"id":        user.ID,
				"nama":      user.Name,
				"nik":       user.NIK,
				"telp":      user.Telepon,
				"kota_asal": user.KotaAsal,
			},
		})
	})
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
		if err := db.Where("telp=?", input.Telepon).Take(&login); err.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Telephone number does not exist",
				"error":   err.Error.Error(),
			})
			return
		}
		if login.Password == hash(input.Password) {
			token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
				"id":  login.ID,
				"exp": time.Now().Add(time.Hour * 7 * 24).Unix(),
			})
			godotenv.Load(".env")
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
				"mess"
			})
		}
	})
}

func hash(input string) string {
	hash := []byte(input)
	hashed, _ := bcrypt.GenerateFromPassword(hash, bcrypt.DefaultCost)
	return string(hashed)
}
