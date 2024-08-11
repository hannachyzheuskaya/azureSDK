package virtualMachine

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"x/internal/app/clients/networkInterface"
	"x/internal/app/clients/networkSecurityGroup"
	"x/internal/app/clients/publicIP"
	"x/internal/app/clients/resourceGroup"
	"x/internal/app/clients/subnet"
	"x/internal/app/clients/virtualMachine"
	"x/internal/app/clients/virtualNetwork"
	"x/internal/app/controller/handlers/util"
	"x/internal/app/processor/response"
	processor "x/internal/app/processor/virtualMachine"
	"x/internal/lib/api"
	sl "x/internal/lib/logger"
	resp "x/internal/lib/response"
)

var (
	clientRG     *resourceGroup.ClientResourceGroup
	clientVNET   *virtualNetwork.ClientVirtualNetwork
	clientSUBNET *subnet.ClientSubnet
	clientIP     *publicIP.ClientPublicIP
	clientNSG    *networkSecurityGroup.ClientSecurityGroup
	clientNIC    *networkInterface.ClientNetworkInterface
	clientVM     *virtualMachine.ClientVirtualMachine
)

type Response struct {
	resp.Response
	response.ProcessedResponse `json:"azure_response"`
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.virtualMachine.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		cred, subscriptionId, err := util.Credential(r)
		if err != nil {
			log.Error(op, "failed to get credential", err)
			render.JSON(w, r, sl.Err(err))
			return
		}

		req := &api.RequestAPI{}
		err = render.DecodeJSON(r.Body, req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, resp.Error("empty request"))
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to decode request"))
			return
		}
		log.Info("request body decoded")

		clientRG, err = resourceGroup.New(subscriptionId, cred)
		clientVNET, err = virtualNetwork.New(subscriptionId, cred)
		clientSUBNET, err = subnet.New(subscriptionId, cred)
		clientIP, err = publicIP.New(subscriptionId, cred)
		clientNSG, err = networkSecurityGroup.New(subscriptionId, cred)
		clientNIC, err = networkInterface.New(subscriptionId, cred)
		clientVM, err = virtualMachine.New(subscriptionId, cred)
		if err != nil {
			log.Error(op, "failed to create virtual machine client", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to create virtual machine client"))
			return
		}

		myProcessor := processor.New(
			clientRG,
			clientVNET,
			clientSUBNET,
			clientIP,
			clientNSG,
			clientNIC,
			clientVM)

		resp := myProcessor.Process(req)
		responseOK(w, r, resp)

	}
}

func responseOK(w http.ResponseWriter, r *http.Request, pr response.ProcessedResponse) {
	render.JSON(w, r, Response{
		Response:          resp.OK(),
		ProcessedResponse: pr,
	})
}
