package handler

// func successResponse(w http.ResponseWriter, payload interface{}) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	response := map[string]interface{}{
// 		"error":  false,
// 		"status": http.StatusOK,
// 		"data":   payload,
// 	}
// 	json.NewEncoder(w).Encode(response)
// }

// func errorResponse(w http.ResponseWriter, status int, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(status)
// 	response := map[string]interface{}{
// 		"status":  status,
// 		"message": message,
// 		"error":   true,
// 	}
// 	json.NewEncoder(w).Encode(response)
// }
