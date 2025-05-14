package department

import (
	"database/sql"
	"encoding/json"
	"errors"
	"labyrinth/logger"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (d DepartmentHandlers) GetDepartmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID, ok := ctx.Value("id").(uuid.UUID)
	if !ok || userID == uuid.Nil {
		logger.NewErrMessage("Invalid user ID in context",
			zap.String("operation", "GetDepartmentHandler"),
			zap.Any("context_values", ctx.Value("id")),
		)
		http.Error(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	userPathId, err := uuid.Parse(vars["user_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid path variable",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("variable", "user_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	if userPathId != userID {
		logger.NewWarnMessage("User ID mismatch",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("context_user_id", userID.String()),
			zap.String("path_user_id", userPathId.String()),
		)
		http.Error(w, "Forbidden: user ID mismatch", http.StatusForbidden)
		return
	}

	companyId, err := uuid.Parse(vars["company_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid company ID format",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("variable", "company_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid company ID format", http.StatusBadRequest)
		return
	}

	departmentId, err := uuid.Parse(vars["department_id"])
	if err != nil {
		logger.NewWarnMessage("Invalid department ID format",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("variable", "department_id"),
			zap.Error(err),
		)
		http.Error(w, "Invalid department ID format", http.StatusBadRequest)
		return
	}

	departmentData, err := bl.Department.GetDepartment(userID, companyId, departmentId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.NewWarnMessage("Department not found",
				zap.String("operation", "GetDepartmentHandler"),
				zap.String("department_id", departmentId.String()),
			)
			http.Error(w, "Department not found", http.StatusNotFound)
			return
		}

		logger.NewErrMessage("Failed to get department",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to get department data", http.StatusInternalServerError)
		return
	}

	fetchedDir, err := fsl.Folder.GetFolder(departmentId)
	if err != nil {
		logger.NewErrMessage("Failed to get department folder",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to get department folder", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":     "success",
		"department": departmentData,
		"folder":     fetchedDir,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.NewErrMessage("Failed to encode response",
			zap.String("operation", "GetDepartmentHandler"),
			zap.String("department_id", departmentId.String()),
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	logger.NewInfoMessage("Department data retrieved successfully",
		zap.String("user_id", userID.String()),
		zap.String("company_id", companyId.String()),
		zap.String("department_id", departmentId.String()),
	)
}
