package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/klokare/evo/x/web"
)

var (
	studiesTemplate = `
{{define "content"}}
<div class="page-header">
	<div class="row">
		<div class="col-lg-12">
			<h1 id="buttons">Studies</h1>
		</div>
	</div>
</div>
<div class="row">
	<!--
		<ul class="breadcrumb">
			<li class="active">Home</li>
		</ul>
		
		-->		    
	<table class="table table-striped table-hover ">
		<thead>
			<tr>
				<th>ID</th>
				<th>Study</th>
				<th>Experiments</th>
				<th>Last Update</th>
				<th>Diversity</th>
				<th>Best.Fitness</th>
				<th>Best.Complexity</th>
			</tr>
		</thead>
		<tbody>
			{{range .}}
			<tr>
				<td>{{.ID}}</td>
				<td><a href="#" onclick="setModal({{.ID}}, '{{.Description}}')" data-toggle="modal" data-target="#edit-modal">{{.Description}}</a></td>
				<td><a href="/studies/{{.ID}}">{{.Experiments}}</a></td>
				<td>{{.LastUpdated}}</td>
				<td>{{.Diversity}}</td>
				<td>{{.Fitness}}</td>
				<td>{{.Complexity}}</td>
			</tr>
			{{end}}
		</tbody>
	</table>
	<a href="/studies/add" class="btn btn-primary btn-xs">add study</a>
</div>
<div class="modal" id="edit-modal">
	<div class="modal-dialog">
		<div class="modal-content">
			<form class="form-horizontal" method="post" action="/studies/edit">
				<fieldset>
					<div class="modal-header">
						<button type="button" class="close" data-dismiss="modal" aria-hidden="true">&times;</button>
						<h4 class="modal-title">Edit Study</h4>
					</div>
					<div class="modal-body">
						<input type="hidden" name="study-id" id="study-id" />
						<div class="form-group">
							<label for="study-desc" class="col-lg-2 control-label">Description</label>
							<div class="col-lg-10">
								<input type="text" class="form-control" id="study-desc" name="study-desc">
							</div>
						</div>
						<div class="form-group">
							<label for="study-del" class="col-lg-2 control-label">Delete</label>
							<div class="col-lg-10">
								<div class="checkbox">
									<label>
									<input type="checkbox" id="study-del" name="study-del"> Select to delete this study
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
	function setModal(sid, desc) {
		document.getElementById("study-id").value = sid;
		document.getElementById("study-desc").value = desc;
	}
</script>
{{end}}
	`
)

func getStudies(w http.ResponseWriter, r *http.Request) {
	log.Println("server.getStudies")
	var err error
	var ss []web.Study
	if ss, err = client.GetStudies(); err != nil {
		if err = json.NewEncoder(w).Encode(err); err != nil {
			log.Println("ERROR", "getStudies", err)
		}
		return
	}

	type study struct {
		ID          int
		Description string
		LastUpdated time.Time
		Experiments int
		Fitness     float64
		Complexity  float64
		Diversity   float64
	}

	data := make([]study, len(ss))
	for i, s := range ss {
		data[i] = study{
			ID:          s.ID,
			Description: s.Description,
			LastUpdated: s.LastUpdated(),
			Experiments: len(s.Experiments),
			Fitness:     s.Best().Fitness,
			Complexity:  float64(s.Best().Complexity()),
			Diversity:   s.Diversity().Mean(),
		}
	}

	t := template.New("page")
	t, _ = t.Parse(layoutTemplate)
	t, _ = t.Parse(studiesTemplate)
	t.Execute(w, data)
}

func setStudy(w http.ResponseWriter, r *http.Request) {

	var err error
	if err = r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("server.SetStudies parse form error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Get the variables
	var sid int
	sidx := r.Form.Get("study-id")
	if sid, err = strconv.Atoi(sidx); err != nil {
		http.Error(w, fmt.Sprintf("server.SetStudies sid %s error: %s", sidx, err.Error()), http.StatusInternalServerError)
		return
	}
	desc := r.Form.Get("study-desc")
	delx := r.Form.Get("study-del")

	// Delete the study
	if delx == "on" {
		if err = client.DelStudy(sid); err != nil {
			http.Error(w, "server.SetStudies error deleting study: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/studies", http.StatusSeeOther)
		return
	}

	// Update the description
	if err = client.SetStudy(sid, desc); err != nil {
		http.Error(w, "server.SetStudies error updating study: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/studies", http.StatusSeeOther)
}

func addStudyApi(w http.ResponseWriter, r *http.Request) {
	log.Println("server.addStudyApi")
	w.Header().Set("Content-Type", "application/json")
	var sid int
	var err error
	if sid, err = client.AddStudy("New study"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var answer struct{ StudyID int }
	answer.StudyID = sid
	if err = json.NewEncoder(w).Encode(answer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addStudy(w http.ResponseWriter, r *http.Request) {
	log.Println("server.addStudy")
	if _, err := client.AddStudy("New study"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/studies", http.StatusSeeOther)
}
