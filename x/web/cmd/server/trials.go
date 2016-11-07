package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/klokare/evo/x/web"
)

var (
	trialsTemplate = `
{{define "content"}}
<div class="page-header">
	<div class="row">
		<div class="col-lg-12">
			<h1 id="buttons">{{.Experiment.Description}}</h1>
		</div>
	</div>
</div>
<div class="row">
	<div class="col-lg-12">
		<ul class="breadcrumb">
			<li><a href="/studies">Home</a></li>
			<li><a href="/studies/{{.Study.ID}}">{{.Study.Description}}</a></li>
			<li class="active">{{.Experiment.Description}}</li>
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
					<th>Trial</th>
					<th>Iterations</th>
					<th>Last Update</th>
					<th>Diversity</th>
					<th>Best.Fitness</th>
					<th>Best.Complexity</th>
				</tr>
			</thead>
			<tbody>
				{{range .Trials}}
				<tr>
					<td>{{.ID}}</td>
					<td>{{.Description}}</td>
					<td><a href="/studies/{{$.Study.ID}}/experiments/{{$.Experiment.ID}}/trials/{{.ID}}?desc={{.DescEnc}}">{{.Iterations}}</a></td>
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
{{end}}	`
)

func getTrials(w http.ResponseWriter, r *http.Request) {

	type study struct {
		ID          int
		Description string
	}

	type experiment struct {
		ID          int
		Description string
	}

	type trial struct {
		ID          int
		Description string
		DescEnc     string
		LastUpdated time.Time
		Iterations  int
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
		Chart      chart
		Study      study
		Experiment experiment
		Trials     []trial
	}

	var err error
	var sid, eid int
	vars := mux.Vars(r)
	if sid, err = strconv.Atoi(vars["sid"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if eid, err = strconv.Atoi(vars["eid"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	var s web.Study
	if s, err = client.GetStudy(sid); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data.Study = study{s.ID, s.Description}

	var e web.Experiment
	if e, err = client.GetExperiment(eid); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	data.Experiment = experiment{e.ID, e.Description}
	data.Trials = make([]trial, len(e.Trials))
	for i, t := range e.Trials {
		if t.Description == "" {
			t.Description = fmt.Sprintf("Trial %d", len(e.Trials)-i)
		}
		data.Trials[i] = trial{t.ID, t.Description, url.QueryEscape(t.Description), t.LastUpdated(), len(t.Iters),
			t.Best().Fitness, float64(t.Best().Complexity()), t.Diversity().Mean()}
	}

	// Build the box chart
	data.Chart.Box.Title = "Fitness by Trial"
	data.Chart.Box.XTitle = "Trial"
	data.Chart.Box.XLabels = make([]string, len(e.Trials))
	data.Chart.Box.Data = make([][]float64, len(e.Trials))
	data.Chart.Box.YMin = 1e10
	data.Chart.Box.YMax = 0
	for i, t := range e.Trials {
		data.Chart.Box.XLabels[len(e.Trials)-i-1] = strconv.Itoa(len(e.Trials) - i)
		x := t.Fitness()
		data.Chart.Box.Data[len(e.Trials)-i-1] = []float64{
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
	data.Chart.Scatter.Data = make([][]float64, len(e.Trials))
	for i, t := range e.Trials {
		f := t.Best().Fitness
		c := float64(t.Best().Complexity())

		data.Chart.Scatter.Data[len(e.Trials)-i-1] = []float64{c, f}
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
	t, _ = t.Parse(trialsTemplate)
	t, _ = t.Parse(chartsTemplate)
	t, _ = t.Parse(boxesTemplate)
	t, _ = t.Parse(scatterTemplate)
	t.Execute(w, data)
}
