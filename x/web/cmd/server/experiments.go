package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/klokare/evo"
	"github.com/klokare/evo/x/web"
)

var (
	experimentsTemplate = `
{{define "content"}}
<div class="page-header">
	<div class="row">
		<div class="col-lg-12">
			<h1 id="buttons">{{.Study.Description}}</h1>
		</div>
	</div>
</div>
<div class="row">
	<div class="col-lg-12">
		<ul class="breadcrumb">
			<li><a href="/studies">Home</a></li>
			<li class="active">{{.Study.Description}}</li>
		</ul>
	</div>
</div>
{{template "charts" .}}
<div class="row">
	<div class="col-lg-12">
		<table class="table table-striped table-hover ">
			<thead>
				<tr>
					<th>ID</th>
					<th>Experiment</th>
					<th>Trials</th>
					<th>Last Update</th>
					<th>Diversity</th>
					<th>Best.Fitness</th>
					<th>Best.Complexity</th>
				</tr>
			</thead>
			<tbody>
				{{range .Experiments}}
				<tr>
					<td>{{.ID}}</td>
					<td><a href="#" onclick="setModal({{.ID}}, '{{.Description}}')" data-toggle="modal" data-target="#edit-modal">{{.Description}}</a></td>
					<td><a href="/studies/{{$.Study.ID}}/experiments/{{.ID}}">{{.Trials}}</a></td>
					<td>{{.LastUpdated}}</td>
					<td>{{.Diversity}}</td>
					<td>{{.Fitness}}</td>
					<td>{{.Complexity}}</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</div>
</div>
<div class="modal" id="edit-modal">
	<div class="modal-dialog">
		<div class="modal-content">
			<form class="form-horizontal" method="post" action="/experiments/edit">
				<fieldset>
					<div class="modal-header">
						<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
						<h4 class="modal-title">Edit Study</h4>
					</div>
					<div class="modal-body">
						<input type="hidden" name="study-id" value="{{.Study.ID}}" />
						<input type="hidden" name="experiment-id" id="experiment-id" />
						<div class="form-group">
							<label for="experiment-desc" class="col-lg-2 control-label">Description</label>
							<div class="col-lg-10">
								<input type="text" class="form-control" id="experiment-desc" name="experiment-desc">
							</div>
						</div>
						<div class="form-group">
							<label for="select" class="col-lg-2 control-label">Move to study</label>
							<div class="col-lg-10">
								<select class="form-control" name="experiment-sid" id="experiment-sid">
									{{range .Studies}}
									<option value="{{.ID}}" {{if eq .ID $.Study.ID}} selected {{end}}>{{.Description}}</option>
									{{end}}
								</select>
							</div>
						</div>						
						<div class="form-group">
							<label for="experiment-del" class="col-lg-2 control-label">Delete</label>
							<div class="col-lg-10">
								<div class="checkbox">
									<label>
									<input type="checkbox" id="experiment-del" name="experiment-del"> Select to delete this experiment
									</label>
								</div>
							</div>
						</div>
					</div>
					<div class="modal-footer">
						<button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
						<button type="submit" class="btn btn-primary">Save changes</button>
					</div>
		</div>
		</fieldset>
		</form>
	</div>
</div>
<script>
	function setModal(eid, desc) {
		document.getElementById("experiment-id").value = eid;
		document.getElementById("experiment-desc").value = desc;
	}
</script>		

{{end}}
	`
)

