package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"sejastip.id/api/util"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"sejastip.id/api"
	"sejastip.id/api/handler"
)

type BankHandler struct {
	uc api.BankUsecase
}

func NewBankHandler(uc api.BankUsecase) BankHandler {
	return BankHandler{uc}
}

func (h *BankHandler) RegisterHandler(r *httprouter.Router) error {
	if r == nil {
		return errors.New("Router must not be nil")
	}

	r.POST("/banks", handler.Decorate(h.CreateBank, handler.AppAuth...))
	r.GET("/banks", handler.Decorate(h.GetBanks, handler.AppAuth...))

	return nil
}

func (h *BankHandler) CreateBank(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	decoder := json.NewDecoder(r.Body)
	var bankForm api.BankForm
	if err := decoder.Decode(&bankForm); err != nil {
		api.Error(w, err)
		return err
	}

	file, extension, err := util.DecodeUploadedBase64File(bankForm.ImageFile)
	if err != nil {
		err = api.SejastipError{
			Message:    fmt.Sprintf("Error parsing file: %v", err),
			ErrorCode:  400,
			HTTPStatus: http.StatusBadRequest,
		}
		api.Error(w, err)
		return err
	}

	ctx := r.Context()

	// upload file
	filename := fmt.Sprintf("%s%s", bankForm.Name, extension)
	filePath, err := h.uc.UploadBankImage(ctx, filename, file)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error in uploading file")
	}

	bank := api.Bank{
		Name:  bankForm.Name,
		Image: filePath,
	}
	err = h.uc.CreateBank(ctx, &bank)
	if err != nil {
		api.Error(w, err)
		return err
	}

	api.OK(w, bank, nil)
	return nil
}

func (h *BankHandler) GetBanks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	helper := api.NewQueryHelper(r)
	limit := helper.GetInt("limit", 10)
	offset := helper.GetInt("offset", 0)

	// Get all banks
	ctx := r.Context()
	banks, total, err := h.uc.GetBanks(ctx, limit, offset)
	if err != nil {
		api.Error(w, err)
		return errors.Wrap(err, "error getting banks")
	}

	meta := api.NewPaginationMeta(limit, offset, int(total))
	api.OK(w, banks, meta)
	return nil
}
