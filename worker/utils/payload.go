package utils

// GetIntPointer: ดึงค่า int แบบ Pointer (รองรับทั้ง int, float64 จาก JSON และ nil)
func GetIntPointer(data map[string]interface{}, key string) *int {
	val, ok := data[key]
	if !ok || val == nil {
		return nil
	}

	if f, ok := val.(float64); ok {
		i := int(f)
		return &i
	}

	if i, ok := val.(int); ok {
		return &i
	}

	return nil
}

// GetString: ดึงค่า string (ถ้าไม่มีคืนค่าว่าง "")
func GetString(data map[string]interface{}, key string) string {
	val, ok := data[key]
	if !ok || val == nil {
		return ""
	}
	if s, ok := val.(string); ok {
		return s
	}
	return ""
}

// GetStringArray: ดึงค่า Array String อย่างปลอดภัย (รองรับทั้ง []string และ []interface{})
func GetStringArray(data map[string]interface{}, key string) []string {
	val, ok := data[key]
	if !ok || val == nil {
		return []string{} // คืนค่า Array ว่าง "{}" เพื่อไม่ให้ error malformed array
	}

	// กรณีเป็น []string อยู่แล้ว
	if arr, ok := val.([]string); ok {
		return arr
	}

	// กรณีเป็น []interface{} (มักเจอตอน decode JSON array)
	if arr, ok := val.([]interface{}); ok {
		strArr := make([]string, len(arr))
		for i, v := range arr {
			if s, ok := v.(string); ok {
				strArr[i] = s
			}
		}
		return strArr
	}

	if s, ok := val.(string); ok && s != "" {
		return []string{s}
	}

	return []string{}
}