package config

import (
	"github.com/gofiber/fiber/v2"
)

func CORS(whiteList []string) fiber.Handler {
	// สร้าง map เพื่อเก็บรายการที่อนุญาตให้ใช้งาน
	whiteListMap := make(map[string]bool)
	for _, origin := range whiteList {
		whiteListMap[origin] = true
	}

	// ส่วน middleware ที่จะทำการตรวจสอบและกำหนดค่า Access-Control-Allow-Origin และอื่น ๆ
	return func(c *fiber.Ctx) error {
		// ตรวจสอบว่าเป็นการเรียก preflight request หรือไม่
		if c.Method() == "OPTIONS" {
			c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Set("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type, Accept")
			c.Set("Access-Control-Max-Age", "86400") // 24 ชั่วโมง
			return c.SendStatus(fiber.StatusNoContent)
		}

		origin := c.Get("Origin")
		// ตรวจสอบว่า origin อยู่ใน whitelist หรือไม่ หากไม่อยู่ให้ return ข้อความผิดพลาด
		if !whiteListMap[origin] && !whiteListMap["*"] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Not allowed by CORS",
			})
		}

		// กำหนดค่า Access-Control-Allow-Origin
		c.Set("Access-Control-Allow-Origin", origin)
		// กำหนดค่า Access-Control-Allow-Credentials
		c.Set("Access-Control-Allow-Credentials", "true")

		return c.Next()
	}
}
