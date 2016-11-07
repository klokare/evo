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
	iterationsTemplate = `
{{define "content"}}
<div class="page-header">
	<div class="row">
		<div class="col-lg-12">
			<h1 id="buttons">{{.Experiment.Description}} - {{.Trial.Description}}</h1>
		</div>
	</div>
</div>
<div class="row">
	<div class="col-lg-12">
		<ul class="breadcrumb">
			<li><a href="/studies">Home</a></li>
			<li><a href="/studies/{{.Study.ID}}">{{.Study.Description}}</a></li>
			<li><a href="/studies/{{.Study.ID}}/experiments/{{.Experiment.ID}}">{{.Experiment.Description}}</a></li>
			<li class="active">{{.Trial.Description}}</li>
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
					<th>Iteration</th>
					<th>Updated</th>
					<th>Species</th>
					<th>Best.Fitness</th>
					<th>Best.Complexity</th>
				</tr>
			</thead>
			<tbody>
				{{range .Iterations}}
				<tr>
					<td>{{.ID}}</td>
					<td>{{.Description}}</td>
					<td>{{.Updated}}</td>
					<td>{{.Diversity}}</td>
					<td>{{.Fitness}}</td>
					<td>{{.Complexity}}</td>
				</tr>
				{{end}}
			</tbody>
		</table>
	</div>
</div>
{{end}}
	`
)

func getIterations(w http.ResponseWriter, r *http.Request) {

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
	}

	type iteration struct {
		ID          int
		Description string
		Updated     time.Time
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
		Trial      trial
		Iterations []iteration
	}

	var err error
	var sid, eid, tid int
	vars := mux.Vars(r)
	if sid, err = strconv.Atoi(vars["sid"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if eid, err = strconv.Atoi(vars["eid"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if tid, err = strconv.Atoi(vars["tid"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tde, _ := url.QueryUnescape(r.FormValue("desc"))

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

	var t web.Trial
	if t, err = client.GetTrial(tid); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if t.Description != "" {
		tde = t.Description
	}
	data.Trial = trial{t.ID, tde}
	data.Iterations = make([]iteration, len(t.Iters))
	for i, it := range t.Iters {
		if it.Description == "" {
			it.Description = fmt.Sprintf("Iteration %d", len(t.Iters)-i)
		}

		data.Iterations[i] = iteration{it.ID, it.Description, it.Updated, it.Best.Fitness, float64(it.Best.Complexity()), it.Diversity[0]}
	}

	// Build the box chart
	data.Chart.Box.Title = "Fitness by Iteration"
	data.Chart.Box.XTitle = "Iteration"
	data.Chart.Box.XLabels = make([]string, len(t.Iters))
	data.Chart.Box.Data = make([][]float64, len(t.Iters))
	data.Chart.Box.YMin = 1e10
	data.Chart.Box.YMax = 0
	for i, it := range t.Iters {
		data.Chart.Box.XLabels[len(t.Iters)-i-1] = strconv.Itoa(len(t.Iters) - i)
		x := it.Fitness
		data.Chart.Box.Data[len(t.Iters)-i-1] = []float64{
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
	data.Chart.Scatter.Data = make([][]float64, len(t.Iters))
	for i, it := range t.Iters {
		f := it.Best.Fitness
		c := float64(it.Best.Complexity())

		data.Chart.Scatter.Data[len(t.Iters)-i-1] = []float64{c, f}
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

	z := template.New("page")
	z, _ = z.Parse(layoutTemplate)
	z, _ = z.Parse(iterationsTemplate)
	z, _ = z.Parse(chartsTemplate)
	z, _ = z.Parse(boxesTemplate)
	z, _ = z.Parse(scatterTemplate)
	z.Execute(w, data)
}
