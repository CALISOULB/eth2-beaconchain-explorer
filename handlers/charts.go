package handlers

import (
	"eth2-exporter/services"
	"eth2-exporter/types"
	"eth2-exporter/utils"
	"eth2-exporter/version"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var chartsTemplate = template.Must(template.New("charts").Funcs(utils.GetTemplateFuncs()).ParseFiles("templates/layout.html", "templates/charts.html"))
var genericChartTemplate = template.Must(template.New("chart").Funcs(utils.GetTemplateFuncs()).ParseFiles("templates/layout.html", "templates/genericchart.html"))

// Charts uses a go template for presenting the page to show charts
func Charts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	chartsPageData := services.LatestChartsPageData()
	if chartsPageData == nil {
		http.Error(w, "The requested data is currently unavailable, please retry in a few seconds", http.StatusServiceUnavailable)
		return
	}

	data := &types.PageData{
		Meta: &types.Meta{
			Title:       fmt.Sprintf("%v - Charts - beaconcha.in - %v", utils.Config.Frontend.SiteName, time.Now().Year()),
			Description: "beaconcha.in makes the Ethereum 2.0. beacon chain accessible to non-technical end users",
			Path:        "/charts",
		},
		ShowSyncingMessage:    services.IsSyncing(),
		Active:                "charts",
		Data:                  chartsPageData,
		Version:               version.Version,
		ChainSlotsPerEpoch:    utils.Config.Chain.SlotsPerEpoch,
		ChainSecondsPerSlot:   utils.Config.Chain.SecondsPerSlot,
		ChainGenesisTimestamp: utils.Config.Chain.GenesisTimestamp,
	}

	err := chartsTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		logger.Errorf("error executing template for %v route: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// GenericChart uses a go template for presenting the page of a generic chart
func GenericChart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	chartsPageData := services.LatestChartsPageData()
	if chartsPageData == nil {
		http.Error(w, "The requested data is currently unavailable, please retry in a few seconds", http.StatusServiceUnavailable)
		return
	}

	vars := mux.Vars(r)
	chartVar := vars["chart"]
	var chartData *types.GenericChartData
	for _, d := range *chartsPageData {
		if d.Path == chartVar {
			chartData = d.Data
			break
		}
	}

	if chartData == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	data := &types.PageData{
		Meta: &types.Meta{
			Title:       fmt.Sprintf("%v - %v Chart - beaconcha.in - %v", chartData.Title, utils.Config.Frontend.SiteName, time.Now().Year()),
			Description: "beaconcha.in makes the Ethereum 2.0. beacon chain accessible to non-technical end users",
			Path:        "/charts/" + chartVar,
		},
		ShowSyncingMessage:    services.IsSyncing(),
		Active:                "charts",
		Data:                  chartData,
		Version:               version.Version,
		ChainSlotsPerEpoch:    utils.Config.Chain.SlotsPerEpoch,
		ChainSecondsPerSlot:   utils.Config.Chain.SecondsPerSlot,
		ChainGenesisTimestamp: utils.Config.Chain.GenesisTimestamp,
	}

	err := genericChartTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		logger.Errorf("error executing template for %v route: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