func getExperiments(w http.ResponseWriter, r *http.Request) {

	var err error
	var sid int
	vars := mux.Vars(r)
	if sid, err = strconv.Atoi(vars["sid"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	type study struct {
		ID          int
		Description string
	}
	type experiment struct {
		ID          int
		Description string
		LastUpdated time.Time
		Trials      int
		Fitness     float64
		Complexity  float64
		Diversity   float64
	}

	type chart struct {
		Box struct {
			Title      string
			XTitle     string
			Data       [][]float64
			XLabels    []string
			YMin, YMax float64
		}
		Scatter struct {
			Title      string
			XTitle     string
			YTitle     string
			Data       [][]float64
			XMin, XMax float64
			YMin, YMax float64
		}
	}

	var data struct {
		Chart       chart
		Study       study
		Studies     []study
		Experiments []experiment
	}

	var s web.Study
	if s, err = client.GetStudy(sid); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data.Study = study{s.ID, s.Description}
	data.Experiments = make([]experiment, len(s.Experiments))
	for i, e := range s.Experiments {
		data.Experiments[i] = experiment{
			ID:          e.ID,
			Description: e.Description,
			LastUpdated: e.LastUpdated(),
			Trials:      len(e.Trials),
			Fitness:     e.Best().Fitness,
			Complexity:  float64(e.Best().Complexity()),
			Diversity:   e.Diversity().Mean(),
		}
	}

	var ss []web.Study
	if ss, err = client.GetStudies(); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data.Studies = make([]study, len(ss))
	for i, s := range ss {
		data.Studies[i] = study{s.ID, s.Description}
	}

	// Build the box chart
	data.Chart.Box.Title = "Fitness by Experiment"
	data.Chart.Box.XTitle = "Experiment"
	data.Chart.Box.XLabels = make([]string, len(s.Experiments))
	data.Chart.Box.Data = make([][]float64, len(s.Experiments))
	data.Chart.Box.YMin = 1e10
	data.Chart.Box.YMax = 0
	for i, e := range s.Experiments {
		data.Chart.Box.XLabels[len(s.Experiments)-i-1] = e.Description
		x := e.Fitness()
		data.Chart.Box.Data[len(s.Experiments)-i-1] = []float64{
			x.Min(), x.Q25(), x.Median(), x.Q75(), x.Max(),
		}
		if data.Chart.Box.YMax < x.Max() {
			data.Chart.Box.YMax = x.Max()
		}
		if data.Chart.Box.YMin > x.Min() {
			data.Chart.Box.YMin = x.Min()
		}
	}

	// Build the scatter chart
	data.Chart.Scatter.Title = "Best Fitness vs Complexity"
	data.Chart.Scatter.XTitle = "Complexity"
	data.Chart.Scatter.XMin = 1e10
	data.Chart.Scatter.XMax = 0
	data.Chart.Scatter.YTitle = "Fitness"
	data.Chart.Scatter.YMin = 1e10
	data.Chart.Scatter.YMax = 0
	data.Chart.Scatter.Data = make([][]float64, len(s.Experiments))
	for i, e := range s.Experiments {
		f := e.Best().Fitness
		c := float64(e.Best().Complexity())

		data.Chart.Scatter.Data[len(s.Experiments)-i-1] = []float64{c, f}
		if data.Chart.Scatter.YMax < f {
			data.Chart.Scatter.YMax = f
		}
		if data.Chart.Scatter.YMin > f {
			data.Chart.Scatter.YMin = f
		}
		if data.Chart.Scatter.XMax < c {
			data.Chart.Scatter.XMax = c
		}
		if data.Chart.Scatter.XMin > c {
			data.Chart.Scatter.XMin = c
		}
	}

	t := template.New("page")
	t, _ = t.Parse(layoutTemplate)
	t, _ = t.Parse(experimentsTemplate)
	t, _ = t.Parse(chartsTemplate)
	t, _ = t.Parse(boxesTemplate)
	t, _ = t.Parse(scatterTemplate)
	t.Execute(w, data)
}

func setExperiment(w http.ResponseWriter, r *http.Request) {

	var err error
	if err = r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("server.SetStudies parse form error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Get the variables
	var sid, eid, newsid int
	if sid, err = strconv.Atoi(r.Form.Get("study-id")); err != nil {
		http.Error(w, fmt.Sprintf("server.setExperiment error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if eid, err = strconv.Atoi(r.Form.Get("experiment-id")); err != nil {
		http.Error(w, fmt.Sprintf("server.setExperiment error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	if newsid, err = strconv.Atoi(r.Form.Get("experiment-sid")); err != nil {
		http.Error(w, fmt.Sprintf("server.setExperiment error: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	desc := r.Form.Get("experiment-desc")
	delx := r.Form.Get("experiment-del")

	// Delete the study
	if delx == "on" {
		if err = client.DelExperiment(eid); err != nil {
			http.Error(w, "server.SetExperiment error deleting experiment: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/studies/%d", sid), http.StatusSeeOther)
		return
	}

	// Update the description
	if err = client.SetExperiment(newsid, eid, desc); err != nil {
		http.Error(w, "server.setExperiment error updating iteration: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/studies/%d", sid), http.StatusSeeOther)
}

func getExperimentApi(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Content-Type", "application/json")

	var err error
	var eid int
	vars := mux.Vars(r)
	if eid, err = strconv.Atoi(vars["eid"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var e web.Experiment
	if e, err = client.GetExperiment(eid); err != nil {
		log.Println("server.getExperimentApi", "get error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(e); err != nil {
		log.Println("server.getExperimentApi", "json error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addExperimentApi(w http.ResponseWriter, r *http.Request) {

	// Extract the description from the request
	var err error
	var sid int
	if sid, err = strconv.Atoi(r.FormValue("sid")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	desc := r.FormValue("desc")

	// Create the experiment
	var eid int
	if eid, err = client.AddExperiment(sid, desc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID
	var answer struct{ ExperimentID int }
	answer.ExperimentID = eid
	if err = json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addTrialApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the exeperiment ID
	var err error
	var eid int
	vars := mux.Vars(r)
	if eid, err = strconv.Atoi(vars["eid"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the description from the request
	desc := r.FormValue("desc")

	// Create the trial
	var tid int
	if tid, err = client.AddTrial(eid, desc); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID
	var answer struct{ TrialID int }
	answer.TrialID = tid
	if err = json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addIterationApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract the trial ID and population
	var err error
	var tid int
	vars := mux.Vars(r)
	if tid, err = strconv.Atoi(vars["tid"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var p evo.Population
	if err = json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the iteration
	var iid int
	if iid, err = client.AddIteration(tid, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the ID
	var answer struct{ IterationID int }
	answer.IterationID = iid
	if err = json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
